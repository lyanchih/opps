package engine

import (
	"errors"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"log"
)

var (
	ErrEngineNotSupport = errors.New("engine is not support")
	ErrEngineExist      = errors.New("engine had been declare")
	engines             map[string]Engine
)

type Engine interface {
	Name() string
	Discovery(chan<- conf.ScenarioReport, []conf.Node, []byte) (string, error)
	HandleHook(data []byte) error
}

func init() {
	engines = make(map[string]Engine)
}

func RegistryEngine(name string, e Engine) error {
	if _, ok := engines[name]; ok {
		return ErrEngineExist
	}

	engines[name] = e
	log.Printf("Registry %s enging succeeded\n", name)
	return nil
}

func TranslateEngine(name string) (Engine, error) {
	e, ok := engines[name]
	if !ok {
		return nil, ErrEngineNotSupport
	}

	return e, nil
}
