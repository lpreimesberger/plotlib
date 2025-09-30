// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import (
	"crypto/rand"
	"testing"

	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
)

func TestNewSolutionAndVerify(t *testing.T) {
	// Generate a new key pair
	_, sk, err := mldsa87.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create a challenge hash
	challengeHash := make([]byte, 32)
	_, err = rand.Read(challengeHash)
	if err != nil {
		t.Fatalf("Failed to create challenge hash: %v", err)
	}

	t.Logf("sk: %v", sk)

	// Create a new solution
	solution, err := NewSolution(challengeHash, 10, sk)
	if err != nil {
		t.Fatalf("Failed to create new solution: %v", err)
	}

	// Verify the solution
	valid, err := solution.Verify()
	if err != nil {
		t.Fatalf("Failed to verify solution: %v", err)
	}

	if !valid {
		t.Errorf("Expected solution to be valid, but it was not")
	}

	// Test with an invalid signature
	_, sk2, err := mldsa87.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	solution2, err := NewSolution(challengeHash, 10, sk2)
	if err != nil {
		t.Fatalf("Failed to create new solution: %v", err)
	}

	// The signature from solution2 should not be valid for solution1
	solution.Signature = solution2.Signature

	t.Logf("Verifying with incorrect signature")
	valid, err = solution.Verify()
	if err != nil {
		t.Fatalf("Failed to verify solution: %v", err)
	}

	if valid {
		t.Errorf("Expected solution to be invalid, but it was valid")
	}
}
