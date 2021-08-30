/*
Copyright © 2021 Ian James Gordon
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"io"
	"os"

	"github.com/gordonianj/blacksite/images"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <image name>",
	Aliases:      []string{"a"},
	Short:        "Add an image",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		awsRegion, err := cmd.Flags().GetString("aws-region")
		if err != nil {
			return err
		}

		return addAction(os.Stdout, awsRegion, args)
	},
}

func addAction(out io.Writer, awsRegion string, args []string) error {
	il := &images.Images{}

	for _, i := range args {
		if err := il.Add(i); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	imagesCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringP("aws-region", "r", "us-west-2", "AWS region")
}
