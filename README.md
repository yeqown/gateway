# gateway

api gateway for golang http server.

## Todos

- [x] HTTP reverse proxy 
- [x] HTTP Cache, support URI with query param and post form
- [] Expansion, support leader node and slave
- [] Permission, RBAC mode
- [] Ratelimit, token bucket alg

## JSON Config file

just run a binary file and do some config so you can just run it easily. proxy config likes: `config.proxy.json`
```json
{
    "port": 8989,
    "logpath": "./logs",
    "proxy_config": {
        "path_rules": [
            {
                "path": "/gw/name",
                "rewrite_path": "/srv/name",
                "method": "GET",
                "server_name": "srv1",
                "combine_req_cfgs": [],
                "need_combine": false
            },
            {
                "path": "/gw/id",
                "rewrite_path": "/srv/id",
                "method": "GET,POST",
                "server_name": "srv1",
                "combine_req_cfgs": [],
                "need_combine": false
            },
            {
                "path": "/gw/combine",
                "rewrite_path": "",
                "method": "GET",
                "server_name": "",
                "combine_req_cfgs": [
                    {
                        "server_name": "srv1",
                        "path": "/srv/id",
                        "field": "combine_id",
                        "method": "GET"
                    },
                    {
                        "server_name": "srv1",
                        "path": "/srv/name",
                        "field": "combine_name",
                        "method": "POST"
                    }
                ],
                "need_combine": true
            }
        ],
        "server_rules": [
            {
                "prefix": "/srv",
                "server_name": "srv1",
                "need_strip_prefix": false
            },
            {
                "prefix": "/striprefix",
                "server_name": "srv1",
                "need_strip_prefix": true
            }
        ],
        "reverse_server_cfgs": {
            "custom_group1": [
                {
                    "name": "srv1",
                    "prefix": "/srv",
                    "addr": "127.0.0.1:8081",
                    "weight": 5
                },
                {
                    "name": "srv1",
                    "prefix": "/srv",
                    "addr": "127.0.0.1:8082",
                    "weight": 5
                }
            ]
        }
    },
    "cacheno_rules": [
        {
            "regular": "^/api/id$"
        }
    ]
}
```

server config likes: `github.com/yeqown/gateway.server.json`
```json
{
    "host": "127.0.0.1",
    "port": "9898"
}
```
