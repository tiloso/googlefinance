package csvdecoding_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tiloso/googlefinance/csvdecoding"
)

type Q struct {
	Open, High, Low, Close float64
	Volume                 int
	Date                   time.Time
}

var streamMixedCase = strings.NewReader(`
Date,Open,High,low,Close,Volume
21-Oct-14,24.22,24.95,23.98,24.90,44253
20-Oct-14,24.38,24.38,23.95,24.19,58890
`)

var streamInvalidType = strings.NewReader(`
Date,Open,High,Low,Close,Volume
21-Oct-14,-,24.95,23.98,24.90,44253
20-Oct-14,24.38,24.38,23.95,24.19,58890
`)

func TestDecodeInvalidArgNoPtr(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("expected err, got nil")
		}
	}()

	arg := "a"
	csvdecoding.New(streamMixedCase).Decode(arg)
}

func TestDecodeInvalidArgPtrToNoSlice(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("expected err, got nil")
		}
	}()

	arg := Q{}
	csvdecoding.New(streamMixedCase).Decode(&arg)
}

func TestDecodeInvalidArgPtrToSliceOfNoStruct(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("expected err, got nil")
		}
	}()

	arg := []*Q{}
	csvdecoding.New(streamMixedCase).Decode(&arg)
}

type Q2 struct {
	Open, Low, Max float64
	high, low      float64
	Close          string
	Volume         int
}

func TestDecode(t *testing.T) {
	qs := []Q2{}
	if err := csvdecoding.New(streamMixedCase).Decode(&qs); err != nil {
		t.Errorf("err: %v")
	}

	if len(qs) != 2 {
		t.Errorf("expected slice with len of 2, got: %v", len(qs))
	}

	q := qs[0]
	if q.Open != 24.22 || q.Low != 23.98 || q.Close != "24.90" || q.Volume != 44253 {
		t.Errorf(
			"expected Open/ Low/ Close/ Volume to be 24.22/ 23.98/ 24.90/ 44253, got %v/ %v/ %v/ %v",
			q.Open, q.Low, q.Close, q.Volume,
		)
	}

	if q.high != 0 || q.low != 0 || q.Max != 0 {
		t.Errorf(
			"expected high/ low/ Max to be 0, got %v/ %v /%v",
			q.high, q.low, q.Max,
		)
	}
}

func TestDecodeWithInvalidType(t *testing.T) {
	qs := []Q{}
	err := csvdecoding.New(streamInvalidType).Decode(&qs)

	if err == nil {
		t.Error("expected err, got nil")
	}
	_, ok := err.(*csvdecoding.UnmarshalTypeError)
	if !ok {
		t.Errorf("expect err type of %v, got %v", "*csvdecoding.UnmarshalTypeError", err)
	}

	if len(qs) != 2 {
		t.Errorf("expected length of qs to be 2, got %v", len(qs))
	}

	if qs[0].Open != 0 || qs[0].High != 24.95 || qs[1].Open != 24.38 {
		t.Errorf(
			"expected Open/ High/ Open2 to be 0/ 24.95/ 24.38, got %v/ %v/ %v",
			qs[0].Open, qs[0].High, qs[1].Open,
		)
	}
}

func Example() {
	r := strings.NewReader(`
Date,Open,High,Low,Close,Volume
21-Oct-14,24.22,24.95,23.98,24.90,44253
20-Oct-14,24.38,24.38,23.95,24.19,58890
`)

	type Quote struct {
		Date                   time.Time
		Open, High, Low, Close float64
		Volume                 uint32
	}

	var quotes []Quote
	if err := csvdecoding.New(r).Decode(&quotes); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("%+v\n", quotes)
	// Output: [{Date:2014-10-21 00:00:00 +0000 UTC Open:24.22 High:24.95 Low:23.98 Close:24.9 Volume:44253} {Date:2014-10-20 00:00:00 +0000 UTC Open:24.38 High:24.38 Low:23.95 Close:24.19 Volume:58890}]
}
