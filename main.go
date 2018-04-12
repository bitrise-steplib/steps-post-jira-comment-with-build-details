package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-post-jira-comment-with-build-details/jira"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	UserName string `env:"user_name,required"`
	APIToken string `env:"api_token,required"`
	BaseURL  string `env:"base_url,required"`
	IsueKeys string `env:"issue_keys,required"`
	Message  string `env:"build_message,required"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	stepconf.Print(cfg)
	fmt.Println()

	encodedToken := generateBase64APIToken(cfg.UserName, cfg.APIToken)
	client := jira.NewClient(encodedToken, cfg.BaseURL)
	issueKeys := strings.Split(cfg.IsueKeys, `|`)

	var comments []jira.Comment
	for _, issueKey := range issueKeys {
		comments = append(comments, jira.Comment{Content: cfg.Message, IssuKey: issueKey})
	}

	if err := client.PostIssueComments(comments); err != nil {
		failf("Posting comments failed with error: %s", err)
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
