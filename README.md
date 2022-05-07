# rediscache

[![Pipeline](https://github.com/ercantopuz/rediscache/actions/workflows/ci.yaml/badge.svg)](https://github.com/ercantopuz/rediscache/actions/workflows/ci.yaml)
[![codecov](https://codecov.io/gh/ercantopuz/rediscache/branch/master/graph/badge.svg?token=DV0KA0K3X8)](https://codecov.io/gh/ercantopuz/rediscache)
[![Release](https://img.shields.io/github/release/ercantopuz/rediscache.svg?style=flat-square)](https://github.com/ercantopuz/rediscache/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/ercantopuz/rediscache)](https://goreportcard.com/report/github.com/ercantopuz/rediscache)

## Installation

This package makes it easy to abstract your Redis connection and use Get and Set methods.

```sh
go get github.com/ercantopuz/rediscache
```

Import it in your code:
```go
import "github.com/ercantopuz/rediscache"
```
### Quick start

```go
package main

import "github.com/ercantopuz/rediscache"

func main() {
	
	// It has all the properties of the redis.Options
	redisCacheOptions := rediscache.Options{
		Addr:               "",
		Username:           "",
		Password:           "",
		DB:                 0,
	}
	redisCache := rediscache.NewRedisCache(
		expire,
		redisCacheOptions)
}
```

```go
package service

import (
	"context"
	"github.com/ercantopuz/rediscache"
)

type service struct {
	redisCache rediscache.Cache
}

func (s *service) GetData() {
	var err error
	var expectedData model.Data{}

	expectedData, err = s.redisCache.Get(
		context.Background(),
		key,
		expectedData)
	
	/****/
	
	err = s.redisCache.Set(
		context.Background(),
		key,
		dataFromRepository)
}
```
