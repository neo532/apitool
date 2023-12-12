package service

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/neo532/apitool/internal/base"
	"github.com/neo532/apitool/internal/proto/entity"
)

//nolint:lll
var serviceTemplate = `
{{- /* delete empty line */ -}}
package {{ .PackageName }}

import (
	{{range $key,$value := .ImportList }}{{ $value }}"{{ $key }}"
	{{ end }}
    pb "{{ .Package }}"
)

type {{ .Service }}{{ .ServiceType }} struct {
    pb.Unimplemented{{ .Service }}Server
	tag string
}

func New{{ .Service }}{{ .ServiceType }}(
) *{{ .Service }}{{ .ServiceType }} {
    return &{{ .Service }}{{ .ServiceType }}{
		tag : "{{ .PackageName }}.{{ .Service }}{{ .ServiceType }}",
	}
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}
{{- if eq .Type 1 }}
func ({{ .ServiceShortAlias }} *{{ .Service }}{{ .ServiceType }}) {{ .Name }}(c context.Context, req *{{ .RequestType }}) (reply *{{ .ReplyType }}, err error) {
    return 
}

{{- else if eq .Type 2 }}
func ({{ .ServiceShortAlias }} *{{ .Service }}{{ .ServiceType }}) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Server) error {
    for {
        req, err := conn.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        
        err = conn.Send(&pb.{{ .Reply }}{})
        if err != nil {
            return err
        }
    }
}

{{- else if eq .Type 3 }}
func ({{ .ServiceShortAlias }} *{{ .Service }}{{ .ServiceType }}) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Server) error {
    for {
        req, err := conn.Recv()
        if err == io.EOF {
            return conn.SendAndClose(&pb.{{ .Reply }}{})
        }
        if err != nil {
            return err
        }
    }
}

{{- else if eq .Type 4 }}
func ({{ .ServiceShortAlias }} *{{ .Service }}{{ .ServiceType }}) {{ .Name }}(req *{{ .RequestType }}, conn pb.{{ .Service }}_{{ .Name }}Server) error {
    for {
        err := conn.Send(&pb.{{ .Reply }}{})
        if err != nil {
            return err
        }
    }
}

{{- end }}
{{- end }}
`

// Service is a proto service.
type Service struct {
	Package     string
	Service     string
	Methods     []*Method
	GoogleEmpty bool

	ServiceType string
	PackageName string

	UseIO      bool
	UseContext bool

	ImportList        base.Import
	PackageDomainList base.PackageDomain
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string
	Request string
	Reply   string

	RequestType string
	ReplyType   string

	// type: unary or stream
	Type entity.MethodType

	ServiceType       string
	ServiceShortAlias string
}

func (s *Service) Alias(m *Method) {
	var pkg string
	var ok bool

	// request
	if m.RequestType, ok = entity.SpecialMap[m.Request]; !ok {
		if m.RequestType, pkg = s.PackageDomainList.ParsePackageInParam(m.Request); pkg != "" {
			s.ImportList = s.ImportList.Special(pkg)
		}
		if !strings.Contains(m.RequestType, ".") {
			m.RequestType = "pb." + m.RequestType
		}
	}

	// reply
	if m.ReplyType, ok = entity.SpecialMap[m.Reply]; !ok {
		if m.ReplyType, pkg = s.PackageDomainList.ParsePackageInParam(m.Reply); pkg != "" {
			s.ImportList = s.ImportList.Special(pkg)
		}
		if !strings.Contains(m.ReplyType, ".") {
			m.ReplyType = "pb." + m.ReplyType
		}
	}
}

func (s *Service) execute() (rst []byte, err error) {

	for _, method := range s.Methods {

		s.Alias(method)

		switch method.Type {
		case entity.UnaryType:
			s.ImportList = s.ImportList.Context()
			if method.Request == entity.AnyPb ||
				method.Reply == entity.AnyPb {
				s.ImportList = s.ImportList.AnyPB()
			}
			if method.Request == entity.EmptyPb ||
				method.Reply == entity.EmptyPb {
				s.ImportList = s.ImportList.EmptyPB()
			}
		case entity.TwoWayStreamsType, entity.RequestStreamsType:
			s.ImportList = s.ImportList.IO()
		case entity.ReturnsStreamsType:
			if method.Request == entity.AnyPb {
				s.ImportList = s.ImportList.AnyPB()
			}
			if method.Request == entity.EmptyPb {
				s.ImportList = s.ImportList.EmptyPB()
			}
		}

		method.ServiceShortAlias = strings.ToLower(s.PackageName[:1])
	}
	s.ImportList = s.ImportList.Show()

	var tmpl *template.Template
	if tmpl, err = template.
		New("service").
		Parse(serviceTemplate); err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, s); err != nil {
		return
	}
	rst = buf.Bytes()
	return
}
