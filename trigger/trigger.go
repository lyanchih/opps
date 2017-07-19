package trigger

import (
	"errors"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"log"
	"sync"
)

var (
	targets           map[string]target = make(map[string]target)
	ErrTriggerNotInit                   = errors.New("Trigger still not been initialized")
)

type target interface {
	trigger([]conf.Node, []byte) error
	init([]byte) error
}

func InitTriggers(ts []conf.Trigger) error {
	for _, t := range ts {
		target, ok := targets[t.Type]
		if !ok {
			log.Printf("Trigger %s type is not exist\n", t.Type)
			continue
		}

		err := target.init(t.Data)
		if err != nil {
			log.Printf("Trigger %s type init failed: err %s", t.Type, err)
			delete(targets, t.Type)
			continue
		}

		log.Printf("Trigger %s type init succeeded\n", t.Type)
	}

	return nil
}

func Trigger(nodes []conf.Node, data []byte, types ...string) {
	wg := &sync.WaitGroup{}

	for _, t := range types {
		go func(t string) {
			defer wg.Done()

			target, ok := targets[t]
			if !ok {
				log.Printf("Trigger %s type is not support\n", t)
				return
			}

			if err := target.trigger(nodes, data); err != nil {
				log.Printf("Trigger %s type failed: %s\n", t, err)
				return
			}
		}(t)
		wg.Add(1)
	}

	wg.Wait()
}
