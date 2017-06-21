package engine

import (
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"testing"
)

type fakeEngine struct{}

func (f *fakeEngine) Name() string {
	return "fake"
}

func (f *fakeEngine) Discovery(chan<- conf.ScenarioReport,
	[]conf.Node, []byte) (string, error) {
	return "", nil
}

func (f *fakeEngine) HandleHook(data []byte) error {
	return nil
}

func TestRegistryEngine(t *testing.T) {
	f := &fakeEngine{}
	err := RegistryEngine(f.Name(), f)
	if err != nil {
		t.Error("First registry engine should be save")
	}

	err = RegistryEngine(f.Name(), f)
	if err != ErrEngineExist {
		t.Error("Error should return when registry same engine twice")
	}
}

func TestTranslateEngine(t *testing.T) {
	e, err := TranslateEngine("notSupport")
	if err != ErrEngineNotSupport {
		t.Error("Error should return if name of engine is not support")
	}

	e, err = TranslateEngine("fake")
	if err != nil {
		t.Error("Registry engine should been place into engines map")
	}

	if e.Name() != "fake" {
		t.Error("Registry engine name should be equale")
	}

}
