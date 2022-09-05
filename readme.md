# ❤️ Kiyora ❤️

[![GoDoc](https://pkg.go.dev/badge/github.com/ItsMyEyes/kiyora?status.svg)](https://github.com/ItsMyEyes/kiyora)
![Versioning](https://img.shields.io/badge/Version-1.0.0-brightgreen)
![License](https://img.shields.io/badge/License-MIT-blue)
![Go](https://img.shields.io/badge/Go-1.18-blue)

Kiyora is a framework written in Go (Golang). It features a martini-like API with performance that is up to 40 times faster thanks to [httprouter](https://github.com/julienschmidt/httprouter). If you need performance and good productivity, you will love me.

## Contents

- [Kiyora Framework](#)
  - [Installation](#installation)
  - [Api](#api)
    - [Version](api/v1/) [[v1](api/v1/)]
      - [Controllers](api/v1/controllers)
      - [Services](api/v1/services/)
  - [Cmd](cmd/)
    - [Http](cmd/http)
    - [Another](#another-else)
  - [Configs](configs/)
  - [Entity](entity/)
  - [Library](library/)
    - [httpserver](library/httpserver/ginserver/)
    - [log](library/log/)
    - [logger](library/logger/v2/)
  - [Models](#models)
    - [requests](models/requests/)
    - [response](models/response/)
  - [Repository](repository/)
  - [Routers](routers/)
  - [Scripts](scripts/)
  - [Utils](utils/)  
    - [constants](utils/constants/)
    - [helpers](utils/helpers/)
    - [validate](utils/validate/)

## Installation

To install kiyora package, you need to install Go and set your Go workspace first.
1. ou first need [Go](https://go.dev/dl/) installed (version 1.15+ is required), then you can use the below Go command to install kiyora.
2. clonning the project
```sh
$ git clone https://github.com/ItsMyEyes/kiyora.git
```
3. cd to project
```sh
$ cd kiyora
```
4. download dependencies
```sh
$ go mod download
```
5. You just need configuration at
   - [Configs](configs/)
     - [configs/cors.go](configs/cors.go)
     - [configs/redis.go](configs/redis.go)
     - [Coming Soon 🔜🔜](#soon)
6. run the project
```sh
$ go run cmd/http/main.go
```
7. open your browser and go to http://localhost:8080

## Thanks to..
<!-- make message thanks -->
- [Gin](https://github.com/gin-gonic/gin)
- [Redis](https://github.com/go-redis/redis)
- [Resty](https://github.com/go-resty/resty)

