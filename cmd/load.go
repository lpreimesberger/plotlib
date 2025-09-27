// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package cmd

import (
	"fmt"
	"strings"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [paths]",
	Short: "Loads plot files from a comma-delimited list of paths.",
	Long: `Loads plot files from a comma-delimited list of paths. 
If a path is a directory, it will be searched recursively for plot files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths := strings.Split(args[0], ",")

		pc, err := storageproof.LoadPlots(paths, verbose)
		if err != nil {
			fmt.Printf("Error loading plots: %s\n", err)
			return
		}

		if verbose {
			fmt.Printf("Loaded %d plot files.\n", len(pc.Plots))
			var totalKeys uint32
			for _, plot := range pc.Plots {
				totalKeys += plot.Header.NumKeys
			}
			fmt.Printf("Total solutions: %d\n", totalKeys)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
