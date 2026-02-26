package cmd

import (
	"fmt"

	"github.com/diegosantochi/wid/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .wid.yaml file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.Init(); err != nil {
			return err
		}
		fmt.Println("Initialized .wid.yaml in current directory.")
		return nil
	},
}
