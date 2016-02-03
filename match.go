// +build linux darwin

package main

import (
	"github.com/upyun/go-sdk/upyun"
	"path/filepath"
)

type MatchConfig struct {
	wildcard      string
	timestampType string /* before, after*/
	timestamp     int
	itemType      string /* file, folder */
}

func (mc *MatchConfig) IsMatched(upInfo *upyun.FileInfo) bool {
	if mc.wildcard != "" {
		if same, _ := filepath.Match(mc.wildcard, upInfo.Name); !same {
			return false
		}
	}
	if mc.timestamp != 0 {
		dist := int(upInfo.Time.Unix()) - mc.timestamp
		typ := mc.timestampType
		if dist < 0 && typ == "before" {
			return false
		}
		if dist > 0 && typ == "after" {
			return false
		}
	}
	if mc.itemType != "" {
		if upInfo.Type != mc.itemType {
			return false
		}
	}
	return true
}
