package api

import "strings"

type FuncExp struct {
	Name   string
	SubExp []*FuncExp
}

func newExp(name string, subexp []*FuncExp) *FuncExp {
	exp := new(FuncExp)
	exp.Name = name
	exp.SubExp = subexp
	return exp
}

func Compile(expr string) *FuncExp {
	i, n := 0, len(expr)

	var compile func() []*FuncExp
	compile = func() []*FuncExp {
		var token string = ""
		fe := make([]*FuncExp, 0)

		for i < n {
			c := expr[i]
			i++
			switch c {
			case '(':
				if len(token) > 0 {
					fe = append(fe, newExp(token, compile()))
				}
				token = ""
				break

			case ',':
				if len(token) > 0 {
					fe = append(fe, newExp(token, nil))
				}
				token = ""
				break
			case ')':
				if len(token) > 0 {
					fe = append(fe, newExp(token, nil))
				}
				return fe
			default:
				token += strings.Trim(string(c), " ")
			}
		}

		if len(token) > 0 {
			fe = append(fe, newExp(token, nil))
		}
		return fe
	}

	exp := compile()
	if len(exp) > 0 {
		return exp[0]
	}
	return nil
}
