package api

type Api interface {
	Value(name string, from int64, to int64) interface{}
	Call(name string, args interface{}) interface{}
}

func Eval(expr string, from int64, to int64, api Api) interface{} {
	var exp *FuncExp = Compile(expr)
	if exp == nil {
		return nil
	}

	var eval func(*FuncExp, chan interface{})
	eval = func(fe *FuncExp, ch chan interface{}) {

		if n := len(fe.SubExp); n > 0 {
			argv, achn := make([]interface{}, n), make([]chan interface{}, n)

			for i := 0; i < n; i++ {
				achn[i] = make(chan interface{})
				go eval(fe.SubExp[i], achn[i])
			}

			for i := 0; i < n; i++ {
				argv[i] = <-achn[i]
				close(achn[i])
			}

			ch <- api.Call(fe.Name, argv)
			return
		}
		ch <- api.Value(fe.Name, from, to)
	}

	vchn := make(chan interface{})
	{
		go eval(exp, vchn)
	}
	value := <-vchn
	close(vchn)

	return value
}
