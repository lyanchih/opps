package trigger

import (
	"errors"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"reflect"
	"testing"
	"time"
)

var errFakeError = errors.New("fake error")

type fakeTrigger struct {
	nodes       []conf.Node
	initData    []byte
	triggerData []byte
	err         error
	wait        chan bool
}

func (f *fakeTrigger) trigger(nodes []conf.Node, data []byte) error {
	f.nodes = nodes
	f.triggerData = data
	f.wait <- true
	return nil
}

func (f *fakeTrigger) init(data []byte) error {
	if f.err != nil {
		return f.err
	}

	f.initData = data
	return nil
}

func TestInitTriggers(t *testing.T) {
	targets["fake"] = &fakeTrigger{}
	targets["err_fake"] = &fakeTrigger{err: errFakeError}

	initData := []byte("just data")
	err := InitTriggers([]conf.Trigger{
		{Type: "fake", Data: initData},
		{Type: "err_fake"},
	})
	if err != nil {
		t.Error("init triggers should not return error")
	}

	_, ok := targets["err_fake"]
	if ok {
		t.Error("Init error trigger should been deleted")
	}

	target, ok := targets["fake"]
	if !ok {
		t.Error("Init succeded trigger should not been deleted")
	}

	f, ok := target.(*fakeTrigger)
	if !ok {
		t.Error("The target should been fake trigger")
	}

	if !reflect.DeepEqual(f.initData, initData) {
		t.Error("Init triiger data should been pass into fake trigger")
	}
}

func TestTrigger(t *testing.T) {
	f := &fakeTrigger{wait: make(chan bool)}
	targets["fake"] = f
	nodes := []conf.Node{
		{Name: "node1"},
		{Name: "node2", Identifiers: []string{"foo", "bar"}},
		{Name: "node3"},
	}
	triggerData := []byte("Just trigger sample")

	Trigger(nodes, triggerData, "fake", "fake2", "other")

	select {
	case <-f.wait:
	case <-time.After(5 * time.Second):
		t.Error("Fake trigger should been trigger instead timeout")
	}

	if !reflect.DeepEqual(f.nodes, nodes) {
		t.Error("Trigger action should pass nodes into target")
	}

	if !reflect.DeepEqual(f.triggerData, triggerData) {
		t.Error("Trigger action should pass data into target")
	}
}
