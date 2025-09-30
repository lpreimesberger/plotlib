// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import (
	"crypto"
	"crypto/rand"
	"encoding/ascii85"
	"errors"

	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
)

// Solution represents a solution to a storage proof challenge
// All binary data is encoded as base85 strings for long-term storage

type Solution struct {
	Hash      string `json:"hash"`
	Distance  int    `json:"distance"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

const shakeOutputLen = 64 // 512 bits for 256-bit security
// Shake256SignerOpts implements crypto.SignerOpts for SHAKE256.
type Shake256SignerOpts struct {
	// The desired output length in bytes.
	// For 256-bit security, use at least 64 bytes (512 bits) of output.
	OutputLen int
}

// HashFunc returns 0 because SHAKE256 is an XOF, not a fixed hash.
// A value of 0 signals to the signer that the hash will be provided as a digest
// with a flexible length, as is the case for SHAKE.
func (o Shake256SignerOpts) HashFunc() crypto.Hash {
	return 0
}

func NewSolution(challengeHash []byte, distance int, sk *mldsa87.PrivateKey) (*Solution, error) {
	pk := sk.Public().(*mldsa87.PublicKey)

	pkBytes, err := pk.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if len(challengeHash) != 32 {
		return nil, errors.New("challenge hash length must be 32 bytes")
	}
	if sk == nil {
		return nil, errors.New("sk cannot be nil")
	}
	opts := &Shake256SignerOpts{OutputLen: shakeOutputLen}
	sig, err := sk.Sign(rand.Reader, challengeHash, opts)
	if err != nil {
		return nil, err
	}

	hashDst := make([]byte, ascii85.MaxEncodedLen(len(challengeHash)))
	encoded := ascii85.Encode(hashDst, challengeHash)
	hashDst = hashDst[:encoded]

	pkDst := make([]byte, ascii85.MaxEncodedLen(len(pkBytes)))
	encoded = ascii85.Encode(pkDst, pkBytes)
	pkDst = pkDst[:encoded]

	sigDst := make([]byte, ascii85.MaxEncodedLen(len(sig)))
	encoded = ascii85.Encode(sigDst, sig)
	sigDst = sigDst[:encoded]

	return &Solution{
		Hash:      string(hashDst),
		Distance:  distance,
		PublicKey: string(pkDst),
		Signature: string(sigDst),
	}, nil
}

func (s *Solution) Verify() (bool, error) {
	hashBytes := make([]byte, 32)
	decoded, _, err := ascii85.Decode(hashBytes, []byte(s.Hash), true)
	if err != nil {
		return false, err
	}
	hashBytes = hashBytes[:decoded]

	pkBytes := make([]byte, mldsa87.PublicKeySize)
	decoded, _, err = ascii85.Decode(pkBytes, []byte(s.PublicKey), true)
	if err != nil {
		return false, err
	}
	pkBytes = pkBytes[:decoded]

	pk := &mldsa87.PublicKey{}
	err = pk.UnmarshalBinary(pkBytes)
	if err != nil {
		return false, err
	}

	sigBytes := make([]byte, 1024*8) // Allocate a large enough buffer for the signature
	decoded, _, err = ascii85.Decode(sigBytes, []byte(s.Signature), true)
	if err != nil {
		return false, err
	}
	sigBytes = sigBytes[:decoded]

	return mldsa87.Verify(pk, hashBytes, nil, sigBytes), nil
}

// BestMatch returns the best solution from a slice of solutions
func BestMatch(solutions []*Solution) *Solution {
	if len(solutions) == 0 {
		return nil
	}

	bestSolution := solutions[0]
	for _, solution := range solutions {
		if solution.Distance < bestSolution.Distance {
			bestSolution = solution
		}
	}

	return bestSolution
}
