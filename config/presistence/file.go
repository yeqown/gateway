package presistence

import (
	"fmt"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/server-common/utils"
)

// NewJSONFileStore ...
func NewJSONFileStore(file string) Store {
	var (
		jsonFile = new(JSONFile)
	)

	if err := utils.LoadJSONConfigFile(file, jsonFile); err != nil {
		err = fmt.Errorf("NewJSONFileStore.LoadJSONConfigFile err: %v", err)
		panic(err)
	}

	return &JSONFileStore{
		JSONFile: jsonFile,
	}
}

// JSONFileStore ...
type JSONFileStore struct {
	*JSONFile
	cfg *Instance
}

// Instance config instance ...
func (f *JSONFileStore) Instance() *Instance {
	if f.cfg == nil {
		cfg := new(Instance)
		cfg.Logpath, cfg.Port = f.loadGatewayCfg()
		cfg.Nocache = f.loadNocache()
		cfg.ProxyReverseServers = f.loadProxyReverseServers()
		cfg.ProxyServerRules = f.loadProxyServerRules()
		cfg.ProxyPathRules = f.loadProxyPathRules()
		f.cfg = cfg
	}

	return f.cfg
}

// loadGatewayCfg ...
func (f *JSONFileStore) loadGatewayCfg() (string, int) {
	return f.Logpath, f.Port
}

// loadProxyReverseServers ...
func (f *JSONFileStore) loadProxyReverseServers() map[string][]rule.ReverseServer {
	m := make(map[string][]rule.ReverseServer)
	for key, srvs := range f.ProxyReverseServers {
		rules := make([]rule.ReverseServer, len(srvs))
		for idx, c := range srvs {
			rules[idx] = c
		}
		m[key] = rules
	}
	return m
}

// loadProxyPathRules ...
func (f *JSONFileStore) loadProxyPathRules() []rule.PathRuler {
	rules := make([]rule.PathRuler, len(f.ProxyPathRules))
	for idx, c := range f.ProxyPathRules {
		rules[idx] = c
	}
	return rules
}

// loadProxyServerRules ...
func (f *JSONFileStore) loadProxyServerRules() []rule.ServerRuler {
	rules := make([]rule.ServerRuler, len(f.ProxyServerRules))
	for idx, c := range f.ProxyServerRules {
		rules[idx] = c
	}
	return rules
}

// loadNocache ...
func (f *JSONFileStore) loadNocache() []rule.Nocacher {
	rules := make([]rule.Nocacher, len(f.Nocache))
	for idx, c := range f.Nocache {
		rules[idx] = c
	}
	return rules
}

// JSONFile for filestore to manage config data ...
type JSONFile struct {
	Logpath             string                              `json:"logpath"`
	Port                int                                 `json:"port"`
	ProxyServerRules    []*rule.ServerCfg                   `json:"server_rules"`
	ProxyPathRules      []*rule.PathCfg                     `json:"path_rules"`
	ProxyReverseServers map[string][]*rule.ReverseServerCfg `json:"reverse_server_cfgs"`
	Nocache             []*rule.NocacheCfg                  `json:"cacheno_rules"`
}
