// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify [solution]",
	Short: "Verifies a storage proof solution.",
	Long:  `Verifies a storage proof solution provided as a JSON string.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		solutionJSON := args[0]

		var solution storageproof.Solution
		err := json.Unmarshal([]byte(solutionJSON), &solution)
		if err != nil {
			fmt.Printf("Error unmarshalling solution: %s\n", err)
			return
		}

		valid, err := solution.Verify()
		if err != nil {
			fmt.Printf("Error verifying solution: %s\n", err)
			return
		}

		if valid {
			fmt.Println("Solution is valid")
		} else {
			fmt.Println("Solution is invalid")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
