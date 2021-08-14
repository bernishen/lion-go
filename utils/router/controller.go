package router

import (
	"encoding/json"
	"github.com/bernishen/lion-go/utils/exception"
	"github.com/bernishen/lion-go/utils/result"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type BaseController struct {
	route      string
	router     *LionRouter
	funcGet    HandleFunctions
	funcPost   HandleFunctions
	funcPut    HandleFunctions
	funcDelete HandleFunctions

	ResetReturn func(data *interface{}, ex *exception.Exception) (*[]byte, *exception.Exception)
}

type HandleFunction struct {
	Func        reflect.Value
	ParamsNum   int
	ParamsOrder []string
	ParamsType  []reflect.Type
	HasBody     bool
}

type HandleFunctions map[int]HandleFunction

func InitController(route string) *BaseController {
	return &BaseController{
		route:      route,
		funcGet:    make(map[int]HandleFunction),
		funcPost:   make(map[int]HandleFunction),
		funcPut:    make(map[int]HandleFunction),
		funcDelete: make(map[int]HandleFunction),

		ResetReturn: nil,
	}
}

func defaultResetReturn(data *interface{}, ex *exception.Exception) (*[]byte, *exception.Exception) {
	var dRet *result.R
	if ex != nil {
		dRet = result.Faid(ex.Message)
	} else {
		dRet = result.OK(data)
	}

	ret, err := json.Marshal(dRet)
	if err != nil {
		return nil, ex
	}
	return &ret, nil
}

func (c *BaseController) FuncRun(method string, params *Parameters, req *http.Request) (*[]byte, *exception.Exception) {
	var handle HandleFunction
	var ok bool
	switch method {
	case http.MethodGet:
		handle, ok = c.funcGet.findFunc(params)
	case http.MethodPost:
		handle, ok = c.funcPost.findFunc(params)
	case http.MethodPut:
		handle, ok = c.funcPut.findFunc(params)
	case http.MethodDelete:
		handle, ok = c.funcDelete.findFunc(params)
	default:
		return nil, exception.NewException(exception.Error, 1001, "Found not function.")
	}
	if !ok {
		msg := "Found not function with " + strconv.Itoa(params.ParamNum) + " general arguments"
		if params.HasBody {
			msg += "[The body exist]."
		} else {
			msg += "[The body does not exist]."
		}
		return nil, exception.NewException(exception.Error, 1002, msg)
	}

	ctx := Context{
		Token:   "",
		Request: req,
	}
	ps := fillParams(&handle, params, &ctx)
	for i, b := range params.ParamsFill {
		if b {
			continue
		}
		alias := ""
		if i < len(*params.ParamsOrder) {
			alias = (*params.ParamsOrder)[i]
		}
		msg := "The parameter is not fill.[alisa:" + alias + ",index:" + strconv.Itoa(i) + "]"
		return nil, exception.NewException(exception.Error, 1003, msg)
	}
	retV := handle.Func.Call(*ps)
	ret := retV[0].Interface()
	ex := retV[1].Interface().(*exception.Exception)
	if c.ResetReturn != nil {
		return c.ResetReturn(&ret, ex)
	}
	return defaultResetReturn(&ret, ex)
}

func (h *HandleFunctions) findFunc(params *Parameters) (HandleFunction, bool) {
	handle1, ok1 := (*h)[params.ParamNum]
	if ok1 && handle1.HasBody == params.HasBody {
		return handle1, true
	}

	var handle2 HandleFunction
	var ok2, retOK bool
	if params.HasBody {
		handle2, ok2 = (*h)[params.ParamNum-1]
		retOK = ok2 && !handle2.HasBody
		if retOK {
			params.RmoveParamData(BodyTag)
		}
	} else {
		handle2, ok2 = (*h)[params.ParamNum+1]
		retOK = ok2 && handle2.HasBody
		if retOK {
			params.AddParamData(BodyTag, "")
		}
	}

	return handle2, ok2
}

func fillParams(handle *HandleFunction, params *Parameters, ctx *Context) *[]reflect.Value {
	pos := handle.ParamsOrder
	pts := handle.ParamsType
	pfs := make([]bool, handle.ParamsNum)
	fos := make([]string, handle.ParamsNum)
	pds := *(params.ParamsData)

	ps := make([]reflect.Value, handle.ParamsNum)
	var fail, context []int
	for i := 0; i < handle.ParamsNum; i++ {
		pt := pts[i]

		alias := pos[i]
		if alias == "" {
			pfs[i] = false
			if pt == contextType1 || pt == contextType2 {
				pfs[i] = true
				fos[i] = alias
				context = append(context, i)
			} else {
				fail = append(fail, i)
			}
			continue
		}

		pd, ok := pds[alias]
		if !ok {
			pfs[i] = false
			if pt == contextType1 || pt == contextType2 {
				pfs[i] = true
				fos[i] = alias
				context = append(context, i)
			} else {
				fail = append(fail, i)
			}
			continue
		}

		ps[i] = *convertParam(&pt, pd)
		pfs[i] = true
		fos[i] = alias
		delete(pds, alias)
	}

	iFail := 0
	maxFail := len(fail)
	for alias, pd := range pds {
		if fail == nil || iFail > maxFail {
			break
		}
		i := fail[iFail]
		pt := pts[i]

		ps[i] = *convertParam(&pt, pd)
		pfs[i] = true
		fos[i] = alias
		iFail++
	}

	params.ParamsOrder = &(handle.ParamsOrder)
	params.ParamsType = &(handle.ParamsType)
	params.ParamsFill = pfs
	params.FillOrder = fos
	ctx.Params = *params
	for _, i := range context {
		switch pts[i] {
		case contextType1:
			ps[i] = reflect.ValueOf(*ctx)
		case contextType2:
			ps[i] = reflect.ValueOf(ctx)
		default:
		}

	}

	return &ps
}

func convertParam(t *reflect.Type, value string) *reflect.Value {
	name := (*t).Name()
	switch name {
	case "string":
		ret := reflect.ValueOf(value)
		return &ret
	default:
		v := []byte(value)
		in := reflect.New(*t).Interface()
		err := json.Unmarshal(v, &in)
		if err != nil {
			return nil
		}
		temp := reflect.ValueOf(in).Elem().Interface()
		ret := reflect.ValueOf(temp)
		return &ret
	}
}

func addFunc(container *HandleFunctions, f *interface{}, paramsOrder *[]string) {
	if f == nil {
		return
	}
	t := reflect.TypeOf(*f)
	outNum := t.NumOut()
	if outNum != 2 {
		panic("The function must has 2 return and the second return value of the function must be '*exception.Exception'.")
	}
	//if t.Out(0) != reflect.TypeOf(tin) {
	//	panic("The first return value of the function must be '*interface{}'.")
	//}
	if t.Out(1) != reflect.TypeOf(&exception.Exception{}) {
		panic("The second return value of the function must be '*exception.Exception'.")
	}

	inNum := t.NumIn()
	hasBody := false
	pOrder := make([]string, inNum)
	for i, s := range *paramsOrder {
		pOrder[i] = s
		if strings.ToLower(s) == BodyTag {
			hasBody = true
		}
	}

	pType := make([]reflect.Type, inNum)
	tagNum := inNum
	bodyIndex := -1
	for i := inNum - 1; i >= 0; i-- {
		pType[i] = t.In(i)
		if pType[i] == contextType1 || pType[i] == contextType2 {
			pOrder[i] = ContextTag
			tagNum--
			continue
		}
		if bodyIndex >= 0 {
			continue
		}
		bodyIndex = i
	}
	orderNum := len(*paramsOrder)
	if bodyIndex >= 0 && orderNum < tagNum {
		pOrder[bodyIndex] = BodyTag
		hasBody = true
	}

	h := HandleFunction{
		reflect.ValueOf(*f),
		inNum,
		pOrder,
		pType,
		hasBody,
	}

	c := *container
	_, ok := c[tagNum]
	if ok {
		panic("A function with " + strconv.Itoa(inNum) + " general arguments already exists.")
	}
	c[tagNum] = h
}

func (c *BaseController) Get(f interface{}, paramsOrder ...string) *BaseController {
	addFunc(&c.funcGet, &f, &paramsOrder)
	return c
}

func (c *BaseController) Post(f interface{}, paramsOrder ...string) *BaseController {
	addFunc(&c.funcPost, &f, &paramsOrder)
	return c
}

func (c *BaseController) Put(f interface{}, paramsOrder ...string) *BaseController {
	addFunc(&c.funcPut, &f, &paramsOrder)
	return c
}

func (c *BaseController) Delete(f interface{}, paramsOrder ...string) *BaseController {
	addFunc(&c.funcDelete, &f, &paramsOrder)
	return c
}

func (c *BaseController) WebSocket(w http.ResponseWriter, req *http.Request, params *map[string]string) {
	//pNum := len(*params)
	//handle, ok := c.funcGet[(pNum + 2)]
	//if !ok {
	//	return
	//}
	//
	//ps := fillParams(&handle, params)
	//pa := make([]reflect.Value, (pNum + 2))
	//pa[0] = reflect.ValueOf(w)
	//pa[1] = reflect.ValueOf(req)
	//for i := 0; i < pNum; i++ {
	//	pa[i+2] = (*ps)[i]
	//}
	//handle.Func.Call(pa)
}
