package gateway

import "github.com/yeqown/gateway/plugin/proxy"

type gatewayConfig struct {
	Port      int              `json:"port"`
	ProxyCfgs []proxy.PathRule `json:"proxy_config"`
}
