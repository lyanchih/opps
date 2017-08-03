package report

import (
	"github.com/spf13/cobra"
)

type commander interface {
	cmd(runE func(*cobra.Command, []string) error) *cobra.Command
}

func AddSubcommands(cmd *cobra.Command, runE func(*cobra.Command, []string) error) error {
	for _, r := range reporters {
		c, ok := r.(commander)
		if !ok {
			continue
		}

		cmd.AddCommand(c.cmd(runE))
	}

	return nil
}
