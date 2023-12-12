package base

import (
	"strings"
)

type Import map[string]string

func NewImportList() Import {
	return make(Import, 6)
}

func (i Import) IO() Import {
	i["io"] = ""
	return i
}

func (i Import) Context() Import {
	i["context"] = ""
	return i
}

func (i Import) Time() Import {
	i["time"] = ""
	return i
}

func (i Import) EmptyPB() Import {
	i["google.golang.org/protobuf/types/known/emptypb"] = ""
	return i
}

func (i Import) AnyPB() Import {
	i["google.golang.org/protobuf/types/known/anypb"] = ""
	return i
}

func (i Import) Special(pkg string) Import {
	i[strings.Trim(pkg, `\"`)] = ""
	return i
}

func (i Import) Show() (m map[string]string) {
	m = make(map[string]string, len(i))
	for k := range i {
		s := strings.SplitN(k, " ", 2)
		var v string
		if len(s) > 1 {
			v = s[0] + " "
			k = s[1]
		}
		m[k] = v
	}
	return
}
