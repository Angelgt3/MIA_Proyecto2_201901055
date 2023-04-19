package comandos

type List struct {
	objects []interface{}
}

func New() *List {
	l := new(List)
	l.objects = []interface{}{}
	return l
}

func (l *List) Len() int {
	return len(l.objects)
}

func (l *List) Clear() {
	l.objects = []interface{}{}
}

func (l *List) Clone() *List {
	n := New()
	for _, v := range l.objects {
		n.objects = append(n.objects, v)
	}
	return n
}

func (l *List) Sublist(fromRange int, toRange int) *List {
	if fromRange > toRange || fromRange < 0 || toRange > l.Len()-1 {
		return nil
	}
	nl := New()
	nl.AddAll(l.objects[fromRange : toRange+1])
	return nl
}

func (l *List) IndexOf(elem interface{}) int {
	for k, v := range l.objects {
		if v == elem {
			return k
		}
	}
	return -1
}

func (l *List) GetValue(pos int) interface{} {
	if pos < 0 || pos > l.Len()-1 {
		return nil
	}
	return l.objects[pos]
}

func (l *List) Contains(elem interface{}) bool {
	for _, v := range l.objects {
		if v == elem {
			return true
		}
	}
	return false
}

func (l *List) Add(elem interface{}) {
	l.objects = append(l.objects, elem)
}

func (l *List) AddAll(elems []interface{}) {
	for _, elem := range elems {
		l.Add(elem)
	}
}

func (l *List) ReplaceAll(fromElem interface{}, toElem interface{}) {
	for k, v := range l.objects {
		if v == fromElem {
			l.objects[k] = toElem
		}
	}
}

func (l *List) RemoveFirst(elem interface{}) {
	for i, val := range l.objects {
		if val == elem {
			newArray := l.objects[:i]
			newArray = append(newArray, l.objects[i+1:]...)
			l.objects = newArray
			return
		}
	}
}

func (l *List) RemoveAtIndex(index int) {
	if index < 0 || index > l.Len()-1 {
		return
	}
	newArray := l.objects[:index]
	newArray = append(newArray, l.objects[index+1:]...)
	l.objects = newArray

}

func (l *List) RemoveAll(elem interface{}) {
	l.SubRemoveStarting(elem, 0)
}

func (l *List) SubRemoveStarting(elem interface{}, start int) {
	if start == l.Len() {
		return
	}
	for i := start; i < l.Len(); i++ {
		if l.objects[i] == elem {
			l.RemoveAtIndex(i)
			l.SubRemoveStarting(elem, i)
		}
	}
}

func (l *List) ToArray() []interface{} {
	return l.objects
}

func (l *List) Equals(o *List) bool {
	if l == o {
		return true
	}
	if l.Len() != o.Len() {
		return false
	}
	for k := range l.objects {
		if l.objects[k] != o.objects[k] {
			return false
		}
	}
	return true
}
