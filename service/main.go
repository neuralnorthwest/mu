package service

import "github.com/spf13/cobra"

func (s *Service) Main() error {
	cmd := &cobra.Command{
		Use:   s.name,
		Short: s.name,
		Long:  s.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.run()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	// Default the mock flag to the current value of s.mockMode. This prevents
	// the flag's default setting from overriding the value set by WithMockMode.
	cmd.PersistentFlags().BoolVar(&s.mockMode, "mock", s.mockMode, "enable mock mode")
	return cmd.Execute()
}
