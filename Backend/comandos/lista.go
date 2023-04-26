package comandos

type List struct {
	objects []interface{}
}

func (l *List) Len() int {
	return len(l.objects)
}

func (l *List) GetValue(pos int) interface{} {
	if pos < 0 || pos > l.Len()-1 {
		return nil
	}
	return l.objects[pos]
}

func (l *List) Add(elem interface{}) {
	l.objects = append(l.objects, elem)
}
