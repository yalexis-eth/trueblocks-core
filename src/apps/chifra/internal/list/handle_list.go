// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package listPkg

import (
	"fmt"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/bykof/gostradamus"
)

func (opts *ListOptions) HandleListAppearances(monitorArray []monitor.Monitor) error {
	for _, mon := range monitorArray {
		count := mon.Count()
		apps := make([]index.AppearanceRecord, count, count)
		err := mon.ReadAppearances(&apps)
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			fmt.Println("No appearances found for", mon.GetAddrStr())
			return nil
		}

		sort.Slice(apps, func(i, j int) bool {
			si := uint64(apps[i].BlockNumber)
			si = (si << 32) + uint64(apps[i].TransactionId)
			sj := uint64(apps[j].BlockNumber)
			sj = (sj << 32) + uint64(apps[j].TransactionId)
			return si < sj
		})

		exportRange := cache.FileRange{First: opts.FirstBlock, Last: opts.LastBlock}
		results := make([]types.SimpleAppearance, 0, mon.Count())
		verboseResults := make([]types.VerboseAppearance, 0, mon.Count())
		for _, app := range apps {
			appRange := cache.FileRange{First: uint64(app.BlockNumber), Last: uint64(app.BlockNumber)}
			if appRange.Intersects(exportRange) {
				if opts.Globals.Verbose {
					ts, err := tslib.FromBnToTs(opts.Globals.Chain, uint64(app.BlockNumber))
					if err != nil {
						return err
					}
					s := types.VerboseAppearance{
						Address:          mon.GetAddrStr(),
						BlockNumber:      app.BlockNumber,
						TransactionIndex: app.TransactionId,
						Timestamp:        ts,
						Date:             gostradamus.FromUnixTimestamp(int64(ts)),
					}
					verboseResults = append(verboseResults, s)
				} else {
					s := types.SimpleAppearance{
						Address:          mon.GetAddrStr(),
						BlockNumber:      app.BlockNumber,
						TransactionIndex: app.TransactionId,
					}
					results = append(results, s)
				}
			}
		}

		// TODO: Fix export without arrays
		if opts.Globals.Verbose {
			err = globals.RenderSlice(&opts.Globals, verboseResults)
			if err != nil {
				return err
			}
		} else {
			err = globals.RenderSlice(&opts.Globals, results)
			if err != nil {
				return err
			}
		}
	}
	return nil
}