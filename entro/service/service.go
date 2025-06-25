package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var base64Encoded = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$"
var commitsApi = "https://api.github.com/repos/%s/%s/commits?page=%d&per_page=10"
var commitApi = "https://api.github.com/repos/%s/%s/commits/%s"

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
	Filename string `json:"filename"`
	Sha      string `json:"sha"`
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
			return nil, err

		}
		for _, file := range resp.Files {
			if isSecreteExists(file.Patch) {
				fmt.Printf("found secret in  %s\n", file.Filename)
				foundSecrets = append(foundSecrets, FoundSecret{Filename: file.Filename, Sha: file.Sha})
			}
		}
	}

	return nil, nil
}

func isSecreteExists(s string) bool {
	re, err := regexp.Compile(base64Encoded)
	if err != nil {
		return false
	}

	return re.MatchString(s)
}
