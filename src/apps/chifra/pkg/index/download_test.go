// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package index

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/manifest"
)

func Test_exclude(t *testing.T) {
	fileNames := map[string]bool{
		"013337527-013340418": true,
		"013340419-013343305": true,
		"013346064-013348861": true,
		"013348862-013351760": true,
	}

	pins := []manifest.ChunkRecord{
		{
			Range: "013337527-013340418",
		},
		{
			Range: "013340419-013343305",
		},
		{
			Range: "013346064-013348861",
		},
		{
			Range: "013348862-013351760",
		},
		{
			Range: "013387069-013389874",
		},
		{
			Range: "013389875-013392800",
		},
	}

	result := exclude(fileNames, pins)

	if len(result) != 2 {
		t.Errorf("Wrong length: %d", len(result))
	}

	if result[0].Range != "013387069-013389874" &&
		result[1].Range != "013389875-013392800" {
		t.Errorf("Bad values: '%s' and '%s'", result[0].Range, result[1].Range)
	}
}