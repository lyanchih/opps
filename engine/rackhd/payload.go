package rackhd

import (
	"encoding/json"
	"errors"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"log"
	"sync"
)

var ErrEngineDataNotValid = errors.New("Engine data is not valid")

type rackhdEnginePayload struct {
	Graph             string `json:"graph"`
	API               string `json:"api"`
	client            *RackhdClient
	name              string
	nodes             []*conf.Node
	status            conf.ReportStatus
	nodeStatusMap     map[string]string
	nodeStatusChannel chan nodeStatus
	reportCh          chan<- conf.ScenarioReport
}

func newRackhdEnginePayload(reportCh chan<- conf.ScenarioReport, nodes []*conf.Node,
	data []byte) (*rackhdEnginePayload, error) {
	if len(data) == 0 {
		return nil, ErrEngineDataNotValid
	}

	payload := &rackhdEnginePayload{}
	err := json.Unmarshal(data, payload)
	if err != nil {
		return nil, err
	}

	c, err := NewRackhdClient(payload.API)
	if err != nil {
		return nil, err
	}

	payload.client = c
	payload.nodes = nodes
	payload.status = conf.ReportPendingStatus
	payload.reportCh = reportCh
	payload.name = newName(payload.API, payload.Graph)
	payload.nodeStatusMap = make(map[string]string)
	payload.nodeStatusChannel = make(chan nodeStatus, 10)

	return payload, err
}

func (p *rackhdEnginePayload) handleStatus() {
	for {
		nodeStatus, ok := <-p.nodeStatusChannel
		if !ok {
			break
		}
		switch nodeStatus.status {
		case "succeeded":
			if p.status != conf.ReportFailedStatus {
				_, ok := p.nodeStatusMap[nodeStatus.node.ID]
				if !ok {
					p.nodeStatusMap[nodeStatus.node.ID] = nodeStatus.status
					p.status = conf.ReportRunningStatus
					if len(p.nodeStatusMap) == len(p.nodes) {
						p.status = conf.ReportSucceededStatus
						log.Printf("Rackhd %s had deploy succeeded\n", p.name)
						p.reportCh <- conf.ScenarioReport{
							Name:   p.name,
							Status: p.status,
						}
						close(p.nodeStatusChannel)
					}
				}
			}
		case "failed":
			p.status = conf.ReportFailedStatus
			log.Printf("Rackhd %s had deploy failed\n", p.name)
			p.reportCh <- conf.ScenarioReport{
				Name:   p.name,
				Status: p.status,
			}
			close(p.nodeStatusChannel)
		default:
			log.Printf("Node %s have incorrect status %s\n",
				nodeStatus.node.ID, nodeStatus.status)
		}
	}
}

func (p *rackhdEnginePayload) handleEvent(event *rackhdEvent) error {
	// request for this event's workflow
	workflow, err := p.client.getWorkflow(event.TypeID)
	if err != nil {
		return err
	}

	// return if event do not belong to this graph
	if !workflow.isMatchGraph(p.Graph) {
		return nil
	}

	// return if event do not belong to any node
	nodeID := event.NodeID
	if len(workflow.Node) != 0 {
		nodeID = workflow.Node
	}
	if len(nodeID) == 0 {
		return nil
	}

	// request for this event's node
	rNode, err := p.client.getNode(nodeID)
	if err != nil {
		return err
	}

	n := p.findNode(rNode)
	if n == nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	p.nodeStatusChannel <- nodeStatus{
		node:   *rNode,
		status: workflow.Status,
	}
	return nil
}

func (p *rackhdEnginePayload) discovery() error {
	log.Printf("Using %s engine to discovery following nodes: %v",
		Name(), p.nodes)

	rackNodes, err := p.client.getNodes()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, n := range p.nodes {
		wg.Add(1)
		go func(n *conf.Node) {
			defer wg.Done()

			rNode := rackNodes.findByIdentifiers(n.Identifiers)
			if rNode == nil {
				log.Printf("Node identifiers %s is not discovered yet\n", n.Identifiers)
				return
			}
			if l, err := p.client.lookup(n.Identifiers[0]); err != nil {
				log.Printf("Lookup %s failed: %s\n", n.Identifiers[0], err)
			} else {
				n.IP, n.MAC = l.IP, l.MAC
			}

			s, err := rNode.getGraphStatus(p.client, p.Graph)
			if err != nil {
				log.Printf("Node identifiers %s verify workflows failed: %s\n", n.Identifiers, err)
				return
			}

			if s == "succeeded" {
				p.nodeStatusChannel <- nodeStatus{
					node:   *rNode,
					status: s,
				}
			}
		}(n)
	}
	wg.Wait()
	return nil
}

func (p *rackhdEnginePayload) findNode(rNode *rackhdNode) *conf.Node {
	for _, n := range p.nodes {
		for _, iden := range n.Identifiers {
			if rNode.Name == iden {
				return n
			}

			for _, iden2 := range rNode.Identifiers {
				if iden2 == iden {
					return n
				}
			}
		}
	}
	return nil
}
