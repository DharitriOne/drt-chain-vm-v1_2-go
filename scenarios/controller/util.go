package scencontroller

import (
	fr "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/fileresolver"
)

// NewDefaultFileResolver yields a new DefaultFileResolver instance.
// Reexported here to avoid having all external packages importing the parser.
// DefaultFileResolver is in parse for local tests only.
func NewDefaultFileResolver() *fr.DefaultFileResolver {
	return fr.NewDefaultFileResolver()
}
