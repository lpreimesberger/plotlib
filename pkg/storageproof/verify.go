// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import (
	"fmt"
	os"
	"time"

	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
	"golang.org/x/crypto/argon2"
)

func Verify(filePath string, verbose bool) error {
	// Open the plot file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the header
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

	// Read the key entries
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

	startTime := time.Now()

	// Verify each key
	for i, ke := range keyEntries {
		if verbose {
			// Calculate ETA
			elapsed := time.Since(startTime)
			progress := float64(i+1) / float64(header.NumKeys)
			eta := time.Duration(float64(elapsed) / progress * (1 - progress))
			fmt.Printf("Verifying key %d of %d (ETA: %s)\r", i+1, header.NumKeys, eta.Round(time.Second))
		}

		// Seek to the private key offset
		_, err = file.Seek(int64(ke.Offset), os.SEEK_SET)
		if err != nil {
			return err
		}

		// Read the private key
		skBytes := make([]byte, mldsa87.PrivateKeySize)
		_, err = file.Read(skBytes)
		if err != nil {
			return err
		}

		sk := &mldsa87.PrivateKey{}
		err = sk.UnmarshalBinary(skBytes)
		if err != nil {
			return err
		}

		// Derive the public key
		pk := sk.Public().(*mldsa87.PublicKey)

		// Generate the public key hash
		pkBytes, err := pk.MarshalBinary()
		if err != nil {
			return err
		}
		hash := argon2.IDKey(pkBytes, []byte("storageproof"), 1, 64*1024, 4, 32)

		// Compare the hashes
		if string(hash) != string(ke.Hash[:]) {
			return fmt.Errorf("key %d: hash mismatch", i)
		}
	}

	fmt.Println("\nVerification successful!")
	return nil
}