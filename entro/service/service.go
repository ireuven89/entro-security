package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var (
	commitsApi          = "https://api.github.com/repos/%s/%s/commits?page=%d&per_page=10"
	commitApi           = "https://api.github.com/repos/%s/%s/commits/%s"
	secretKeyPattern    = regexp.MustCompile("[A-Za-z0-9+/]{40}")
	AccessKeyPattern    = regexp.MustCompile("(?:AKIA|ASIA|AIDA|AROA)[A-HJ-NP-Z2-7]{16,20}")
	SessionTokenPattern = regexp.MustCompile("[A-Za-z0-9+/]{60,}={0,2}")
)

type CommitsResponse struct {
	Commits []Commit
}
type Commit struct {
	Sha string `json:"sha"`
}

type ScanRequest struct {
	RepositoryOwner string
	RepositoryName  string
	RepositoryToken string
}

type CommitResponse struct {
	Sha   string         `json:"sha"`
	Files []FileResponse `json:"files"`
}

type FileResponse struct {
	Sha      string `json:"sha"`
	Filename string `json:"filename"`
	Patch    string `json:"patch"`
}

type ScanResponse struct {
	FoundSecrets []FoundSecret `json:"found_secrets"`
}

type FoundSecret struct {
	Filename string   `json:"filename"`
	Sha      string   `json:"sha"`
	Secret   []Secret `json:"secrets"`
}

type Secret struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
}

func StartScan(owner, repo, repoToken string) (ScanResponse, error) {
	var scanResponse ScanResponse
	var page int

	for {
		commits, err := getCommits(fmt.Sprintf(commitsApi, owner, repo, page), repoToken)

		if err != nil {
			return scanResponse, err
		}

		if len(commits) == 0 {
			break
		}

		foundSecrets, err := scanCommits(commits, owner, repo, repoToken)
		if err != nil {
			return scanResponse, err

		}

		scanResponse.FoundSecrets = append(scanResponse.FoundSecrets, foundSecrets...)

		page++
	}

	return scanResponse, nil
}

func makeHttpGet(url, token string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got error response with status: %v, body: %v", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	return nil
}

func getCommits(repoUrl, token string) ([]Commit, error) {
	var commits []Commit

	if err := makeHttpGet(repoUrl, token, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

func getCommit(url, token string) (*CommitResponse, error) {
	var commitResponse CommitResponse

	if err := makeHttpGet(url, token, &commitResponse); err != nil {
		return nil, err
	}

	return &commitResponse, nil
}

func scanCommits(commits []Commit, owner, repo, token string) ([]FoundSecret, error) {
	var foundSecrets []FoundSecret

	for _, commit := range commits {
		resp, err := getCommit(fmt.Sprintf(commitApi, owner, repo, commit.Sha), token)
		if err != nil {
			return nil, fmt.Errorf("error executing request: %v", err)

		}
		for _, file := range resp.Files {
			secrets := findFileSecrets(file.Patch)

			if len(secrets) > 0 {
				foundSecrets = append(foundSecrets, FoundSecret{Filename: file.Filename, Sha: file.Sha, Secret: secrets})
			}
		}
	}

	return foundSecrets, nil
}

func findFileSecrets(s string) []Secret {
	var secrets []Secret

	if accessKey := AccessKeyPattern.FindString(s); accessKey != "" {
		secrets = append(secrets, Secret{Type: "access_key", Secret: accessKey})
	}

	if sessionKey := SessionTokenPattern.FindString(s); sessionKey != "" {
		secrets = append(secrets, Secret{Type: "session_key", Secret: sessionKey})
	}

	if secretKey := secretKeyPattern.FindString(s); secretKey != "" {

		secrets = append(secrets, Secret{Type: "secret_key", Secret: secretKey})
	}

	return secrets
}
