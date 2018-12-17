package gateway

import (
	"log"

	"github.com/yeqown/gateway/plugin/cache"
	"github.com/yeqown/gateway/plugin/proxy"
	"github.com/yeqown/server-common/utils"
)

// Config contains all plugin rules or gateway config
type Config struct {
	Port         int           `json:"port"`
	Logpath      string        `json:"logpath"`
	ProxyCfg     *proxy.Config `json:"proxy_config"`
	CachenoRules []cache.Rule  `json:"cacheno_rules"`
}

// LoadConfig to load gateway.Config from file
func LoadConfig(file string) (*Config, error) {
	log.Printf("load config from file: %s", file)
	var c = new(Config)
	if err := utils.LoadJSONConfigFile(file, c); err != nil {
		return nil, err
	}
	return c, nil
}
