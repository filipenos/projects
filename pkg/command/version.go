package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/filipenos/projects/pkg/log"
	"github.com/spf13/cobra"
)

// Version, Commit This variables is filled on build time
var (
	Version string
	Commit  string
)

var checkUpdate bool

// GitHubRelease representa a estrutura do release no GitHub
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Name    string `json:"name"`
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of projects",
	Run: func(cmd *cobra.Command, args []string) {
		currentVersion := Version
		if currentVersion == "" {
			currentVersion = "dev"
		}

		log.Infof("Version %s, Build %s", currentVersion, Commit)

		if checkUpdate {
			log.Infof("Checking for updates...")
			latestRelease, err := getLatestRelease()
			if err != nil {
				log.Warnf("Failed to check for updates: %v", err)
				return
			}

			latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
			currentClean := strings.TrimPrefix(currentVersion, "v")

			if latestVersion != currentClean && latestVersion != "" {
				log.Infof("New version available: %s (current: %s)", latestVersion, currentClean)
				log.Infof("Download: %s", latestRelease.HTMLURL)
				log.Infof("Update with: go install github.com/filipenos/projects@latest")
			} else {
				log.Infof("You are using the latest version")
			}
		}
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&checkUpdate, "check-update", "c", false, "Check for new versions on GitHub")
	rootCmd.AddCommand(versionCmd)
}

// getLatestRelease busca o último release do GitHub
func getLatestRelease() (*GitHubRelease, error) {
	url := "https://api.github.com/repos/filipenos/projects/releases/latest"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}
