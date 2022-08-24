// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package scrapeCfg

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type ScrapeSettings struct {
	Apps_per_chunk uint64 `json:"appsPerChunk"` // The number of appearances to build into a chunk before consolidating it
	Snap_to_grid   uint64 `json:"snapToGrid"`   // An override to apps_per_chunk to snap-to-grid at every modulo of this value, this allows easier corrections to the index
	First_snap     uint64 `json:"firstSnap"`    // The first block at which snap_to_grid is enabled
	Unripe_dist    uint64 `json:"unripeDist"`   // The distance (in blocks) from the front of the chain under which (inclusive) a block is considered unripe
	Channel_count  uint64 `json:"blockChanCnt"` // Number of concurrent block processing channels
	Allow_missing  bool   `json:"allowMissing"` // Do not report errors for blockchain that contain blocks with zero addresses
	// EXISTING_CODE
	// EXISTING_CODE
}

var defaultSettings = ScrapeSettings{
	Apps_per_chunk: 200000,
	Snap_to_grid:   100000,
	First_snap:     0,
	Unripe_dist:    28,
	Channel_count:  20,
	Allow_missing:  false,
	// EXISTING_CODE
	// EXISTING_CODE
}

var Unset = ScrapeSettings{
	Apps_per_chunk: utils.NOPOS,
	Snap_to_grid:   utils.NOPOS,
	First_snap:     utils.NOPOS,
	Unripe_dist:    utils.NOPOS,
	Channel_count:  utils.NOPOS,
	Allow_missing:  false,
	// EXISTING_CODE
	// EXISTING_CODE
}

func (s *ScrapeSettings) isDefault(chain, fldName string) bool {
	def := GetDefault(chain)
	switch fldName {
	case "Apps_per_chunk":
		return s.Apps_per_chunk == def.Apps_per_chunk
	case "Snap_to_grid":
		return s.Snap_to_grid == def.Snap_to_grid
	case "First_snap":
		return s.First_snap == def.First_snap
	case "Unripe_dist":
		return s.Unripe_dist == def.Unripe_dist
	case "Channel_count":
		return s.Channel_count == def.Channel_count
	case "Allow_missing":
		return s.Allow_missing == def.Allow_missing
	}

	// EXISTING_CODE
	// EXISTING_CODE

	return false
}

func (s *ScrapeSettings) TestLog(chain string, test bool) {
	logger.TestLog(!s.isDefault(chain, "Apps_per_chunk"), "Apps_per_chunk: ", s.Apps_per_chunk)
	logger.TestLog(!s.isDefault(chain, "Snap_to_grid"), "Snap_to_grid: ", s.Snap_to_grid)
	logger.TestLog(!s.isDefault(chain, "First_snap"), "First_snap: ", s.First_snap)
	logger.TestLog(!s.isDefault(chain, "Unripe_dist"), "Unripe_dist: ", s.Unripe_dist)
	logger.TestLog(!s.isDefault(chain, "Channel_count"), "Channel_count: ", s.Channel_count)
	logger.TestLog(!s.isDefault(chain, "Allow_missing"), "Allow_missing: ", s.Allow_missing)
	// EXISTING_CODE
	// EXISTING_CODE
}

func GetDefault(chain string) ScrapeSettings {
	ret := defaultSettings
	// EXISTING_CODE
	if chain == "mainnet" {
		ret.Apps_per_chunk = 2000000
		ret.First_snap = 2300000
	}
	// EXISTING_CODE
	return ret
}

const configFilename = "blockScrape.toml"

// GetSettings retrieves scrape config from (in order) default, config, environment, optionally provided cmdLine
func GetSettings(chain string, cmdLine *ScrapeSettings) (ScrapeSettings, error) {
	type TomlFile struct {
		Settings ScrapeSettings
	}

	// Start with the defalt values...
	ret := GetDefault(chain)

	tt := reflect.TypeOf(defaultSettings)
	fieldList, _, _ := utils.GetFields(&tt, "txt", true)

	configFn := filepath.Join(config.GetPathToChainConfig(chain), configFilename)
	if file.FileExists(configFn) {
		var t TomlFile
		// ...pick up values from toml file...
		if _, err := toml.Decode(utils.AsciiFileToString(configFn), &t); err != nil {
			return ScrapeSettings{}, err
		}
		ret.overlay(chain, t.Settings)
	}

	// ...check the environment...
	for _, field := range fieldList {
		envKey := toEnvStr(field)
		envValue := os.Getenv(envKey)
		if envValue != "" {
			fName := utils.MakeFirstUpperCase(field)
			fld := reflect.ValueOf(&ret).Elem().FieldByName(fName)
			if fld.Kind() == reflect.String {
				fld.SetString(envValue)
			} else if fld.Kind() == reflect.Bool {
				if envValue == "true" {
					fld.SetBool(true)
				}
			} else {
				if v, err := strconv.ParseUint(envValue, 10, 64); err == nil {
					fld.SetUint(v)
				}
			}
		}
	}

	if cmdLine != nil {
		ret.overlay(chain, *cmdLine)
	}

	return ret, nil
}

func (base *ScrapeSettings) overlay(chain string, overlay ScrapeSettings) {
	if !overlay.isDefault(chain, "Apps_per_chunk") && overlay.Apps_per_chunk != 0 && overlay.Apps_per_chunk != utils.NOPOS {
		base.Apps_per_chunk = overlay.Apps_per_chunk
	}
	if !overlay.isDefault(chain, "Snap_to_grid") && overlay.Snap_to_grid != 0 && overlay.Snap_to_grid != utils.NOPOS {
		base.Snap_to_grid = overlay.Snap_to_grid
	}
	if !overlay.isDefault(chain, "First_snap") && overlay.First_snap != 0 && overlay.First_snap != utils.NOPOS {
		base.First_snap = overlay.First_snap
	}
	if !overlay.isDefault(chain, "Unripe_dist") && overlay.Unripe_dist != 0 && overlay.Unripe_dist != utils.NOPOS {
		base.Unripe_dist = overlay.Unripe_dist
	}
	if !overlay.isDefault(chain, "Channel_count") && overlay.Channel_count != 0 && overlay.Channel_count != utils.NOPOS {
		base.Channel_count = overlay.Channel_count
	}
	if !overlay.isDefault(chain, "Allow_missing") && overlay.Allow_missing {
		base.Allow_missing = overlay.Allow_missing
	}

	// EXISTING_CODE
	// EXISTING_CODE
}

// EXISTING_CODE
func AllowMissing(chain string) bool {
	s, _ := GetSettings(chain, nil)
	return s.Allow_missing
}

func toEnvStr(name string) string {
	return "TB_SETTINGS_" + strings.ToUpper(strings.Replace(name, "_", "", -1))
}

// EXISTING_CODE