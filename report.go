package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.adlinktech.com/lyan.hung/opps/report"
	"log"
)

func newReportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Report cluster nodes information",
		RunE:  runReport,
	}

	return cmd
}

func runReport(cmd *cobra.Command, args []string) error {
	ss, err := initOpps()
	if err != nil {
		return err
	}

	r, err := report.GetReporter("ansible")
	if err != nil {
		log.Println("Get reporter failed with:", err)
		return err
	}

	data, err := r.Report(ss)
	if err != nil {
		log.Printf("Report ansible format failed with: %s\n", err)
		return err
	}

	fmt.Println(string(data))
	return nil
}
