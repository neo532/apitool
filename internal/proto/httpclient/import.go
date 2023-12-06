package httpclient

type Import map[string]struct{}

func NewImportList() Import {
	return make(Import, 6)
}

func (i Import) IO() Import {
	i["io"] = struct{}{}
	return i
}

func (i Import) Context() Import {
	i["context"] = struct{}{}
	return i
}

func (i Import) Time() Import {
	i["time"] = struct{}{}
	return i
}

func (i Import) EmptyPB() Import {
	i["google.golang.org/protobuf/types/known/emptypb"] = struct{}{}
	return i
}

func (i Import) AnyPB() Import {
	i["google.golang.org/protobuf/types/known/anypb"] = struct{}{}
	return i
}

func (i Import) Special(pkg string) Import {
	i[pkg] = struct{}{}
	return i
}
