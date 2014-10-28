package googlefinance_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tiloso/googlefinance"
)

type Quote struct {
	Date                   time.Time
	Open, High, Low, Close float64
	Volume                 uint32
}

func TestDateImport(t *testing.T) {
	date, _ := time.Parse("2006-Jan-02", "2014-Oct-21")
	quotes := []Quote{}

	if err := googlefinance.Date(date).Key("fra:dbk").Get(&quotes); err != nil {
		t.Errorf("not ok: %v", err)
	}

	if len(quotes) != 1 {
		t.Errorf("got %v quotes, expected %v", len(quotes), 1)
	}

	if quotes[0].Open != 24.22 {
		t.Errorf("got open of %v, expected %v", quotes[0].Open, 24.22)
	}

	if quotes[0].Date != date {
		t.Errorf("got date %v, expected %v", quotes[0].Date, date)
	}
}

func TestRangeImport(t *testing.T) {
	start, _ := time.Parse("2006-Jan-02", "2014-Oct-13")
	end, _ := time.Parse("2006-Jan-02", "2014-Oct-17")
	quotes := []Quote{}

	if err := googlefinance.Range(start, end).Key("fra:dbk").Get(&quotes); err != nil {
		t.Errorf("not ok: %v", err)
	}

	if len(quotes) != 5 {
		t.Errorf("got %v quotes, expected %v", len(quotes), 5)
	}
}

func Example() {
	type Quote struct {
		Date                   time.Time
		Open, High, Low, Close float64
		Volume                 uint32
	}
	var qs []Quote
	d, _ := time.Parse("2-Jan-06", "22-Oct-14")

	if err := googlefinance.Date(d).Key("NASDAQ:GOOG").Get(&qs); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("%+v\n", qs)
	// Output: [{Date:2014-10-22 00:00:00 +0000 UTC Open:529.89 High:539.8 Low:528.8 Close:532.71 Volume:2917183}]
}

func ExampleRange() {
	type Quote struct {
		Date                   time.Time
		Open, High, Low, Close float64
		Volume                 int
	}
	var qs []Quote
	start, _ := time.Parse("2006-Jan-02", "2014-Oct-13")
	end, _ := time.Parse("2006-Jan-02", "2014-Oct-14")

	if err := googlefinance.Range(start, end).Key("fra:dbk").Get(&qs); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("%+v\n", qs)
	// Output: [{Date:2014-10-14 00:00:00 +0000 UTC Open:24.84 High:25.22 Low:24.66 Close:25.07 Volume:44751} {Date:2014-10-13 00:00:00 +0000 UTC Open:25.03 High:25.34 Low:24.82 Close:25.15 Volume:49379}]
}
