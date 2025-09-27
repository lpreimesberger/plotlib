// SPDX-License-Identifier: Apache-2.0
// Copyright 2025 Caprica LLC

package storageproof

import "encoding/binary"

const Version = 1

// Header defines the structure of the plot file header.
// The header will be followed by the key data.
// The key data will be a sequence of private keys.
// The offsets in the header will point to the start of each private key.
// The hashes in the header will be the Argon2 hash of the corresponding public key.

type Header struct {
	Version       uint32
	NumKeys       uint32
	LibVersion    [32]byte // Fixed-size array for a 32-character string
}

// KeyEntry defines the structure of the key lookup table in the header.

type KeyEntry struct {
	Offset uint64
	Hash   [32]byte // Assuming a 32-byte hash output
}

func (h *Header) MarshalBinary() ([]byte, error) {
	b := make([]byte, 40)
	binary.LittleEndian.PutUint32(b[0:4], h.Version)
	binary.LittleEndian.PutUint32(b[4:8], h.NumKeys)
	copy(b[8:40], h.LibVersion[:])
	return b, nil
}

func (h *Header) UnmarshalBinary(data []byte) error {
	h.Version = binary.LittleEndian.Uint32(data[0:4])
	h.NumKeys = binary.LittleEndian.Uint32(data[4:8])
	copy(h.LibVersion[:], data[8:40])
	return nil
}

func (ke *KeyEntry) MarshalBinary() ([]byte, error) {
	b := make([]byte, 40)
	binary.LittleEndian.PutUint64(b[0:8], ke.Offset)
	copy(b[8:40], ke.Hash[:])
	return b, nil
}

func (ke *KeyEntry) UnmarshalBinary(data []byte) error {
	ke.Offset = binary.LittleEndian.Uint64(data[0:8])
	copy(ke.Hash[:], data[8:40])
	return nil
}