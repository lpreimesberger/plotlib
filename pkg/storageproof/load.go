// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
)

type PlotCollection struct {
	Plots map[string]*PlotInfo
}

type PlotInfo struct {
	*Header
	KeyEntries []KeyEntry
}

func LoadPlots(paths []string, verbose bool) (*PlotCollection, error) {
	pc := &PlotCollection{
		Plots: make(map[string]*PlotInfo),
	}

	for _, path := range paths {
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasPrefix(info.Name(), "sp") && strings.HasSuffix(info.Name(), ".plot") {
				if verbose {
					fmt.Printf("Loading plot file: %s\n", filePath)
				}

				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()

				headerBytes := make([]byte, 40)
				_, err = file.Read(headerBytes)
				if err != nil {
					return err
				}

				header := &Header{}
				err = header.UnmarshalBinary(headerBytes)
				if err != nil {
					return err
				}

				keyEntries := make([]KeyEntry, header.NumKeys)
				for i := uint32(0); i < header.NumKeys; i++ {
					keBytes := make([]byte, 40)
					_, err = file.Read(keBytes)
					if err != nil {
						return err
					}
					err = keyEntries[i].UnmarshalBinary(keBytes)
					if err != nil {
						return err
					}
				}

				pc.Plots[filePath] = &PlotInfo{
					Header:     header,
					KeyEntries: keyEntries,
				}
			}

			return nil
		})
	}

	return pc, nil
}

func (pc *PlotCollection) LookUp(challengeHash []byte) ([]byte, int, *mldsa87.PrivateKey, error) {
	var bestMatch []byte
	var bestDistance int = -1
	var bestPlotPath string
	var bestKeyEntry KeyEntry

	for plotPath, plotInfo := range pc.Plots {
		for _, keyEntry := range plotInfo.KeyEntries {
			distance := hammingDistance(challengeHash, keyEntry.Hash[:])
			if bestDistance == -1 || distance < bestDistance {
				bestDistance = distance
				bestMatch = keyEntry.Hash[:]
				bestPlotPath = plotPath
				bestKeyEntry = keyEntry
			}
		}
	}

	if bestDistance == -1 {
		return nil, -1, nil, nil // No plots loaded
	}

	// Now retrieve the private key
	file, err := os.Open(bestPlotPath)
	if err != nil {
		return nil, -1, nil, err
	}
	defer file.Close()

	_, err = file.Seek(int64(bestKeyEntry.Offset), 0)
	if err != nil {
		return nil, -1, nil, err
	}

	skBytes := make([]byte, mldsa87.PrivateKeySize)
	_, err = file.Read(skBytes)
	if err != nil {
		return nil, -1, nil, err
	}

	sk := &mldsa87.PrivateKey{}
	err = sk.UnmarshalBinary(skBytes)
	if err != nil {
		return nil, -1, nil, err
	}

	return bestMatch, bestDistance, sk, nil
}

func hammingDistance(a, b []byte) int {
	if len(a) != len(b) {
		return -1 // Or handle error appropriately
	}

	distance := 0
	for i := range a {
		xor := a[i] ^ b[i]
		for xor > 0 {
			distance++
			xor &= xor - 1
		}
	}
	return distance
}