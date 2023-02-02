// Copyright 2023 Scott M. Long
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
