# Plotlib

Plotlib is a Go library and command-line tool for creating and verifying proof-of-storage plot files.

## Installation

```bash
go get github.com/lpreimesberger/plotlib
```

## CLI Usage

### `plot`

Generates a new plot file.

```bash
plotlib plot [kValue] [destDir]
```

*   `kValue`: The number of keys to generate in thousands.
*   `destDir`: The destination directory for the plot file.

### `verify`

Verifies a plot file.

```bash
plotlib verify [filePath]
```

*   `filePath`: The path to the plot file.

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

### `Plot`

```go
package main

import (
	"github.com/lpreimesberger/plotlib/pkg/storageproof"
)

func main() {
	err := storageproof.Plot("./plots", 1, true)
	if err != nil {
		panic(err)
	}
}
```

### `Verify`

```go
package main

import (
	"github.com/lpreimesberger/plotlib/pkg/storageproof"
)

func main() {
	err := storageproof.Verify("./plots/sp1...plot", true)
	if err != nil {
		panic(err)
	}
}
```

### `LoadPlots`

```go
package main

import (
	"fmt"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
)

func main() {
	pc, err := storageproof.LoadPlots([]string{"./plots"}, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d plot files.\n", len(pc.Plots))
}
```

### `LookUp`

```go
package main

import (
	"fmt"

	"github.com/lpreimesberger/plotlib/pkg/storageproof"
)

func main() {
	pc, err := storageproof.LoadPlots([]string{"./plots"}, false)
	if err != nil {
		panic(err)
	}

	challengeHash := []byte("some 32-byte hash")
	match, distance, sk, err := pc.LookUp(challengeHash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Best match: %x\n", match)
	fmt.Printf("Distance: %d\n", distance)
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
    *   `Hash` ([32]byte)
3.  **Private Keys:** The raw private keys.

