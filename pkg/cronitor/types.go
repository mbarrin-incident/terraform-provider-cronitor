// Copyright (c) HashiCorp, Inc.

package cronitor

type Monitor struct {
	Name             string   `json:"name"`
	Assertions       []string `json:"assertions"`
	Disabled         bool     `json:"disabled"`
	FailureTolerance *int     `json:"failure_tolerance"`
	GraceSeconds     *int     `json:"grace_seconds"`
	Group            *string  `json:"group"`
	Key              string   `json:"key"`
	Notify           []string `json:"notify"`
	Paused           bool     `json:"paused"`
	Platform         string   `json:"platform"`
	RealertInterval  string   `json:"realert_interval"`
	Request          struct {
		URL             string            `json:"url"`
		Headers         map[string]string `json:"headers"`
		Cookies         map[string]string `json:"cookies"`
		Body            string            `json:"body"`
		Method          string            `json:"method"`
		TimeoutSeconds  int               `json:"timeout_seconds"`
		Regions         []string          `json:"regions"`
		FollowRedirects bool              `json:"follow_redirects"`
		VerifySsl       bool              `json:"verify_ssl"`
	} `json:"request"`
	Running           bool     `json:"running"`
	Schedule          string   `json:"schedule"`
	ScheduleTolerance *int     `json:"schedule_tolerance"`
	Tags              []string `json:"tags"`
	Timezone          *string  `json:"timezone"`
	Type              string   `json:"type"`
	Environments      []string `json:"environments"`
}
