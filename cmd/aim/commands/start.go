package commands

import (
	aimos "github.com/BlockPILabs/aa-scan/internal/os"
	"github.com/spf13/cobra"
)

// AddFlags exposes some common configuration options on the command-line
// These are exposed for convenience of commands embedding
func AddFlags(cmd *cobra.Command) {

	// rpc flags
	cmd.Flags().String("api.laddr", config.Api.ListenAddress, "api listen address. Port required")
	cmd.Flags().Bool("api.unsafe", config.Api.Unsafe, "enabled unsafe api methods")
	cmd.Flags().String("api.pprof_laddr", config.Api.PprofListenAddress, "pprof listen address (https://golang.org/pkg/net/http/pprof)")

}

// NewStartCmd returns the command that allows the CLI to start a node.
func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"node", "run"},
		Short:   "Run the aim api",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("start api")

			// Stop upon receiving SIGTERM or CTRL-C.
			aimos.TrapSignal(logger, func() {
				logger.Info("end api")
			})

			// Run forever.
			select {}
		},
	}

	AddFlags(cmd)
	return cmd
}
