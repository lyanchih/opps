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
)

func init() {
	reportCh = make(chan conf.ScenarioReport, 10)
	scenarioMap = make(map[string]*conf.Scenario)
}

func initConfig() error {
	log.Println("Ready to parse config file", config)
	_, err := conf.ParseConf(config)
	return err
}

func initScenarios() error {
	cfg, err := conf.GetConf()
	if err != nil {
		return err
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

	return nil
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
	if err := initConfig(); err != nil {
		log.Println("Init config with failed:", err)
		return err
	}

	if err := initScenarios(); err != nil {
		log.Println("Init scenarios with failed:", err)
		return err
	}

	if err := initTriggers(); err != nil {
		log.Println("Init triggers with failed:", err)
		return err
	}

	go handleScenarioReport()
	return nil
}

func handleScenarioReport() {
	for {
		select {
		case r := <-reportCh:
			s, ok := scenarioMap[r.Name]
			if !ok {
				log.Printf("Report name %s is not declear at any scenario\n", r.Name)
			}
			log.Printf("Scenario name %s had been %s status\n",
				r.Name, r.Status)
			trigger.Trigger(s.Nodes, r.Data, s.Trigger...)
		}
	}
}
