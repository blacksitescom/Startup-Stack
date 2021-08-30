/*
Copyright © 2021 Ian James Gordon
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/gordonianj/blacksite/images"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"l"},
	Short:   "Describe the existing images",
	RunE: func(cmd *cobra.Command, args []string) error {
		awsRegion, err := cmd.Flags().GetString("aws-region")
		if err != nil {
			return err
		}

		return describeAction(os.Stdout, awsRegion, args)
	},
}

func describeAction(out io.Writer, awsRegion string, args []string) error {
	il := &images.Images{}

	if err := il.Describe(awsRegion); err != nil {
		return err
	}

	for _, i := range il.Images {
		if _, err := fmt.Fprintln(out, i); err != nil {
			return err
		}
	}

	return nil

}

func init() {
	imagesCmd.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	describeCmd.Flags().StringP("aws-region", "r", "us-west-2", "AWS EC2 region")
}
