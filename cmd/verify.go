// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package cmd

import (
	"fmt"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify [filePath]",
	Short: "Verifies a plot file.",
	Long:  `Verifies a plot file by checking the header and all the keys.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		err := storageproof.Verify(filePath, verbose)
		if err != nil {
			fmt.Printf("Error verifying: %s\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
