// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package cmd

import (
	"fmt"
	"strconv"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
	"github.com/spf13/cobra"
)

// plotCmd represents the plot command
var plotCmd = &cobra.Command{
	Use:   "plot [kValue] [destDir]",
	Short: "Generates a new plot file.",
	Long: `Generates a new plot file with a given K value.
The K value represents the number of keys to generate in thousands.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kValue, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid K value")
			return
		}

		destDir := args[1]

		err = storageproof.Plot(destDir, uint32(kValue), verbose)
		if err != nil {
			fmt.Printf("Error plotting: %s\n", err)
			return
		}

		fmt.Println("Plot file generated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(plotCmd)
}
