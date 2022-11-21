package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/helpers"
	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/output"
)

var (
	githubURL = "https://api.github.com/app/installations"
)

type GitHubAppInfo struct {
	AppID      string
	AppInstID  string
	PrivateKey []byte
}

func GetInstallationToken(token string, appInstID string) (*string, error) {

	u := strings.Join([]string{githubURL, appInstID, "access_tokens"}, "/")

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

	client := &http.Client{}

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

	b, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return nil, errors.New(fmt.Sprintf("Invalid response code : %d", res.StatusCode))
	}

	if err := json.Unmarshal(b, &resBody); err != nil {
		return nil, errors.New(fmt.Sprintf("Problem unmarshalling Body : %s\n", err.Error()))
	}

	return &resBody.Token, nil
}

func GenToken(gh GitHubAppInfo, outputFormat string) {
	key, err := helpers.LoadPEMFromBytes(gh.PrivateKey)
	if err != nil {
		response := output.Response{
			Token:  "error",
			Status: "failed",
			Info:   fmt.Sprintf("Unable to load PEM - %s", err.Error()),
			Output: outputFormat,
		}
		output.ReturnResponse(response)
		return //fmt.Println(jwt)

	}

	jwt := helpers.IssueJWTFromPEM(key, gh.AppID)

	token, err := GetInstallationToken(jwt, gh.AppInstID)
	if err != nil {
		response := output.Response{
			Token:  "error",
			Status: "failed",
			Info:   fmt.Sprintf("Unable to get installation token : %s", err),
			Output: outputFormat,
		}
		output.ReturnResponse(response)
		return
	}
	response := output.Response{
		Token:  *token,
		Status: "success",
		Info:   "token retrieved successfully",
		Output: outputFormat,
	}

	output.ReturnResponse(response)
}
