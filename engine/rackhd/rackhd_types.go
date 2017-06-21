package rackhd

import (
	"fmt"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"time"
)

type rackhdWorkflows []rackhdWorkflow

func (ws rackhdWorkflows) findLatestGraph(graph string) *rackhdWorkflow {
	var tmpWs rackhdWorkflow
	for _, w := range ws {
		if w.isMatchGraph(graph) && w.UpdatedAt.After(tmpWs.UpdatedAt) {
			tmpWs = w
		}
	}

	return &tmpWs
}

type rackhdWorkflow struct {
	ID             string    `json:"id"`
	InstanceID     string    `json:"instanceId"`
	Node           string    `json:"node"`
	Name           string    `json:"name"`
	InjectableName string    `json:"injectableName"`
	HideStatus     string    `json:"_status"`
	Status         string    `json:"status"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Context        struct {
		GraphID string `json:"graphId"`
	} `json:"context"`
}

func (w rackhdWorkflow) isMatchGraph(graph string) bool {
	return (w.InstanceID == graph ||
		w.InjectableName == graph ||
		w.Context.GraphID == graph)
}

type rackhdNodes []rackhdNode

func (nodes rackhdNodes) getIdentifierMap() map[string]rackhdNode {
	idenMaps := make(map[string]rackhdNode)
	for _, n := range nodes {
		for _, iden := range n.Identifiers {
			idenMaps[iden] = n
		}
	}
	return idenMaps
}

func (nodes rackhdNodes) findByIdentifiers(identifiers []string) *rackhdNode {
	idenMaps := nodes.getIdentifierMap()

	for _, iden := range identifiers {
		n, ok := idenMaps[iden]
		if ok {
			return &n
		}
	}

	return nil
}

type rackhdNode struct {
	ID           string   `json:"id"`
	Name         string   ` json:"name"`
	Identifiers  []string `json:"identifiers"`
	CatalogsURL  string   `json:"catalogs"`
	WorkflowsURL string   `json:"workflows"`
}

func (rn rackhdNode) getGraphStatus(api *RackhdClient, graph string) (string, error) {
	ws, err := api.getWorkflows(rn.WorkflowsURL)
	if err != nil {
		return "", err
	}

	w := ws.findLatestGraph(graph)
	if w == nil {
		return "", fmt.Errorf("Workflow %s do not have graph %s\n",
			rn.WorkflowsURL, graph)
	}

	// find workflows with node will return _status
	return w.HideStatus, nil
}

func (rNode rackhdNode) findMatchConfNode(nodes []conf.Node) *conf.Node {
	for _, n := range nodes {
		for _, iden := range n.Identifiers {
			if rNode.Name == iden {
				return &n
			}

			for _, iden2 := range rNode.Identifiers {
				if iden2 == iden {
					return &n
				}
			}
		}
	}
	return nil
}

type nodeStatus struct {
	node   rackhdNode
	status string
}
