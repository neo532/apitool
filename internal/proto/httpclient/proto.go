package httpclient

import (
	"log"
	"os"
	"strings"

	"github.com/emicklei/proto"
)

type Proto struct {
	FilePath       string
	PackageName    string
	Package        string
	MessageNameMap map[string]struct{}
	Services       []*Service
	WraperMap      map[string]struct{}
}

// 是否已经存在
func (pb *Proto) FmtWraperName(method *Method) (reply string) {
	return method.Reply + "Wraper"
}

func (pb *Proto) IsNeedWraper(method *Method) (b bool) {
	if method.RespTpl == "" && method.Reply != empty {
		if _, ok := pb.MessageNameMap[pb.FmtWraperName(method)]; !ok {
			return true
		}
	}
	return
}

func (pb *Proto) NewWraper(method *Method) {
	// TODO 将返回值的data层也抹掉，code错误放到err里面，这样grpc也可以完全兼容
	key := "{{ .Reply }}"
	tpl := `
// ========== httpClient append ==========
message ` + key + `Wraper { 
    int32 code = 1;
    string message = 2;
    ` + key + ` data = 3;
}
// ========== /httpClient append ==========
`
	if pb.WraperMap == nil {
		pb.WraperMap = make(map[string]struct{}, 2)
	}
	pb.WraperMap[strings.ReplaceAll(tpl, key, method.Reply)] = struct{}{}
	return
}

func (pb *Proto) ReadProtoFile() (definition *proto.Proto, err error) {
	reader, err := os.Open(pb.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err = parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (pb *Proto) AppendWraper() (err error) {
	if len(pb.WraperMap) == 0 {
		return
	}

	var f *os.File
	f, err = os.OpenFile(pb.FilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	for row := range pb.WraperMap {
		if _, err = f.WriteString(row); err != nil {
			return
		}
	}
	return
}
