package rackhd

import (
	"gitlab.adlinktech.com/lyan.hung/opps/utils"
	"net/url"
	"path"
)

type RackhdClient struct {
	u url.URL
}

func NewRackhdClient(api string) (*RackhdClient, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	return &RackhdClient{
		u: *u,
	}, nil
}

func (c *RackhdClient) apiJoin(elems ...string) string {
	u := c.u
	u.Path = path.Join(append([]string{u.Path}, elems...)...)
	return u.String()
}

func (c *RackhdClient) hostJoin(elems ...string) string {
	u := c.u
	u.Path = path.Join(elems...)
	return u.String()
}

func (c *RackhdClient) getWorkflows(wsURL string) (ws rackhdWorkflows, err error) {
	err = utils.Get(c.hostJoin(wsURL), &ws)
	return
}

func (c *RackhdClient) getWorkflow(id string) (w *rackhdWorkflow, err error) {
	w = &rackhdWorkflow{}
	if err = utils.Get(c.apiJoin("workflows", id), w); err != nil {
		w = nil
	}
	return
}

func (c *RackhdClient) getNodes() (ns rackhdNodes, err error) {
	err = utils.Get(c.apiJoin("nodes"), &ns)
	return
}

func (c *RackhdClient) getNode(id string) (n *rackhdNode, err error) {
	n = &rackhdNode{}
	if err = utils.Get(c.apiJoin("nodes", id), n); err != nil {
		n = nil
	}
	return
}

func (c *RackhdClient) lookup(q string) (l *rackhdLookup, err error) {
	ls := []rackhdLookup{}
	if err = utils.Get(c.apiJoin("lookups")+"?q="+q, &ls); err != nil || len(ls) == 0 {
		l = nil
	} else {
		l = &ls[0]
	}
	return
}
