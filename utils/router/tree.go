package router

import (
	"github.com/bernishen/exception"
	"strings"
)

type tree struct {
	root *treenode
}

func initTree() *tree {
	return &tree{
		root: &treenode{
			route:    "/",
			subNodes: make(map[string]*treenode),
		},
	}
}

type treenode struct {
	route      string
	subNodes   map[string]*treenode
	param      *[]string
	controller ILionController
}

func (t *tree) find(route string) (*ILionController, *Parameters, *exception.Exception) {
	n := t.root
	if n == nil {
		return nil, nil, exception.NewException(exception.Error, 1001, "No routes could be found.")
	}

	paths := strings.Split(route, "/")
	pathNum := len(paths)
	if pathNum == 0 {
		return nil, nil, exception.NewException(exception.Error, 1002, "Null route.["+route+"]")
	}

	dic := make(map[string]string)
	paramNum := 0
	for i := 0; i < pathNum; i++ {
		p0 := strings.Trim(paths[i], " ")
		p := strings.ToLower(p0)
		if p == "" {
			continue
		}

		n_, ok := n.subNodes[p]
		if ok {
			n = n_
			continue
		}

		if n == nil || n.param == nil {
			return nil, nil, exception.NewException(exception.Error, 1003, "Found not '"+route+"'.[error location:"+p+"]")
		}
		if paramNum >= len(*(*n).param) {
			return nil, nil, exception.NewException(exception.Error, 1003, "Found not '"+route+"'.[error location:"+p+"]")
		}
		pname := strings.Trim((*(*n).param)[paramNum], " ")
		paramNum++
		if pname == "" {
			return nil, nil, exception.NewException(exception.Error, 1003, "Found not '"+route+"'.[error location:"+p+"]")
		}

		pnameLen := len(pname)
		dic[pname[1:pnameLen]] = p0
	}
	param := Parameters{
		paramNum,
		false,
		nil,
		nil,
		nil,
		nil,
		&dic,
	}
	return &n.controller, &param, nil
}

func (t *tree) add(item *RouterItem) {
	path := strings.ToLower(item.Route)
	paths := strings.Split(path, "/")
	pathNum := len(paths)
	if pathNum == 0 {
		return
	}

	n := t.root
	for i := 0; i < pathNum; i++ {
		p := strings.Trim(paths[i], " ")
		if p == "" {
			continue
		}

		if p[0:1] == ":" {
			n.addParam(p)
			if i == (pathNum - 1) {
				n.controller = item.Controller
			}
		} else {
			if i == (pathNum - 1) {
				n = n.addSubNode(p, &item.Controller)
			} else {
				n = n.addSubNode(p, nil)
			}
		}
	}
}

func (n *treenode) addParam(route string) {
	var p []string
	if (*n).param == nil {
		p = []string{
			route,
		}
	} else {
		num := len(*(*n).param)
		p = make([]string, num+1)
		for i := 0; i < num; i++ {
			p[i] = (*(*n).param)[i]
		}
		p[num] = route
	}
	n.param = &p
}

func (n *treenode) addSubNode(route string, controller *ILionController) *treenode {
	child, ok := n.subNodes[route]
	if !ok {
		path := n.route
		if path != "/" {
			path = path + "/"
		}
		child = &treenode{
			route:    path + route,
			subNodes: make(map[string]*treenode),
		}
	}

	if controller != nil {
		if child.controller != nil {
			panic("路由[" + route + "]处重复添加。")
		}
		child.controller = *controller
	}
	n.subNodes[route] = child
	return child
}
