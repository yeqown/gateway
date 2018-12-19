package filestore

import (
	"fmt"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
)

// NewPathRule func
func (f *JSONFileStore) NewPathRule(r rule.PathRuler) error {
	cnt := f.PathRulesCount()
	rIdx := fmt.Sprintf("%d", cnt)
	r.SetID(rIdx)
	// update cfg
	f.cfg.ProxyPathRules = append(f.cfg.ProxyPathRules, r)
	// update f.JSONFile
	reqs := r.CombineReqCfgs()
	combReqs := make([]*CombineReqCfg, len(reqs))
	for idx, req := range reqs {
		combReqs[idx] = &CombineReqCfg{
			SrvName: req.ServerName(),
			PPath:   req.Path(),
			FField:  req.Field(),
			MMethod: req.Method(),
			Idx:     fmt.Sprintf("%d", idx),
		}
	}

	cfg := &PathCfg{
		PPath:        r.Path(),
		RRewritePath: r.RewritePath(),
		MMethod:      r.Method(),
		SrvName:      r.ServerName(),
		CombReqs:     combReqs,
		NeedComb:     r.NeedCombine(),
		Idx:          rIdx,
	}
	f.ProxyPathRules = append(f.ProxyPathRules, cfg)
	f.updatedC <- true
	return nil
}

// PathRuleByID ...
func (f *JSONFileStore) PathRuleByID(id string) rule.PathRuler {
	idx := utils.Atoi(id)
	return f.cfg.ProxyPathRules[idx]
}

// PathRulesPage func
func (f *JSONFileStore) PathRulesPage(start, end int) []rule.PathRuler {
	cnt := f.PathRulesCount()
	if start >= cnt {
		return nil
	}

	if end >= cnt {
		end = cnt
	}

	return f.cfg.ProxyPathRules[start:end]
}

// PathRulesCount func
func (f *JSONFileStore) PathRulesCount() int {
	return len(f.cfg.ProxyPathRules)
}

// UpdatePathRule func ...
func (f *JSONFileStore) UpdatePathRule(id string, r rule.PathRuler) error {
	idx := utils.Atoi(id)

	f.ProxyPathRules[idx].MMethod = r.Method()
	f.ProxyPathRules[idx].NeedComb = r.NeedCombine()
	f.ProxyPathRules[idx].PPath = r.Path()
	f.ProxyPathRules[idx].RRewritePath = r.RewritePath()
	f.ProxyPathRules[idx].SrvName = r.ServerName()

	combReqs := r.CombineReqCfgs()
	reqs := make([]*CombineReqCfg, len(combReqs))
	for i, combReq := range combReqs {
		reqs[i] = &CombineReqCfg{
			SrvName: combReq.ServerName(),
			PPath:   combReq.Path(),
			FField:  combReq.Field(),
			MMethod: combReq.Method(),
			Idx:     fmt.Sprintf("%d", i),
		}
	}

	f.ProxyPathRules[idx].CombReqs = reqs
	r.SetID(id)
	f.cfg.ProxyPathRules[idx] = r
	f.updatedC <- true
	return nil
}

// DelPathRule func
func (f *JSONFileStore) DelPathRule(id string) error {
	idx := utils.Atoi(id)
	length := len(f.ProxyPathRules)
	// update JSONFile
	f.ProxyPathRules[idx] = f.ProxyPathRules[length-1]
	f.ProxyPathRules[length-1] = nil
	f.ProxyPathRules = f.ProxyPathRules[:length-1]
	// updatecfg
	f.cfg.ProxyPathRules[idx] = f.cfg.ProxyPathRules[length-1]
	f.cfg.ProxyPathRules[length-1] = nil
	f.cfg.ProxyPathRules = f.cfg.ProxyPathRules[:length-1]
	f.updatedC <- true

	// 重新设置ID修复ID不对应的情况
	// for idx, r := range f.cfg.ProxyPathRules {
	// 	r.SetID(fmt.Sprintf("%d", idx))
	// }

	return nil
}
