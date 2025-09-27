// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup [paths] [hash]",
	Short: "Looks up a hash in the plot files.",
	Long: `Looks up a hash in the plot files. 
If a hash is not provided, it will run a test suite.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a comma-delimited list of paths.")
			return
		}
		paths := strings.Split(args[0], ",")

		pc, err := storageproof.LoadPlots(paths, true)
		if err != nil {
			fmt.Printf("Error loading plots: %s\n", err)
			return
		}

		if len(pc.Plots) == 0 {
			fmt.Println("No plot files found.")
			return
		}

		if len(args) == 2 {
			// Lookup a specific hash
			hash, err := hex.DecodeString(args[1])
			if err != nil {
				fmt.Printf("Invalid hash: %s\n", err)
				return
			}

			match, distance, _, err := pc.LookUp(hash)
			if err != nil {
				fmt.Printf("Error looking up hash: %s\n", err)
				return
			}

			fmt.Printf("Best match: %x\n", match)
			fmt.Printf("Distance: %d\n", distance)
		} else {
			// Run test suite
			fmt.Println("Running test suite...")

			// Positive case
			fmt.Println("\n--- Positive Case ---")
			var knownHash []byte
			for _, plot := range pc.Plots {
				knownHash = plot.KeyEntries[0].Hash[:]
				break
			}
			match, distance, _, err := pc.LookUp(knownHash)
			if err != nil {
				fmt.Printf("Error looking up hash: %s\n", err)
			} else {
				fmt.Printf("Looking up: %x\n", knownHash)
				fmt.Printf("Best match: %x\n", match)
				fmt.Printf("Distance: %d\n", distance)
			}

			// Near miss case
			fmt.Println("\n--- Near Miss Case ---")
			nearMissHash := make([]byte, len(knownHash))
			copy(nearMissHash, knownHash)
			nearMissHash[0] ^= 0x01 // Flip a bit
			match, distance, _, err = pc.LookUp(nearMissHash)
			if err != nil {
				fmt.Printf("Error looking up hash: %s\n", err)
			} else {
				fmt.Printf("Looking up: %x\n", nearMissHash)
				fmt.Printf("Best match: %x\n", match)
				fmt.Printf("Distance: %d\n", distance)
			}

			// Complete miss case
			fmt.Println("\n--- Complete Miss Case ---")
			randomHash := make([]byte, 32)
			rand.Read(randomHash)
			match, distance, _, err = pc.LookUp(randomHash)
			if err != nil {
				fmt.Printf("Error looking up hash: %s\n", err)
			} else {
				fmt.Printf("Looking up: %x\n", randomHash)
				fmt.Printf("Best match: %x\n", match)
				fmt.Printf("Distance: %d\n", distance)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
