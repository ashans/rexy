# Rexy

[![Build](https://github.com/ashans/rexy/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ashans/rexy/actions/workflows/ci.yml)

Simple Configurable reverse proxy for proxying multiple http services based on context path

## Build & run
### Prerequisites
- GO environment setup
1. Clone repository
2. Run `go get`
3. Run `go build .`
4. Add `config.yaml` file to downloaded location with configuration and run binary (See examples below)
5. Execute **rexy** binary

## Sample use cases
1. match all requests starts with `/api` and proxy to `reqres.in`   
http://localhost:3000/api/users →  https://reqres.in/api/users
```yaml
server:
  port: 3000
endpoints:
  - context: /api
    forward:
      protocol: https
      host: reqres.in
      rewrite: false
```
2. match all requests starts with `/api` and proxy to `reqres.in` by removing matched context   
   http://localhost:3000/api/users →  https://reqres.in/users
```yaml
server:
  port: 3000
endpoints:
  - context: /api
    forward:
      protocol: https
      host: reqres.in
      rewrite: true
```
3. match all requests starts with `/api` and proxy to `reqres.in` by removing matched context and adding service specific context   
   http://localhost:3000/api/users →  https://reqres.in/rest/v2/users
```yaml
server:
  port: 3000
endpoints:
  - context: /api
    forward:
      protocol: https
      host: reqres.in
      rewrite: false
      context: /rest/v2
```
4. match all requests starts with `/api` and proxy to `reqres.in` with custom port   
   http://localhost:3000/api/users →  https://reqres.in:8080/api/users
```yaml
server:
  port: 3000
endpoints:
  - context: /api
    forward:
      protocol: https
      host: reqres.in
      port: 880
      rewrite: false
```

## TODO
- [ ] Add support for complex request matching
- [ ] Add simple request manipulation
- [ ] Add SSL/TLS support
- [ ] Add unit tests