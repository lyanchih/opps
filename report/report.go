package report

import (
	"fmt"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"log"
)

var reporters = make(map[string]Reporter)

type Reporter interface {
	Report(ss []*conf.Scenario) ([]byte, error)
}

func addReporter(id string, r Reporter) {
	if r == nil {
		log.Println("Please do not add nil reporter")
	}

	if _, ok := reporters[id]; ok {
		log.Println("Reporter %s had been added before\n", id)
		return
	}

	reporters[id] = r
}

func GetReporter(id string) (Reporter, error) {
	r, ok := reporters[id]
	if !ok {
		return nil, fmt.Errorf("Reporter %s still not implement yet", id)
	}

	return r, nil
}
