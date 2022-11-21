package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/github"
	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/helpers"
	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/output"
	"github.com/marcportabellaclotet-mt/gh-app-token-generator/pkg/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "gh-app-token-generator",
		Short:         fmt.Sprintf("GitHub App token generator tool (version %s)", version.VERSION),
		RunE:          genToken,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version.VERSION,
	}
	privateKey   []byte
	err          error
	OutputFormat string
)

func Init() {
	rootCmd.PersistentFlags().StringP("gh-app-id", "a", "", `GitHub App id value.
Alternatively, you can set this parameter using GH_APP_ID environment variable.
`)
	rootCmd.PersistentFlags().StringP("gh-app-installation-id", "i", "", `GitHub App installation id value.
Alternatively, you can set this parameter using GH_APP_INSTALLATION_ID environment variable.
`)
	rootCmd.PersistentFlags().StringP("gh-app-private-key", "p", "", `GitHub App private key value in base64 format.
Alternatively, you can set this parameter using GH_APP_PRIVATE_KEY environment variable.
`)
	rootCmd.PersistentFlags().StringP("gh-app-private-key-path", "f", "", `GitHub App private key file path (gh-app-private-key flag takes precedence).
Alternatively, you can set this parameter using GH_APP_PRIVATE_KEY_PATH environment variable.
`)
	rootCmd.PersistentFlags().StringP("output", "o", "json", `Output Format [json|txt].
json -> Returns a json struct : {"token":"ghtoken","status":"failed|success","info":""}
txt -> Returns the github token when cmd is successful. Otherwise returns error info.
`)

	err := rootCmd.Execute()
	if err != nil {
		response := output.Response{
			Token:  "error",
			Status: "failed",
			Info:   err.Error(),
			Output: OutputFormat,
		}
		output.ReturnResponse(response)
	}
}

func genToken(ccmd *cobra.Command, args []string) error {

	appID, _ := ccmd.Flags().GetString("gh-app-id")
	appInstID, _ := ccmd.Flags().GetString("gh-app-installation-id")
	appPrivateKeyb64, _ := ccmd.Flags().GetString("gh-app-private-key")
	appPrivateKeyPath, _ := ccmd.Flags().GetString("gh-app-private-key-path")
	OutputFormat, _ = ccmd.Flags().GetString("output")
	appID = helpers.CheckEnvVars(appID, "GH_APP_ID")
	appInstID = helpers.CheckEnvVars(appInstID, "GH_APP_INSTALLATION_ID")
	appPrivateKeyb64 = helpers.CheckEnvVars(appPrivateKeyb64, "GH_APP_PRIVATE_KEY")
	appPrivateKeyPath = helpers.CheckEnvVars(appPrivateKeyPath, "GH_APP_PRIVATE_KEY_PATH")
	if appID == "" {
		return errors.New("gh-app-id parameter is not set")
	}
	if appInstID == "" {
		return errors.New("gh-app-installation-id parameter is not set")
	}
	if appPrivateKeyb64 != "" {
		privateKey, err = base64.StdEncoding.DecodeString(appPrivateKeyb64)
		if err != nil {
			return errors.New("PEM secret should be base64 encoded")
		}
	} else if appPrivateKeyPath == "" {
		return errors.New("gh-private-key or gh-private-key-path parameters are not set")
	} else {
		privateKey, err = ioutil.ReadFile(appPrivateKeyPath)
		if err != nil {
			return errors.New(err.Error())
		}
	}
	ghApp := github.GitHubAppInfo{
		AppID:      appID,
		AppInstID:  appInstID,
		PrivateKey: privateKey,
	}
	github.GenToken(ghApp, OutputFormat)
	return nil
}
