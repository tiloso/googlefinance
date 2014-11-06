// Package csvdecoding implements decoding of CSV streams (from googlefinance) into slices of
// Go structs.
//
// It is build on top of Go's standard library package encoding/csv and
// takes some inspiration from Unmarshal of encoding/json.
//
// The package has only been tested with CSV streams from googlefinance
// so far and it only allows to decode CSV data into a slice of structs.
package csvdecoding

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// A Decoder reads and decodes CSV objects from an input stream.
type Decoder struct {
	cr     *csv.Reader
	fields []*field
	err    error
}

// New returns a new decoder that reads from r.
func New(r io.Reader) *Decoder {
	return &Decoder{
		cr: csv.NewReader(r),
	}
}

// UnmarshalTypeError will be returned from Decode when a given CSV value cannot
// be unmarshaled into a go value of a given type.
type UnmarshalTypeError struct {
	Value string
	Type  reflect.Type
}

func (e *UnmarshalTypeError) Error() string {
	return fmt.Sprintf("csv: cannot unmarshal %v into go value of type %v", e.Value, e.Type.String())
}

// Decode reads the CSV-encoded value from its input and stores it in the value
// pointed to by v.
//
// If a CSV value is not appropriate for a given type of a struct field, decode
// skips the field and tries to finish decoding as good as it can. If no more
// serious errors are encountered it returns an UnmarshalTypeError describing
// the first occurence of the error.
//
// Decode capitalises CSV headers and maps them against struct field names. (E.g.
// columns named `foo` and `Foo` will be matched agains a struct field name `Foo`,
// the latter occurence overwrites the first one.)
func (d *Decoder) Decode(v interface{}) error {
	resultv := reflect.ValueOf(v)

	if resultv.Kind() != reflect.Ptr ||
		resultv.Elem().Kind() != reflect.Slice ||
		resultv.Elem().Type().Elem().Kind() != reflect.Struct {
		panic("argument must be a pointer to a slice of structs")
	}

	slicev := resultv.Elem()
	elemt := slicev.Type().Elem()

	if err := d.parseHeader(v); err != nil {
		return fmt.Errorf("err decoder.parseHeader: %v", err)
	}

	for {
		records, err := d.cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		elemp := reflect.New(elemt)
		d.unmarshal(records, elemp.Interface())
		slicev = reflect.Append(slicev, elemp.Elem())
	}

	resultv.Elem().Set(slicev)
	return d.err
}

func (d *Decoder) unmarshal(records []string, v interface{}) {
	for i, field := range d.fields {
		if field == nil {
			continue
		}

		if err := field.unmarshal(records[i], v); err != nil && d.err == nil {
			d.err = err
		}
	}
}

func (d *Decoder) parseHeader(v interface{}) error {
	records, err := d.cr.Read()
	if err != nil {
		return fmt.Errorf("err read header row: %v", err)
	}

	slicep := reflect.ValueOf(v)
	slicev := slicep.Elem()

	elemt := slicev.Type().Elem()
	elemv := reflect.New(elemt).Elem()

	d.fields = make([]*field, len(records))
	for i, name := range records {
		// remove byte order mark
		cn := strings.Trim(name, "\ufeff")

		// capitalise to map exported names only / case insensitive match
		cn = strings.ToUpper(cn[:1]) + cn[1:]
		rv := elemv.FieldByName(cn)

		if rv.IsValid() {
			d.fields[i] = &field{
				name: cn,
				kind: rv.Kind(),
				typ:  rv.Type(),
			}
		} else {
			d.fields[i] = nil
		}
	}
	return nil
}

type field struct {
	name string
	kind reflect.Kind
	typ  reflect.Type
}

func newUTE(v string, t reflect.Type) error {
	return &UnmarshalTypeError{
		Value: v,
		Type:  t,
	}
}

func (f *field) unmarshal(v string, i interface{}) error {
	elemv := reflect.ValueOf(i).Elem()

	switch f.kind {
	case reflect.String:
		elemv.FieldByName(f.name).SetString(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		pv, err := strconv.ParseUint(v, 10, 64)

		if err != nil {
			return newUTE(v, f.typ)
		}
		elemv.FieldByName(f.name).SetUint(pv)
	case reflect.Float32, reflect.Float64:
		// TODO strconv.ParseFloat(v, v.Type().Bits())?
		pv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return newUTE(v, f.typ)
		}
		elemv.FieldByName(f.name).SetFloat(pv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return newUTE(v, f.typ)
		}
		elemv.FieldByName(f.name).SetInt(pv)
	case reflect.Struct:
		if f.typ == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse("2-Jan-06", v)
			if err != nil {
				return newUTE(v, f.typ)
			}
			elemv.FieldByName(f.name).Set(reflect.ValueOf(t))
		}
	}

	return nil
}
