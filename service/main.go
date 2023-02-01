package service

import "github.com/spf13/cobra"

// MainCommand returns the main cobra.Command for the service. This allows you
// to customize the command (perhaps adding flags) before invoking it with
// Execute.
func (s *Service) MainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   s.name,
		Short: s.name,
		Long:  s.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	// Default the mock flag to the current value of s.mockMode. This prevents
	// the flag's default setting from overriding the value set by WithMockMode.
	cmd.PersistentFlags().BoolVar(&s.mockMode, "mock", s.mockMode, "enable mock mode")
	return cmd
}

// Main invokes the main cobra.Command for the service. If you need to customize
// the command, see MainCommand.
func (s *Service) Main() error {
	return s.MainCommand().Execute()
}
