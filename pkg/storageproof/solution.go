// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import (
	"crypto/rand"
	"encoding/ascii85"

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

func NewSolution(challengeHash []byte, distance int, sk *mldsa87.PrivateKey) (*Solution, error) {
	pk := sk.Public().(*mldsa87.PublicKey)

	pkBytes, err := pk.MarshalBinary()
	if err != nil {
		return nil, err
	}

	sig, err := sk.Sign(rand.Reader, challengeHash, nil)
	if err != nil {
		return nil, err
	}

	hashDst := make([]byte, ascii85.MaxEncodedLen(len(challengeHash)))
	ascii85.Encode(hashDst, challengeHash)

	pkDst := make([]byte, ascii85.MaxEncodedLen(len(pkBytes)))
	ascii85.Encode(pkDst, pkBytes)

	sigDst := make([]byte, ascii85.MaxEncodedLen(len(sig)))
	ascii85.Encode(sigDst, sig)

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

	sigBytes := make([]byte, 4096) // Allocate a large enough buffer for the signature
	decoded, _, err = ascii85.Decode(sigBytes, []byte(s.Signature), true)
	if err != nil {
		return false, err
	}
	sigBytes = sigBytes[:decoded]

	return mldsa87.Verify(pk, hashBytes, sigBytes, nil), nil
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