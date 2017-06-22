package trigger

import (
	"encoding/json"
	"fmt"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"gitlab.adlinktech.com/lyan.hung/opps/utils"
	"strings"
)

const (
	slackHost            = "https://hooks.slack.com/services/"
	slackTargetName      = "slack"
	slackDefaultUsername = "opps-bot"
	slackDefaultIcon     = ":ghost:"
)

func init() {
	targets["slack"] = &slackTarget{}
}

type slackTarget struct {
	*slackPayload
}

func (t *slackTarget) init(data []byte) error {
	p, err := newSlackPayload(data)
	if err != nil {
		return err
	}

	t.slackPayload = p
	return nil
}

func (t *slackTarget) trigger(nodes []conf.Node, data []byte) error {
	if t.slackPayload == nil {
		return ErrTriggerNotInit
	}

	p := *t.slackPayload
	p.Text = fmt.Sprintf("Slack Trigger nodes %s had became succeeded status",
		nodes)
	bs, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	return utils.PostOnly(p.URL, bs)
}

type slackPayload struct {
	URL       string `json:"url"`
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

func newSlackPayload(data []byte) (*slackPayload, error) {
	p := &slackPayload{}
	err := json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(p.URL, slackHost) {
		return nil, fmt.Errorf("Slack url %s is not valid\n", p.URL)
	}

	if len(p.Username) == 0 {
		p.Username = slackDefaultUsername
	}

	if len(p.IconEmoji) == 0 {
		p.IconEmoji = slackDefaultIcon
	}

	return p, nil
}
