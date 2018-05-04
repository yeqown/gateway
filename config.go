package gateway

import "github.com/yeqown/gateway/plugin/proxy"

type gatewayConfig struct {
	Port      int            `json:"port"`
	ProxyCfgs []proxy.Config `json:"proxy_config"`
}
