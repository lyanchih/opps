package main

import (
	"github.com/spf13/cobra"
)

var config string

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opps",
		Short: "OS provision progress supervisor",
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&config, "conf", "", "etc/opps.json", "opps config")

	cmd.AddCommand(newServeCommand())
	cmd.AddCommand(newReportCommand())
	return cmd
}

func waitDone(cmd *cobra.Command) {
	if cmd.Name() != serveCmdUse {
		closeServe()
	}

	select {
	case <-serverDone:
	}
}
