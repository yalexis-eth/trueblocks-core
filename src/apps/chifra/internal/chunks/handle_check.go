// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"fmt"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func (opts *ChunksOptions) HandleChunksCheck(blockNums []uint64) error {
	filenameChan := make(chan cache.IndexFileInfo)

	var nRoutines int = 1
	go cache.WalkCacheFolder(opts.Globals.Chain, cache.Index_Bloom, filenameChan)

	filenames := []string{}
	for result := range filenameChan {
		switch result.Type {
		case cache.Index_Bloom:
			hit := false
			for _, block := range blockNums {
				h := result.Range.BlockIntersects(block)
				hit = hit || h
				if hit {
					break
				}
			}
			if len(blockNums) == 0 || hit {
				filenames = append(filenames, result.Path)
			}
		case cache.None:
			nRoutines--
			if nRoutines == 0 {
				close(filenameChan)
			}
		}
	}

	sort.Slice(filenames, func(i, j int) bool {
		return filenames[i] < filenames[j]
	})

	allow_missing := config.GetBlockScrapeSettings(opts.Globals.Chain).Allow_missing

	nChecks := 0
	nChecksFailed := 0
	notARange := cache.FileRange{First: utils.NOPOS, Last: utils.NOPOS}
	if len(filenames) > 0 {
		prev := notARange
		for _, filename := range filenames {
			fR, _ := cache.RangeFromFilename(filename)
			if prev == notARange {
				prev = fR
			} else if prev != fR {
				nChecks++
				if !fR.Follows(prev, !allow_missing) {
					nChecksFailed++
					fmt.Println(fR, "does not sequentially follow", prev)
				}
			}
			prev = fR
		}
		fmt.Printf("Checked %d chunks, %d failed checks.\n", nChecks, nChecksFailed)
	}

	return nil
}

// TODO: BOGUS We don't check blooms
// TODO: BOGUS We don't check internal consistency of the data files
// TODO: BOGUS - is every address in the index, in the bloom?
// TODO: BOGUS - are there missing blocks (if allow_missing is off)
// TODO: BOGUS We don't check if Pinata has files that aren't needed