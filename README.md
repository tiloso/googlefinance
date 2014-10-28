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
  var qs []Quote
  d, _ := time.Parse("2-Jan-06", "22-Oct-14")

  if err := googlefinance.Date(d).Key("NASDAQ:GOOG").Get(&qs); err != nil {
    fmt.Printf("err: %v\n", err)
  }
  fmt.Printf("%+v\n", qs)
  // Output: [{Date:2014-10-22 00:00:00 +0000 UTC Open:529.89 High:539.8 Low:528.8 Close:532.71 Volume:2917183}]
}
```

For more details please take a look at [godoc.org](http://godoc.org/github.com/tiloso/googlefinance)
