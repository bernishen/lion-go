package router

import (
	"reflect"
	"strings"
)

const (
	BodyTag    = "route_body"
	ContextTag = "route_ctx"
)

var (
	contextType1 = reflect.TypeOf(Context{})
	contextType2 = reflect.TypeOf(&Context{})
)

type Parameters struct {
	ParamNum    int
	HasBody     bool
	ParamsOrder *[]string
	ParamsType  *[]reflect.Type
	ParamsFill  []bool
	FillOrder   []string
	ParamsData  *map[string]string
}

type IParams interface {
	AddParamData(name string, value string)
	RmoveParamData(name string)
}

func (p *Parameters) AddParamData(name string, value string) {
	if strings.Trim(name, " ") == "" {
		return
	}
	_, ok := (*p.ParamsData)[name]
	if ok {
		return
	}
	(*p.ParamsData)[name] = value
	p.ParamNum++
	if name == BodyTag {
		p.HasBody = true
	}
}

func (p *Parameters) RmoveParamData(name string) {
	n := strings.Trim(name, " ")
	if n == "" {
		return
	}
	_, ok := (*p.ParamsData)[n]
	if !ok {
		return
	}
	delete(*p.ParamsData, n)
	p.ParamNum--
	p.checkBody()
}

func (p *Parameters) checkBody() {
	p.HasBody = false
	for key, _ := range *p.ParamsData {
		name := strings.Trim(key, " ")
		if name == "" && name != BodyTag {
			continue
		}
		p.HasBody = true
		return
	}
}
