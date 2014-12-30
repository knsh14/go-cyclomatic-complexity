package complexity

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

type Ast struct {
	Name     string
	Pos      token.Pos
	Attrs    map[string]string
	Children []*Ast
}

type AstConverter interface {
	ToAst() *Ast
}

func BuildAst(prefix string, n interface{}) (astobj *Ast, err error) {
	v := reflect.ValueOf(n)
	t := v.Type()

	a := Ast{Attrs: Attrs(prefix, n), Children: []*Ast{}}

	if node, ok := n.(ast.Node); ok {
		a.Pos = node.Pos()
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if v.IsValid() == false {
		return nil, nil
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:

		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)

			child, err := BuildAst(fmt.Sprintf("%d", i), f.Interface())
			if err != nil {
				return nil, err
			}
			a.Children = append(a.Children, child)
		}
	case reflect.Map:
		for _, kv := range v.MapKeys() {
			f := v.MapIndex(kv)

			child, err := BuildAst(fmt.Sprintf("%v", kv.Interface()), f.Interface())
			if err != nil {
				return nil, err
			}
			a.Children = append(a.Children, child)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fo := f
			name := t.Field(i).Name

			if f.Kind() == reflect.Ptr {
				f = f.Elem()
			}

			if f.IsValid() == false {
				continue
			}

			if _, ok := v.Interface().(ast.Object); !ok && f.Kind() == reflect.Interface {

				switch f.Interface().(type) {
				case ast.Decl, ast.Expr, ast.Node, ast.Spec, ast.Stmt:

					child, err := BuildAst(name, f.Interface())
					if err != nil {
						return nil, err
					}
					a.Children = append(a.Children, child)
					continue
				}
			}

			switch f.Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
				child, err := BuildAst(name, fo.Interface())
				if err != nil {
					return nil, err
				}
				a.Children = append(a.Children, child)

			default:
				a.Attrs[name] = fmt.Sprintf("%v", f.Interface())
			}
		}
	}

	return &a, nil
}

func Attrs(prefix string, n interface{}) map[string]string {
	attrs := make(map[string]string)

	if prefix != "" {
		attrs["Prefix"] = prefix
	}
	attrs["Type"] = fmt.Sprintf("%T", n)

	v := reflect.ValueOf(n)
	t := v.Type()

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if v.IsValid() == false {
		return nil
	}

	switch v.Kind() {

	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		attrs["length"] = fmt.Sprintf("%d", v.Len())

	case reflect.Struct:
		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				name := t.Field(i).Name
				switch name {
				case "Name", "Kind", "Tok", "Op":
					attrs[name] = fmt.Sprintf("%v", f.Interface())
				}
			}
		}
	default:
		attrs[fmt.Sprintf("%s", n)] = fmt.Sprintf("%s", n)
	}
	return attrs
}
