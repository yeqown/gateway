package filestore

import (
	"fmt"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
)

// NewNocacheRule func
func (f *JSONFileStore) NewNocacheRule(c rule.Nocacher) error {
	cnt := len(f.Nocache)
	idx := fmt.Sprintf("%d", cnt)
	c.SetID(idx)
	// update cfg
	f.cfg.Nocache = append(f.cfg.Nocache, c)
	// update f.JSONFile
	cfg := &NocacheCfg{
		Regexp: c.Regular(),
		Idx:    idx,
	}
	f.Nocache = append(f.Nocache, cfg)
	return nil
}

// DelNocacheRule func
func (f *JSONFileStore) DelNocacheRule(id string) error {
	idx := utils.Atoi(id)
	length := len(f.Nocache)
	// update JSONFile
	f.Nocache[idx] = f.Nocache[length-1]
	f.Nocache[length-1] = nil
	f.Nocache = f.Nocache[:length-1]
	// updatecfg
	f.cfg.Nocache[idx] = f.cfg.Nocache[length-1]
	f.cfg.Nocache[length-1] = nil
	f.cfg.Nocache = f.cfg.Nocache[:length-1]

	// 重新设置ID修复ID不对应的情况
	// for idx, r := range f.cfg.Nocache {
	// 	r.SetID(fmt.Sprintf("%d", idx))
	// }

	return nil
}

// UpdateNocacheRule func ...
func (f *JSONFileStore) UpdateNocacheRule(id string, c rule.Nocacher) error {
	idx := utils.Atoi(id)
	c.SetID(id)
	// update JSONFile
	f.Nocache[idx].Regexp = c.Regular()
	// update cfg
	f.cfg.Nocache[idx] = c

	f.updatedC <- true
	return nil
}

// NocacheRules ...
func (f *JSONFileStore) NocacheRules(start, end int) []rule.Nocacher {
	cnt := f.NocacheRulesCount()
	if start >= cnt {
		return nil
	}

	if end >= cnt {
		end = cnt
	}

	return f.cfg.Nocache[start:end]
}

// NocacheRuleByID ...
func (f *JSONFileStore) NocacheRuleByID(id string) rule.Nocacher {
	idx := utils.Atoi(id)
	return f.cfg.Nocache[idx]
}

// NocacheRulesCount ...
func (f *JSONFileStore) NocacheRulesCount() int {
	return len(f.Nocache)
}
