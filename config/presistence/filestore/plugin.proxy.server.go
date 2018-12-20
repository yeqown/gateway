package filestore

import (
	"fmt"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
)

// NewServerRule ...
func (f *JSONFileStore) NewServerRule(r rule.ServerRuler) error {
	cnt := f.ServerRulesCount()
	r.SetID(fmt.Sprintf("%d", cnt))
	// update cfg
	f.cfg.ProxyServerRules = append(f.cfg.ProxyServerRules, r)
	// update f.JSONFile
	cfg := &ServerCfg{
		PPrefix:          r.Prefix(),
		SServerName:      r.ServerName(),
		NNeedStripPrefix: r.NeedStripPrefix(),
		Idx:              r.ID(),
	}
	f.ProxyServerRules = append(f.ProxyServerRules, cfg)
	f.updatedC <- true
	return nil
}

// DelServerRule func
func (f *JSONFileStore) DelServerRule(id string) error {
	idx := utils.Atoi(id)
	length := len(f.ProxyServerRules)
	// update JSONFile
	f.ProxyServerRules[idx] = f.ProxyServerRules[length-1]
	f.ProxyServerRules[length-1] = nil
	f.ProxyServerRules = f.ProxyServerRules[:length-1]
	// updatecfg
	f.cfg.ProxyServerRules[idx] = f.cfg.ProxyServerRules[length-1]
	f.cfg.ProxyServerRules[length-1] = nil
	f.cfg.ProxyServerRules = f.cfg.ProxyServerRules[:length-1]

	f.updatedC <- true

	// 重新设置ID修复ID不对应的情况
	// for idx, r := range f.cfg.ProxyServerRules {
	// 	r.SetID(fmt.Sprintf("%d", idx))
	// }
	return nil
}

// UpdateServerRule func ...
func (f *JSONFileStore) UpdateServerRule(id string, r rule.ServerRuler) error {
	// will panic
	idx := utils.Atoi(id)

	f.ProxyServerRules[idx].NNeedStripPrefix = r.NeedStripPrefix()
	f.ProxyServerRules[idx].PPrefix = r.Prefix()
	f.ProxyServerRules[idx].SServerName = r.ServerName()

	r.SetID(id)
	f.cfg.ProxyServerRules[idx] = r
	f.updatedC <- true
	return nil
}

// ServerRuleByID ...
func (f *JSONFileStore) ServerRuleByID(id string) rule.ServerRuler {
	idx := utils.Atoi(id)
	return f.cfg.ProxyServerRules[idx]
}

// ServerRulesPage func
func (f *JSONFileStore) ServerRulesPage(start, end int) []rule.ServerRuler {
	cnt := f.ServerRulesCount()
	if start >= cnt {
		return nil
	}
	if end > cnt {
		end = cnt
	}

	return f.cfg.ProxyServerRules[start:end]
}

// ServerRulesCount func
func (f *JSONFileStore) ServerRulesCount() int {
	return len(f.cfg.ProxyServerRules)
}
