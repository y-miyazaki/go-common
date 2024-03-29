package slack

import "fmt"

// SlackMessageLimit sets slack message limit.
const SlackMessageLimit int = 1800

// GetSlackAWSTitle returns the slack title.
func GetSlackAWSTitle(function, accountID, region, service, env string) string {
	title := function
	// AccountID
	if accountID != "" {
		title = fmt.Sprintf("%s | AccountID: %s", title, accountID)
	}
	// Region
	if region != "" {
		title = fmt.Sprintf("%s | Region: %s", title, region)
	}
	// Service
	if service != "" {
		title = fmt.Sprintf("%s | Service: %s", title, service)
	}
	// Env
	if env != "" {
		title = fmt.Sprintf("%s | Env: %s", title, env)
	}
	return title
}

// GetSlackLog return slack format log.
func GetSlackLog(log string) string {
	if len([]byte(log)) > SlackMessageLimit {
		return fmt.Sprintf("```%s\n...(too long message)```", log[:SlackMessageLimit])
	}
	return fmt.Sprintf("```%s```", log)
}
