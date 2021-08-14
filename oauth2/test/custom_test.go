package test

import (
	"log"
	"testing"
)

type tvar struct {
	v1 string
	v2 string
}

var t0 = tvar{
	"123",
	"abc",
}

func TestVar(t *testing.T) {
	t3:= &t0
	log.Printf("t0:%p", &t0)
	log.Printf("t3:%p", t3)
	log.Printf("t3_:%p", &t3)
	t1 := t0
	log.Printf("t1:%p", &t1)
	log.Printf("t3:%p", t3)
	log.Printf("t3_:%p", &t3)
	t1.v2 = "789"
	log.Printf("t0:%p", &t0)
	log.Printf("t1:%p", &t1)
	log.Printf("t3:%p", t3)
	log.Printf("t3_:%p", &t3)
}
