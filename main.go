package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-post-jira-comment-with-build-details/network"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type config struct {
	JiraUserName string `env:"jira_user_name,required"`
	JiraAPIToken string `env:"jira_api_token,required"`
	JiraBaseURL  string `env:"jira_base_url,required"`
	JiraIsueKeys string `env:"jira_issue_keys,required"`
	Comment      string `env:"comment,required"`
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	stepconf.Print(cfg)
	fmt.Println()

	encodedToken := generateBase64APIToken(cfg.JiraUserName, cfg.JiraAPIToken)

	client := network.New(encodedToken, cfg.JiraBaseURL)

	issueKeys := strings.Split(cfg.JiraIsueKeys, `|`)

	var comments []network.Comment
	for _, issueKey := range issueKeys {
		comments = append(comments, network.Comment{Content: cfg.Comment, IssuKey: issueKey})
	}

	if err := client.PostIssueComments(comments); err != nil {
		failf("posting comments failed with error: %s", err)
	}

	fmt.Println()
}

func generateBase64APIToken(userName string, apiToken string) string {
	v := userName + `:` + apiToken
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
