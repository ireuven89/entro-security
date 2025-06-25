package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Secret   string `json:"secret"`
}

func StartScan(owner, repo, repoToken string) (ScanResponse, error) {
	var scanResponse ScanResponse
	var commits []string
	var page int

	for {
		req, err := http.NewRequest("GET", fmt.Sprintf(commitsApi, owner, repo, page), nil)

		if err != nil {
			fmt.Println("Error creating request:", err)
			return ScanResponse{}, fmt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("Authorization", "token "+repoToken)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Error executing request:", err)
			return ScanResponse{}, fmt.Errorf("error executing request: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return ScanResponse{}, fmt.Errorf("error reading response: %v", err)
		}

		if err = json.Unmarshal(body, &commits); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return ScanResponse{}, fmt.Errorf("error unmarshalling response: %v", err)
		}

		if len(commits) == 0 {
			return ScanResponse{}, nil
		}
	}

	return scanResponse, nil
}

func scanCommit(owner, repo, commitSha string) ([]FoundSecret, error) {
	//	req, err := http.NewRequest("GET", fmt.Sprintf(commitApi, owner, repo, commitSha), nil)

	/*	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}*/

	return nil, nil
}

func isSecreteExists(s string) (bool, error) {
	re, err := regexp.Compile(base64Encoded)
	if err != nil {
		return false, err
	}

	return re.MatchString(s), nil
}
