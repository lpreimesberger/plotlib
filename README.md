# Plotlib

Plotlib is a Go library and command-line tool for creating and verifying proof-of-storage plot files for the Shadowy Apparatus.

## Installation

To install the command-line tool:

```bash
go install github.com/lpreimesberger/plotlib
```

To use the library in your own project:

```bash
go get github.com/lpreimesberger/plotlib
```

## CLI Usage

### Global Flags

*   `-v`, `--verbose`: Enable verbose output.

### `plot`

Generates a new plot file.

```bash
plotlib plot [kValue] [destDir]
```

*   `kValue`: The number of keys to generate in thousands.
*   `destDir`: The destination directory for the plot file.

### `verify`

Verifies a storage proof solution.

```bash
plotlib verify [solution]
```

*   `solution`: A JSON string representing the solution.

### `load`

Loads plot files from a comma-delimited list of paths.

```bash
plotlib load [paths]
```

*   `paths`: A comma-delimited list of directories or plot files.

### `lookup`

Looks up a hash in the plot files.

```bash
plotlib lookup [paths] [hash]
```

*   `paths`: A comma-delimited list of directories or plot files.
*   `hash` (optional): The hash to look up. If not provided, a test suite is run.

### `benchmarklookup`

Benchmarks the lookup function.

```bash
plotlib benchmarklookup [paths]
```

*   `paths`: A comma-delimited list of directories or plot files.

## Library Usage

The following is a brief example of how to use the `plotlib` library.

```go
package main

import (
	"fmt"
	"log"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
)

func main() {
	// 1. Create a new plot file
	err := storageproof.Plot("./plots", 1, true)
	if err != nil {
		log.Fatalf("Failed to plot: %v", err)
	}

	// 2. Load the plot files
	pc, err := storageproof.LoadPlots([]string{"./plots"}, true)
	if err != nil {
		log.Fatalf("Failed to load plots: %v", err)
	}

	// 3. Define a challenge hash
	challengeHash := []byte("some 32-byte hash")

	// 4. Look up the challenge hash in the plot files
	solution, err := pc.LookUp(challengeHash)
	if err != nil {
		log.Fatalf("Failed to lookup: %v", err)
	}

	fmt.Printf("Best match: %s
", solution.Hash)
	fmt.Printf("Distance: %d
", solution.Distance)

	// 5. Verify the solution
	valid, err := solution.Verify()
	if err != nil {
		log.Fatalf("Failed to verify solution: %v", err)
	}

	if valid {
		fmt.Println("Solution is valid")
	} else {
		fmt.Println("Solution is invalid")
	}
}
```

## Plot File Format

The plot file has the following structure:

1.  **Header:**
    *   `Version` (uint32)
    *   `NumKeys` (uint32)
    *   `LibVersion` ([32]byte)
2.  **Key Entries:** A list of `KeyEntry` structs:
    *   `Offset` (uint64)
    *   `Hash` ([32]byte) - The Argon2 hash of the corresponding public key.
3.  **Private Keys:** The raw private keys.

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for details.