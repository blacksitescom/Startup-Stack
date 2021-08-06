/*
Copyright © 2021 Ian James Gordon
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// scloudCmd represents the scloud command
var scloudCmd = &cobra.Command{
	Use:   "scloud",
	Short: "scloud creates and manages secure cloud environments",
	Long: `scloud raises a custom cloud environment with a zero-trust
architecture, based on your inputs, and then manages the development 
and operations lifecycle of the environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scloud called")
	},
}

func init() {
	rootCmd.AddCommand(scloudCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scloudCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scloudCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
