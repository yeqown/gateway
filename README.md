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
    "proxy": [
        {
            "path": "/admin",
            "method": "GET",
            "servers": [
                {
                    "addr": "127.0.0.1:8981/api/target/1",
                    "weight": 5
                },
                {
                    "addr": "127.0.0.1:8982/api/target/1",
                    "weight": 5
                }
            ],
            "combines": [
                {
                    "addr": "127.0.0.1:8981/api/target/1",
                    "field": "field1",
                    "method":"GET"
                },
                {
                    "addr": "127.0.0.1:8981/api/target/2",
                    "field": "field2",
                    "method":"GET"
                },
            ],
            "need_combine": true
        },
        // ....
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
