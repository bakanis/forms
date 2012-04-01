package forms

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
)

var (
	converter_error = errors.New("converter_error")

	int_converter ConverterFunc = func(in string) (out interface{}, err error) {
		out = 2
		return
	}
	error_converter ConverterFunc = func(in string) (out interface{}, err error) {
		out = in
		err = converter_error
		return
	}
)

func fatal_converter(t *testing.T) ConverterFunc {
	return func(in string) (out interface{}, err error) {
		t.Fail()
		out = in
		return
	}
}

func TestConverterInt(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: int_converter},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}
	rval := res.Value.(map[string]interface{})
	if ex, ok := rval["foo"]; !ok || ex.(int) != 2 {
		t.Fatalf("Expected %v. Got %v", 2, ex)
	}
}

func TestConverterError(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: error_converter},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
	if ex, ok := res.Errors["foo"]; !ok || ex != converter_error {
		t.Fatalf("Expected %v. Got %v", converter_error, ex)
	}
}

func TestConverterNotCalledOnValidatorError(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{
				Name:       "foo",
				Validators: []Validator{error_validator},
				Converter:  fatal_converter(t),
			},
		},
	}
	f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
}

//test specific converters
func TestConvertersIntConverter(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: IntConverter},
		},
	}
	for i := -10; i <= 10; i++ {
		res := f.Load(create_req(url.Values{
			"foo": {fmt.Sprint(i)},
		}))
		if ex, ok := res.Errors["foo"]; ok || ex != nil {
			t.Errorf("Expected %v. Got %v", nil, ex)
			continue
		}
		rval := res.Value.(map[string]interface{})
		if ex, ok := rval["foo"].(int); !ok || ex != i {
			t.Errorf("Expected %v. Got %v", i, ex)
		}
	}
}

func TestConvertersFloat32Converter(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: Float32Converter},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"3.14159"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}
	rval := res.Value.(map[string]interface{})
	if ex, ok := rval["foo"].(float32); !ok || ex != 3.14159 {
		t.Fatalf("Expected %v. Got %v", 3.14159, ex)
	}
}

func TestConvertersFloat64Converter(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: Float64Converter},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"3.14159"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}
	rval := res.Value.(map[string]interface{})
	if ex, ok := rval["foo"].(float64); !ok || ex != 3.14159 {
		t.Fatalf("Expected %v. Got %v", 3.14159, ex)
	}
}
