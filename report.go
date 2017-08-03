package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.adlinktech.com/lyan.hung/opps/report"
	"io/ioutil"
	"log"
)

var reportOutputFile string

func newReportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Report cluster nodes information",
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&reportOutputFile, "report-output", "", "", "Report cluster information file")

	err := report.AddSubcommands(cmd, runReport)
	if err != nil {
		log.Println("Add report subcommands failed:", err)
		return nil
	}
	return cmd
}

func runReport(cmd *cobra.Command, args []string) error {
	ss, err := initOpps()
	if err != nil {
		return err
	}

	r, err := report.GetReporter(cmd.Use)
	if err != nil {
		log.Println("Get reporter failed with:", err)
		return err
	}

	data, err := r.Report(ss)
	if err != nil {
		log.Printf("Report ansible format failed with: %s\n", err)
		return err
	}

	if len(reportOutputFile) != 0 {
		log.Println("Write report output to file", reportOutputFile)
		if err := ioutil.WriteFile(reportOutputFile, data, 0644); err != nil {
			log.Println("Write report output failed:", err)
			return err
		}
	}
	fmt.Println(string(data))
	return nil
}
