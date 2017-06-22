package trigger

import (
	"fmt"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"testing"
)

// Target test sample
func TestInitTrigger(t *testing.T) {
	target, ok := targets["slack"]
	if !ok {
		t.Error("Slack target should been added first")
	}

	url := slackHost + "foobar"
	initData := []byte(`{"url": "` + url + `", "channel": "#c", "username": "foo", "icon_emoji": "bar"}`)
	InitTriggers([]conf.Trigger{{Type: "slack", Data: initData}})
	target, ok = targets["slack"]
	if !ok {
		t.Error("Slack target should not been deleted after init trigger")
	}

	s, ok := target.(*slackTarget)
	if !ok {
		t.Error("Target can not been covert to slack target")
	}

	p := s.slackPayload
	if p == nil {
		t.Error("slack payload should not been nil after init")
	}

	if p.URL != url || p.Channel != "#c" ||
		p.Username != "foo" || p.IconEmoji != "bar" {
		t.Errorf(`slack payload init failed,
expect: %s
real: %v`, string(initData), p)
	}
}

func TestSlackInit(t *testing.T) {
	target := &slackTarget{}

	if target.init(nil) == nil {
		t.Error("Error should return if init with nil parameter")
	}

	if target.init([]byte(`{"url":"fake_url"}`)) == nil {
		t.Error("Error should return if init with wrong url parameter")
	}

	if err := target.init([]byte(fmt.Sprintf(`{"url":"%sfoobar"}`, slackHost))); err != nil {
		t.Error("slack target init failed: ", err)
	}
}

func TestSlackTrigger(t *testing.T) {
	target := &slackTarget{}

	err := target.trigger(nil, nil)
	if err != ErrTriggerNotInit {
		t.Error("Error should return if not init trigger first")
	}
}

func TestSlackPayload(t *testing.T) {
	dataTemplate := `{"url": "%s", "channel": "%s", "username": "%s", "icon_emoji": "%s"}`
	p, err := newSlackPayload(
		[]byte(fmt.Sprintf(dataTemplate, "https://fake", "", "", "")))
	if err == nil {
		t.Error("Error should return if init with wrong url parameter")
	}

	p, err = newSlackPayload(
		[]byte(fmt.Sprintf(dataTemplate, slackHost+"foobar", "", "", "")))
	if err != nil {
		t.Error("new slack payload failed: ", err)
	}

	if p.Username != slackDefaultUsername {
		t.Error("slack username should been default if it is empty")
	}

	if p.IconEmoji != slackDefaultIcon {
		t.Error("slack icon should been default if it is empty")
	}

	p, err = newSlackPayload(
		[]byte(fmt.Sprintf(dataTemplate, slackHost+"foobar", "chan", "foo", "icon")))
	if err != nil {
		t.Error("new slack payload failed: ", err)
	}

	if p.Channel != "chan" {
		t.Error("slack channel should been same as channel in input data")
	}

	if p.Username != "foo" {
		t.Error("slack channel should been same as channel in input data")
	}

	if p.IconEmoji != "icon" {
		t.Error("slack channel should been same as channel in input data")
	}
}
