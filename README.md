# api-gateway
api gateway for golang http server

# Todos

###* Permission Manager
	Internal RBAC modules

###* Load Balance
	ignore this

###* Request Limit
	limitation rules as following:
	1. Request as spider
	2. IP addr in black list
	3. Request with token (what token? )

###* HTTP Proxy
	with the config file `config.proxy.json` to redirect request to microserver.
	and do some assemble works with differrent microserver. this can be config by the file or
	do this config with an api that is servered by the `github.com/yeqown/gateway`?

###* RPC Caller
	ignore this tmeporarily

###* Cache Pool
	to cache what? or this is just a functional module?

###* Support expansion
	how to design this `github.com/yeqown/gateway` in `master -> slave-node` mode, so I can expand `github.com/yeqown/gateway` easily ?

###* Configurable proxy rule
	must be support file or api method to config proxy or assemble microserver

# Usage

just run a binary file and do some condfig so you can just run it easily


proxy config likes:
```json
# config.proxy.json

{
	"proxy": [
		{
			"listen_path": "/admin",
			"target": "https://api.host.com",
			"strip_listen_path": true			
		},
		{
			"listen_path": "/health",
			"target": "https://api.host.com",
			"strip_listen_path": false	
		}
	]
}
```

server config likes:
```json
# github.com/yeqown/gateway.server.json

{
	"host": "127.0.0.1",
	"port": "9898"
}
```

