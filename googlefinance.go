// Package googlefinance implements a simple interface for fetching historical
// prices from google.com/finance.
package googlefinance

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/tiloso/googlefinance/csvstream"
)

const (
	baseURL     = "http://www.google.com/finance/historical?"
	queryLayout = "q=%v&startdate=%v&enddate=%v&output=csv"
)

// Query encapsulates query details
type Query struct {
	date time.Time
	span []time.Time
	key  string
}

// Range returns a range based query
func Range(start, end time.Time) *Query {
	return &Query{
		span: []time.Time{start, end},
	}
}

// Date returns a date based query
func Date(d time.Time) *Query {
	return &Query{
		date: d,
	}
}

func (q *Query) clone() *Query {
	out := &Query{}
	*out = *q
	return out
}

// Key clones the query and sets the provided key on the new instance.
// Stocks' Keys often have following schema: `<exchange>:<symbol>` (e.g.
// "NASDAQ:GOOG" for Google Inc at NASAQ Stock Exchange) and can be found on
// google finance for stocks or other financial data. (Make sure historical
// prices are available as csv export.)
func (q *Query) Key(k string) *Query {
	out := q.clone()
	out.key = k
	return out
}

// Get executes the query and unmarshals fetched data into the provided
// interface
func (q *Query) Get(v interface{}) error {
	rc, err := q.get()
	if err != nil {
		return err
	}
	defer rc.Close()
	return csvstream.NewDecoder(rc).Decode(v)
}

func encDate(d time.Time) string {
	return url.QueryEscape(d.Format("Jan 2, 2006"))
}

func (q *Query) get() (io.ReadCloser, error) {
	addr := baseURL
	if !q.date.IsZero() {
		fd := encDate(q.date)
		addr += fmt.Sprintf(queryLayout, url.QueryEscape(q.key), fd, fd)
	} else {
		addr += fmt.Sprintf(
			queryLayout,
			url.QueryEscape(q.key),
			encDate(q.span[0]),
			encDate(q.span[1]),
		)
	}

	res, err := http.Get(addr)
	if err != nil {
		return nil, fmt.Errorf("err trying to get %v: %v", addr, err)
	}

	if res.StatusCode != http.StatusOK {
		bt, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, fmt.Errorf("err trying to dump response: %s\n", err)
		}
		return nil, fmt.Errorf("err trying to fetch csv: %v\n%s\n", addr, bt)
	}

	return res.Body, nil
}
