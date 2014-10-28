# googlefinance [![GoDoc](https://godoc.org/github.com/tiloso/googlefinance?status.svg)](http://godoc.org/github.com/tiloso/googlefinance)

Package googlefinance implements a simple interface for fetching historical
prices from google.com/finance.

## Install
```sh
go get github.com/tiloso/googlefinance
```

## Introduction
An example on how to use it:
```go
package main

import (
  "fmt"
  "time"

  "github.com/tiloso/googlefinance"
)

type Quote struct {
  Date                    time.Time
  Open, High, Low, Close  float64
  Volume                  uint32
}

func main() {
  var quotes []Quote
  start, _ := time.Parse("")
  end, _ := time.Parse("")

  if err := googlefinance.Range(start, end).Key("NASDAQ:GOOG").Get(&quotes); err != nil {
    fmt.Printf("err: %v\n", err)
  }
  fmt.Printf("%+v\n", quotes)
}
```

For more details please take a look at [godoc.org](http://godoc.org/github.com/tiloso/googlefinance)
