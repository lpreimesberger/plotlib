// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const libVersion = "0.0.1"

func Plot(destDir string, kValue uint32, verbose bool) error {
	numKeys := kValue * 1000

	// Generate a new UUID for the plot file
	guid := uuid.New()
	fileName := fmt.Sprintf("sp%d%s.plot", Version, guid.String())
	filePath := fmt.Sprintf("%s/%s", destDir, fileName)

	// Create the plot file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Write a placeholder for the header
	h := &Header{
		Version: Version,
		NumKeys: numKeys,
	}
	copy(h.LibVersion[:], libVersion)

	headerBytes, err := h.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = file.Write(headerBytes)
	if err != nil {
		return err
	}

	// Write a placeholder for the key entries
	keyEntries := make([]KeyEntry, numKeys)
	keyEntriesBytes := make([]byte, 40*numKeys)
	_, err = file.Write(keyEntriesBytes)
	if err != nil {
		return err
	}

	startTime := time.Now()

	// Generate keys and write them to the file
	for i := uint32(0); i < numKeys; i++ {
		if verbose {
			// Calculate ETA
			elapsed := time.Since(startTime)
			progress := float64(i+1) / float64(numKeys)
			eta := time.Duration(float64(elapsed) / progress * (1 - progress))
			fmt.Printf("Plotting key %d of %d (ETA: %s)\r", i+1, numKeys, eta.Round(time.Second))
		}

		// Generate a new key pair
		pk, sk, err := mldsa87.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		// Get the current offset
		offset, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		// Write the private key to the file
		skBytes, err := sk.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = file.Write(skBytes)
		if err != nil {
			return err
		}

		// Generate the public key hash
		pkBytes, err := pk.MarshalBinary()
		if err != nil {
			return err
		}
		hash := argon2.IDKey(pkBytes, []byte("storageproof"), 1, 64*1024, 4, 32)

		// Update the key entry
		keyEntries[i].Offset = uint64(offset)
		copy(keyEntries[i].Hash[:], hash)
	}

	// Go back to the beginning of the file and write the final header and key entries
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = file.Write(headerBytes)
	if err != nil {
		return err
	}

	for _, ke := range keyEntries {
		keBytes, err := ke.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = file.Write(keBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// HammingDistance calculates the hamming distance between two byte slices
func HammingDistance(a, b []byte) int {
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
