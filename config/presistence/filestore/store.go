package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/server-common/utils"
)

var _ presistence.Store = &JSONFileStore{}

var (
	// ErrGroupNotExist ...
	ErrGroupNotExist = errors.New("filestore.reverseServer group not exist")
	// ErrGroupExists ....
	ErrGroupExists = errors.New("filestore.reverseServer group exsits")
)

// NewJSONFileStore ...
func NewJSONFileStore(filename string) presistence.Store {
	var (
		jsonFile = new(JSONFile)
	)

	if err := utils.LoadJSONConfigFile(filename, jsonFile); err != nil {
		err = fmt.Errorf("NewJSONFileStore.LoadJSONConfigFile err: %v", err)
		panic(err)
	}

	s := &JSONFileStore{
		JSONFile: jsonFile,
		filename: filename,
		updatedC: make(chan bool, 5), // has buffer, avoid deadlock
	}

	// related to updatedC
	go s.presist()

	return s
}

// JSONFile for filestore to manage config data ...
type JSONFile struct {
	Logpath             string                         `json:"logpath"`
	Port                int                            `json:"port"`
	ProxyServerRules    []*ServerCfg                   `json:"server_rules"`
	ProxyPathRules      []*PathCfg                     `json:"path_rules"`
	ProxyReverseServers map[string][]*ReverseServerCfg `json:"reverse_server_cfgs"`
	Nocache             []*NocacheCfg                  `json:"cacheno_rules"`
}

// JSONFileStore ...
type JSONFileStore struct {
	*JSONFile
	filename string
	cfg      *presistence.Instance
	updatedC chan bool // update channel
}

// Instance config instance ...
func (f *JSONFileStore) Instance() *presistence.Instance {
	if f.cfg == nil {
		cfg := new(presistence.Instance)
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
			c.SetID(fmt.Sprintf("%s#%d", key, idx))
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
		c.SetID(fmt.Sprintf("%d", idx))
		rules[idx] = c
	}
	return rules
}

// loadProxyServerRules ...
func (f *JSONFileStore) loadProxyServerRules() []rule.ServerRuler {
	rules := make([]rule.ServerRuler, len(f.ProxyServerRules))
	for idx, c := range f.ProxyServerRules {
		c.SetID(fmt.Sprintf("%d", idx))
		rules[idx] = c
	}
	return rules
}

// loadNocache ...
func (f *JSONFileStore) loadNocache() []rule.Nocacher {
	rules := make([]rule.Nocacher, len(f.Nocache))
	for idx, c := range f.Nocache {
		c.SetID(fmt.Sprintf("%d", idx))
		rules[idx] = c
	}
	return rules
}

// UpdateGateConfig func ...
func (f *JSONFileStore) UpdateGateConfig(logpath string, port int) error {
	f.cfg.Port, f.Port = port, port
	f.cfg.Logpath, f.Logpath = logpath, logpath
	return nil
}

func (f *JSONFileStore) presist() {
	fd, err := os.OpenFile(f.filename, os.O_RDWR, os.ModePerm)
	defer fd.Close()

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-f.updatedC:
			log.Println("updated JSONFile")
			byts, err := json.MarshalIndent(f.JSONFile, "", "\t")
			if err != nil {
				panic(err)
			}

			// clear file then write
			if err := fd.Truncate(0); err == nil {
				_, err = fd.Write(byts)
				if err != nil {
					panic(err)
				}
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
