package httpclient

import (
	"io"
	"log"
	"os"
	"path/filepath"
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
	CacheTpl       map[string]string
}

// 是否已经存在
func (pb *Proto) FmtWraperName(method *Method) (reply string) {
	return method.Reply + "Wraper"
}

func (pb *Proto) IsNeedWraper(method *Method) (b bool) {
	if method.RespTpl != "" && method.Reply != empty {
		if _, ok := pb.MessageNameMap[pb.FmtWraperName(method)]; !ok {
			return true
		}
	}
	return
}

func (pb *Proto) GetTpl(method *Method) (tpl string, err error) {
	if method.RespTpl == "" {
		return
	}
	var ok bool
	filePath := filepath.Dir(pb.FilePath) + "/" + method.RespTpl + ".tpl"
	if tpl, ok = pb.CacheTpl[filePath]; ok {
		return
	}
	var f *os.File
	f, err = os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		return
	}
	defer f.Close()
	var fb []byte
	fb, err = io.ReadAll(f)
	if err != nil {
		return
	}
	pb.CacheTpl[filePath] = string(fb)
	tpl = string(fb)
	return
}

func (pb *Proto) NewWraper(method *Method, tpl string) {
	// TODO 将返回值的data层也抹掉，code错误放到err里面，这样grpc也可以完全兼容
	tpl = `
// ========== httpClient append ==========
` + tpl + `
// ========== /httpClient append ==========
`
	if pb.WraperMap == nil {
		pb.WraperMap = make(map[string]struct{}, 2)
	}
	pb.WraperMap[strings.ReplaceAll(tpl, "{{ .Reply }}", method.Reply)] = struct{}{}
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
