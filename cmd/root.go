package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diegosantochi/wid/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wid",
	Short: "What Was I Doing",
	Long:  `WID helps you remember what you were doing by keeping a list of tasks with titles, descriptions, and statuses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := tui.New()
		if err != nil {
			return err
		}
		_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
