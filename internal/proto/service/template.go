package service

import (
	"bytes"
	"html/template"
	"strings"
)

//nolint:lll
var serviceTemplate = `
{{- /* delete empty line */ -}}
package {{ .PackageName }}

import (
    {{- if .UseContext }}
    "context"
    {{- end }}
    {{- if .UseIO }}
    "io"
    {{- end }}

    pb "{{ .Package }}"
    {{- if .GoogleEmpty }}
    "google.golang.org/protobuf/types/known/emptypb"
    {{- end }}
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
func ({{ .ServiceShortAlias }} *{{ .Service }}{{ .ServiceType }}) {{ .Name }}(c context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) (reply {{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end }}, err error) {
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
func ({{ .ServiceShortAlias }} *{{ .Service }}{{ .ServiceType }}) {{ .Name }}(req {{ if eq .Request $s1 }}*emptypb.Empty
{{ else }}*pb.{{ .Request }}{{ end }}, conn pb.{{ .Service }}_{{ .Name }}Server) error {
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

type MethodType uint8

const (
	unaryType          MethodType = 1
	twoWayStreamsType  MethodType = 2
	requestStreamsType MethodType = 3
	returnsStreamsType MethodType = 4
)

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
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string
	Request string
	Reply   string

	// type: unary or stream
	Type MethodType

	ServiceType       string
	ServiceShortAlias string
}

func (s *Service) execute() ([]byte, error) {
	const empty = "google.protobuf.Empty"
	buf := new(bytes.Buffer)
	for _, method := range s.Methods {
		if (method.Type == unaryType && (method.Request == empty || method.Reply == empty)) ||
			(method.Type == returnsStreamsType && method.Request == empty) {
			s.GoogleEmpty = true
		}
		if method.Type == twoWayStreamsType || method.Type == requestStreamsType {
			s.UseIO = true
		}
		if method.Type == unaryType {
			s.UseContext = true
		}
		method.ServiceShortAlias = strings.ToLower(s.PackageName[:1])
	}
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
