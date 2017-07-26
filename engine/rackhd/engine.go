package rackhd

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"gitlab.adlinktech.com/lyan.hung/opps/engine"
	"log"
)

const engineName = "rackhd"

var (
	rackhdPayloadNameMap = make(map[string]*rackhdEnginePayload)
)

func Name() string {
	return engineName
}

func newName(api, graph string) string {
	baseName := fmt.Sprintf("%s:%s:%s", Name(), api, graph)
	n := baseName

	i := 0

	for _, ok := rackhdPayloadNameMap[n]; ok; _, ok = rackhdPayloadNameMap[n] {
		n = fmt.Sprintf("%s_%d", baseName, i)
		i += 1
	}
	return n
}

func init() {
	err := engine.RegistryEngine(Name(), newEngine())
	if err != nil {
		log.Printf("Engine %s registry failed: %s", Name(), err)
	}
}

func newEngine() engine.Engine {
	e := &rackhdEngine{}
	return e
}

type rackhdEngine struct {
	payloads []*rackhdEnginePayload
}

func (e *rackhdEngine) Name() string {
	return Name()
}

func (e *rackhdEngine) Discovery(reportCh chan<- conf.ScenarioReport, nodes []*conf.Node,
	data []byte) (string, error) {
	p, err := newRackhdEnginePayload(reportCh, nodes, data)
	if err != nil {
		return "", err
	}

	go p.handleStatus()
	if err := p.discovery(); err != nil {
		log.Printf("Engine %s discovery failed: %s\n", Name(), err)
		return "", err
	}

	e.payloads = append(e.payloads, p)
	rackhdPayloadNameMap[p.name] = p
	return p.name, nil
}

func (e *rackhdEngine) HandleHook(data []byte) error {
	log.Println("Hook request:", string(data))
	event := &rackhdEvent{}
	err := json.Unmarshal(data, event)
	if err != nil {
		return err
	}

	if event.Type != "graph" || event.Action != "finished" {
		return nil
	}

	if len(event.TypeID) == 0 {
		return errors.New("Hook's value InstanceID should not been empty")
	}

	for _, p := range e.payloads {
		go func(p *rackhdEnginePayload) {
			err := p.handleEvent(event)
			if err != nil {
				log.Printf("Node %s handle hook with api serever %s failed with: %s",
					event.NodeID, p.API, err)
			}
		}(p)
	}

	return nil
}
