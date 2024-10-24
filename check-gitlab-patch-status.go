package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const API_URL = "https://%s/api/v4/version"
const VERSION_CHECK_URL = "https://version.gitlab.com/check.json?gitlab_info=%s"

type GitLabStatus struct {
	Version string `json:"version"`
}

type VersionCheck struct {
	RecommendedVersions    []string `json:"latest_stable_versions"`
	Severity               string   `json:"severity"`
	CriticalVulnerabilitiy bool     `json:"critical_vulnerability"`
}

const (
	RETURN_STATE_OK int = iota
	RETURN_STATE_WARNING
	RETURN_STATE_CRITICAL
	RETURN_STATE_UNKNOWN
)

const (
	SEVERITY_SUCCESS = "success"
	SEVERITY_WARNING = "warning"
	SEVERITY_DANGER  = "danger"
)

var host = flag.String("host", "", "Hostname of the GitLab installation")
var token = flag.String("token", "", "Private Access Token")

func init() {
	flag.StringVar(host, "H", "", "Hostname of the GitLab installation")
	flag.StringVar(token, "t", "", "Private Access Token")
}

func main() {
	flag.Parse()

	versionCheck := getVersionCriticality(getLocalGitLabVersion(*host, *token))

	if versionCheck.CriticalVulnerabilitiy {
		fmt.Println("CRITICAL: A critical vulnerability has been found in your GitLab installation. Recommended Versions are: " + strings.Join(versionCheck.RecommendedVersions, ", "))
		os.Exit(RETURN_STATE_CRITICAL)
	} else {
		switch versionCheck.Severity {
		case SEVERITY_DANGER:
			fmt.Println("WARN: It is strongly recommended to update your GitLab installation. Recommended Versions are: " + strings.Join(versionCheck.RecommendedVersions, ", "))
			os.Exit(RETURN_STATE_WARNING)
		case SEVERITY_WARNING:
			fmt.Println("INFO: A new version of GitLab is available. Recommended Versions are: " + strings.Join(versionCheck.RecommendedVersions, ", "))
			os.Exit(RETURN_STATE_OK)
		case SEVERITY_SUCCESS:
			fmt.Println("OK: Your GitLab installation is up to date.")
			os.Exit(RETURN_STATE_OK)
		}
	}

	os.Exit(RETURN_STATE_OK)
}

func getLocalGitLabVersion(host string, token string) (version []byte) {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(API_URL, host), nil)
	request.Header.Set("PRIVATE-TOKEN", token)

	response, requestErr := http.DefaultClient.Do(request)
	if requestErr != nil {
		fmt.Println(requestErr)
		os.Exit(RETURN_STATE_UNKNOWN)
	}

	body, parseErr := io.ReadAll(response.Body)
	if parseErr != nil {
		fmt.Println(parseErr)
		os.Exit(RETURN_STATE_UNKNOWN)
	}

	var status GitLabStatus
	unmarshalErr := json.Unmarshal(body, &status)

	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
		os.Exit(RETURN_STATE_UNKNOWN)
	}

	versionJSON, _ := json.Marshal(status)
	return versionJSON
}

func getVersionCriticality(version []byte) VersionCheck {
	encodedVersion := base64.StdEncoding.EncodeToString(version)
	response, requestErr := http.Get(fmt.Sprintf(VERSION_CHECK_URL, encodedVersion))

	if requestErr != nil {
		fmt.Println(requestErr)
		os.Exit(RETURN_STATE_UNKNOWN)
	}

	body, parseErr := io.ReadAll(response.Body)
	if parseErr != nil {
		fmt.Println(parseErr)
		os.Exit(RETURN_STATE_UNKNOWN)
	}

	var versionCheck VersionCheck
	unmarshalErr := json.Unmarshal(body, &versionCheck)

	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
		os.Exit(RETURN_STATE_UNKNOWN)
	}

	return versionCheck
}
