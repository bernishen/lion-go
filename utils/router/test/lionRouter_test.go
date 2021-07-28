package test

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Perple struct {
	Name string
	Age  int
}

func TestMapProperty(t *testing.T) {
	var a1 bool
	t.Log(reflect.TypeOf(a1).Name())
	var a2 uint
	t.Log(reflect.TypeOf(a2).Name())
	var a3 uint8
	t.Log(reflect.TypeOf(a3).Name())
	var a4 uint16
	t.Log(reflect.TypeOf(a4).Name())
	var a5 uint32
	t.Log(reflect.TypeOf(a5).Name())
	var a6 uint64
	t.Log(reflect.TypeOf(a6).Name())
	var a7 int
	t.Log(reflect.TypeOf(a7).Name())
	var a8 int8
	t.Log(reflect.TypeOf(a8).Name())
}

func TestType(t *testing.T) {
	var a *int
	at := reflect.ValueOf(a)
	typeOfA := reflect.TypeOf(a)
	aIns := reflect.New(typeOfA)
	ak := at.Kind()
	rk := aIns.Kind()
	if ak != reflect.Ptr && rk == reflect.Ptr {
		rIn := aIns.Elem().Interface()

		t.Log(rIn)
	}

	t.Log(aIns.Type(), aIns.Kind())
}

func TestType2(t *testing.T) {
	in := Perple{
		"Tom",
		10,
	}
	vbyte, _ := json.Marshal(in)
	value := string(vbyte)
	v := []byte(value)
	var a Perple
	json.Unmarshal(v, &a)

	vt := reflect.TypeOf(in)
	vin := reflect.New(vt).Elem().Interface()
	json.Unmarshal(v, &vin)

	t.Log(v)
}
