# evaluator

# retry

[![Release](https://img.shields.io/github/release/mashiike/evaluator.svg?style=flat-square)](https://github.com/mashiike/evaluator/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
![test workflow](https://github.com/mashiike/evaluator/actions/workflows/test.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/mashiike/evaluator?style=flat-square)](https://goreportcard.com/report/github.com/mashiike/evaluator)
[![GoDoc](https://godoc.org/github.com/mashiike/evaluator?status.svg&style=flat-square)](http://godoc.org/github.com/mashiike/evaluator)


A simple library for expression evaluation.
### SYNOPSIS

```golang
	e, err := evaluator.New("(var1 + 0.5) * var2")
    if err != nil {
        log.Fatal(err)
    }
	ans, err := e.Eval(evaluator.Variables{"var1": 0.5, "var2": 3})
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println(ans)
```

see [godoc.org/github.com/mashiike/evaluator](https://godoc.org/github.com/mashiike/evaluator).

## Author

Copyright (c) 2017 KAYAC Inc.

## LICENSE

MIT
