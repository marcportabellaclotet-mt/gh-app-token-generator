# Github App Token Generator

A simple standalone app written in go to generate a temporal github token using github app credentials
Most of the logic for this app is borrowed from [mlioo/go-github-app-token-gen](https://github.com/mlioo/go-github-app-token-generator) github action

## Usage

```
GitHub App token generator tool

Usage:
  gh-app-token-generator [flags]

Flags:
  -a, --gh-app-id string                 GitHub App id value.
                                         Alternatively, you can set this parameter using GH_APP_ID environment variable.

  -i, --gh-app-installation-id string    GitHub App installation id value.
                                         Alternatively, you can set this parameter using GH_APP_INSTALLATION_ID environment variable.

  -p, --gh-app-private-key string        GitHub App private key value in base64 format.
                                         Alternatively, you can set this parameter using GH_APP_PRIVATE_KEY environment variable.

  -f, --gh-app-private-key-path string   GitHub App private key file path (gh-app-private-key flag takes precedence).
                                         Alternatively, you can set this parameter using GH_APP_PRIVATE_KEY_PATH environment variable.

  -h, --help                             help for gh-app-token-generator

  -o, --output string                    Output Format [json|txt|export].
                                         json -> Returns a json struct : {"token":"ghtoken","status":"failed|success","info":""}
                                         txt -> Returns the github token when cmd is successful. Otherwise    returns error info.
                                         export -> Returns export GH_TOKEN command when cmd is successful. Otherwise returns error info.
                                        (default "json")

````
## Example

Using gh-app-private-key in base64 format:
```
# gh-app-token-generator -gh-app-id 123456 -gh-app-installation-id 12345678 --gh-app-private-key LSxxxxx

{"token":"ghs_xxxxxxxxxxxxxxxxx","status":"success","info":"token retrieved successfully"}

```
Using gh-app-private-key-path parameter to read github app key from file:
```
# gh-app-token-generator -gh-app-id 123456 -gh-app-installation-id 12345678 --gh-app-private-key-path /tmp/key.em

{"token":"ghs_xxxxxxxxxxxxxxxxx","status":"success","info":"token retrieved successfully"}

```
Using export output to generate export GH_TOKEN command:
```
# gh-app-token-generator -gh-app-id 123456 -gh-app-installation-id 12345678 --gh-app-private-key-path /tmp/key.em -o export

export GH_TOKEN="ghs_xxxxxxxxxxxxxxxxx"

```
