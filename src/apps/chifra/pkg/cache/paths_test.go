// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package cache

import (
	"path"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
)

func TestCacheLayout(t *testing.T) {
	indexPath := config.GetPathToIndex(GetTestChain())
	cachePath := config.GetPathToCache(GetTestChain())

	// TODO: turn these back on
	tests := []struct {
		on        bool
		name      string
		cacheType CacheType
		param     string
		expected  CachePath
		path      string
		wantErr   bool
	}{
		{
			on:    true,
			name:  "index chunk path",
			param: "0010000000-0010200000",
			expected: CachePath{
				Type:      Index_Final,
				RootPath:  indexPath,
				Subdir:    "finalized/",
				Extension: ".bin",
			},
			path:    "finalized/0010000000-0010200000.bin",
			wantErr: false,
		},
		{
			on:    true,
			name:  "Bloom filter path",
			param: "0010000000-0010200000",
			expected: CachePath{
				Type:      Index_Bloom,
				RootPath:  indexPath,
				Subdir:    "blooms/",
				Extension: ".bloom",
			},
			path:    "blooms/0010000000-0010200000.bloom",
			wantErr: false,
		},
		{
			on:    false,
			name:  "Block cache path",
			param: "001001001",
			expected: CachePath{
				Type:      Cache_Block,
				RootPath:  cachePath,
				Subdir:    "blocks/",
				Extension: ".bin",
			},
			path:    "blocks/00/10/01/001001001.bin",
			wantErr: false,
		},
		{
			on:    false,
			name:  "Transaction cache path",
			param: "1001001.20",
			expected: CachePath{
				Type:      Cache_Tx,
				RootPath:  cachePath,
				Subdir:    "txs/",
				Extension: ".bin",
			},
			path:    "txs/00/10/01/001001001-00020.bin",
			wantErr: false,
		},
		// TraceFn:      $CACHE_PATH/traces/00/10/01/001001001-00020-10.bin
		// NeighborFn:   $CACHE_PATH/neighbors/00/10/01/001001001-00020.bin
		// ReconFn:      $CACHE_PATH/recons/c011/a724/00e58ecd99ee497cf89e3775d4bd732f/000000012.00013.bin
	}

	for _, tt := range tests {
		if !tt.on {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			cachePath := NewCachePath(GetTestChain(), tt.expected.Type)
			if cachePath.Extension != tt.expected.Extension {
				t.Error("Wrong extension", cachePath.Extension)
			}
			if cachePath.Subdir != tt.expected.Subdir {
				t.Error("Wrong subdir", cachePath.Subdir)
			}
			p := cachePath.GetFullPath(tt.param)
			if p != path.Join(tt.expected.RootPath, tt.path) {
				t.Error("Wrong full path", p)
			}
		})
	}
}

// GetTestChain is duplicated in multiple packages to avoid dependancies. See
// https://stackoverflow.com/questions/49789055/
func GetTestChain() string {
	return "mainnet"
}
