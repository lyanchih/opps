package main

import (
	"github.com/spf13/cobra"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"gitlab.adlinktech.com/lyan.hung/opps/engine"
	"gitlab.adlinktech.com/lyan.hung/opps/trigger"
	"log"
)

var (
	reportCh    chan conf.ScenarioReport
	scenarioMap map[string]*conf.Scenario
	reportDone  chan bool
)

func init() {
	reportCh = make(chan conf.ScenarioReport, 10)
	scenarioMap = make(map[string]*conf.Scenario)
	reportDone = make(chan bool, 1)
}

func initOpps() (ss []*conf.Scenario, err error) {
	if err := initConfig(); err != nil {
		log.Println("Init config with failed:", err)
		return nil, err
	}

	ss, err = initScenarios()
	if err != nil {
		log.Println("Init scenarios with failed:", err)
		return nil, err
	}

	return ss, nil
}

func initConfig() error {
	log.Println("Ready to parse config file", config)
	_, err := conf.ParseConf(config)
	return err
}

func initScenarios() ([]*conf.Scenario, error) {
	cfg, err := conf.GetConf()
	if err != nil {
		return nil, err
	}

	for _, s := range cfg.Scenarios {
		if len(s.Engine) == 0 {
			log.Printf("Scenaria object %v should have non-empty engine\n", s)
			continue
		}

		e, err := engine.TranslateEngine(s.Engine)
		if err != nil {
			log.Printf("Translate engine %s failed: %s\n", s.Engine, err)
			continue
		}

		id, err := e.Discovery(reportCh, s.Nodes, s.Data)
		if err != nil {
			log.Printf("Discovery engine %s with data (%s) failed: %s\n",
				e.Name(), string(s.Data), err)
			continue
		} else if len(id) == 0 {
			log.Printf("Discovery engine %s should not return empty id\n", e.Name())
			continue
		}

		_, ok := scenarioMap[id]
		if ok {
			log.Printf("Scenraio %s had been discovery in engine %s\n",
				id, e.Name())
			continue
		}

		s.Name = id
		scenarioMap[id] = s
	}

	return cfg.Scenarios, nil
}

func initTriggers() error {
	cfg, err := conf.GetConf()
	if err != nil {
		return err
	}

	trigger.InitTriggers(cfg.Triggers)
	return nil
}

func runOpps(cmd *cobra.Command, args []string) error {
	ss, err := initOpps()
	if err != nil {
		return err
	}

	if err := initTriggers(); err != nil {
		log.Println("Init triggers with failed:", err)
		return err
	}

	go handleScenarioReport(ss)
	return nil
}

func handleScenarioReport(ss []*conf.Scenario) {
	if len(ss) == 0 {
		log.Println("This cluster do not have any scenario to wait report")
		return
	}

	reportStatus := make(map[string]conf.ReportStatus)
	for {
		select {
		case r := <-reportCh:
			s, ok := scenarioMap[r.Name]
			if !ok {
				log.Printf("Report name %s is not declear at any scenario\n", r.Name)
			}
			reportStatus[r.Name] = r.Status
			log.Printf("Scenario name %s had been %s status\n",
				r.Name, r.Status)
			trigger.Trigger(conf.CopyNodes(s.Nodes), r.Data, s.Trigger...)

			for i, s := range ss {
				st, ok := reportStatus[s.Name]
				if !ok || st != conf.ReportSucceededStatus {
					break
				} else if i == len(ss)-1 {
					reportDone <- true
					log.Println("All scenario had been succeed status, so finish opps")
					return
				}
			}
		}
	}
}
