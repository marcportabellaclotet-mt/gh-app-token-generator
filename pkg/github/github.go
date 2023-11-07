package github

import (
	"encoding/json"
	"strconv"

	"fmt"
	"io"

	"net/http"
	"strings"
	"time"

	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/helpers"
	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/output"
)

const (
	githubURL              = "https://api.github.com"
	githubInstallationsURL = "https://api.github.com/app/installations"
)

var (
	AppInst string
)

type GitHubAppInfo struct {
	AppID      string
	AppInstID  string
	GHRepo     string
	PrivateKey []byte
}

func GetInstallationID(token string, repo string) (*string, error) {

	u := strings.Join([]string{githubURL, "repos", repo, "installation"}, "/")

	var resBody struct {
		Id int `json:"id"`
	}
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return nil, fmt.Errorf("Invalid response code : %d", res.StatusCode)
	}
	if err := json.Unmarshal(b, &resBody); err != nil {
		return nil, fmt.Errorf("Problem unmarshalling Body : %s\n", err.Error())
	}
	r := strconv.Itoa(resBody.Id)
	return helpers.String(r), nil
}

func GetInstallationToken(token string, appInstID string) (*string, error) {

	u := strings.Join([]string{githubInstallationsURL, appInstID, "access_tokens"}, "/")

	var resBody struct {
		Token       string    `json:"token"`
		ExpiresAt   time.Time `json:"expires_at"`
		Permissions struct {
			Contents     string `json:"contents"`
			Metadata     string `json:"metadata"`
			PullRequests string `json:"pull_requests"`
		} `json:"permissions"`
		RepositorySelection string `json:"repository_selection"`
	}

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return nil, fmt.Errorf("Invalid response code : %d", res.StatusCode)
	}

	if err := json.Unmarshal(b, &resBody); err != nil {
		return nil, fmt.Errorf("Problem unmarshalling Body : %s\n", err.Error())
	}

	return &resBody.Token, nil
}

func GenToken(gh GitHubAppInfo, outputFormat string) error {
	key, err := helpers.LoadPEMFromBytes(gh.PrivateKey)
	if err != nil {
		return err
	}

	jwt := helpers.IssueJWTFromPEM(key, gh.AppID)
	if gh.AppInstID != "" {
		AppInst = gh.AppInstID
	} else {
		a, err := GetInstallationID(jwt, gh.GHRepo)
		if err != nil {
			return err
		}
		AppInst = *helpers.String(*a)
	}
	token, err := GetInstallationToken(jwt, AppInst)
	if err != nil {
		return err
	}
	response := output.Response{
		Token:  *token,
		Status: "success",
		Info:   "token retrieved successfully",
		Output: outputFormat,
	}

	output.ReturnResponse(response)
	return nil
}
