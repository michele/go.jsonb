[![Build Status](https://travis-ci.org/michele/go.jsonb.svg?branch=master)](https://travis-ci.org/michele/go.jsonb) [![Go Report Card](https://goreportcard.com/badge/github.com/michele/go.jsonb)](https://goreportcard.com/report/github.com/michele/go.jsonb)

## Postrgres' JSONB implementation for Go

### Usage

Include it in your project:

```
go get github.com/michele/go.jsonb

or

glide get github.com/michele/go.jsonb
```

then use it in your project:

```
import (
  jsonb "github.com/michele/go.jsonb"
)

...

type Example struct {
  Field jsonb.JSONB
}
```