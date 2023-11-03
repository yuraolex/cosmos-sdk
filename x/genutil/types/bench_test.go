package types

import (
	"path/filepath"
	"testing"
)

var sink any = nil

func BenchmarkAppGenesisParsing(b *testing.B) {
	b.ReportAllocs()

	shortPaths := []string{
		"app_genesis.json",
		"big_app_genesis.json",
		"cmt_genesis.json",
		"big_cmt_genesis.json",
	}

	paths := make([]string, 0, len(shortPaths))
	for _, p := range shortPaths {
		paths = append(paths, filepath.Join("testdata", p))
	}

	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			sink, err = AppGenesisFromFile(path)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	if sink == nil {
		b.Fatal("Benchmark did not run!")
	}
}
