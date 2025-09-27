// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package cmd

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
	"github.com/spf13/cobra"
)

// benchmarklookupCmd represents the benchmarklookup command
var benchmarklookupCmd = &cobra.Command{
	Use:   "benchmarklookup [paths]",
	Short: "Benchmarks the lookup function.",
	Long: `Benchmarks the lookup function by generating 1024 random hashes 
and looking them up in the plot files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paths := strings.Split(args[0], ",")

		pc, err := storageproof.LoadPlots(paths, false) // Don't need verbose output for loading
		if err != nil {
			fmt.Printf("Error loading plots: %s\n", err)
			return
		}

		if len(pc.Plots) == 0 {
			fmt.Println("No plot files found.")
			return
		}

		fmt.Println("Benchmarking lookup function...")

		const numLookups = 1024
		randomHashes := make([][]byte, numLookups)
		for i := 0; i < numLookups; i++ {
			randomHashes[i] = make([]byte, 32)
			rand.Read(randomHashes[i])
		}

		startTime := time.Now()

		for i := 0; i < numLookups; i++ {
			_, _, _, err := pc.LookUp(randomHashes[i])
			if err != nil {
				fmt.Printf("Error looking up hash: %s\n", err)
				// We continue the benchmark even if one lookup fails
			}
		}

		totalTime := time.Since(startTime)
		avgTime := totalTime / numLookups

		fmt.Printf("\n--- Benchmark Results ---\n")
		fmt.Printf("Total lookups: %d\n", numLookups)
		fmt.Printf("Total time: %s\n", totalTime)
		fmt.Printf("Average lookup time: %s\n", avgTime)
	},
}

func init() {
	rootCmd.AddCommand(benchmarklookupCmd)
}
