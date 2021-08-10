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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List the existing images",
	RunE: func(cmd *cobra.Command, args []string) error {
		imagesFile, err := cmd.Flags().GetString("images-file")
		if err != nil {
			return err
		}

		return listAction(os.Stdout, imagesFile, args)
	},
}

func listAction(out io.Writer, imagesFile string, args []string) error {
	il := &images.ImagesList{}

	if err := il.Load(imagesFile); err != nil {
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
	imagesCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
