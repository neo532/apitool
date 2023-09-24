package httpclient

import (
	"bytes"
	"html/template"
	"strings"
)

const (
	empty = "google.protobuf.Empty"
)

var RespTplMap = map[string]string{
	"default": `
message {{ .Reply }}Wraper {
    int32 code = 1;
    string message = 2;
    {{ .Reply }} data = 3;
}`,
}

var serviceTemplate = `
{{- /* delete empty line */ -}}
// Code generated by tool. DO NOT EDIT.
// Command : apitool httpclient {{ .ProtoFileName }}
package {{ .PackageName }}

import (
	{{- if .UseContext }}
	"context"
	{{- end }}
	{{- if .UseIO }}
	"io"
	{{- end }}
	{{/* "encoding/json" */}}
	//"net/http"
	{{- if ne .HasImportTime "" }}
	"time"{{- end }}

	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}

	kithttp "github.com/neo532/apitool/transport/http"
	"github.com/neo532/apitool/transport/http/xhttp"
)

type {{ .Service }}XHttpClient struct {
	*kithttp.XClient
	wrapper *kithttp.Wrapper
}

func New{{ .Service }}XHttpClient(clt xhttp.Client) (xclt *{{ .Service }}XHttpClient) {
	xclt = &{{ .Service }}XHttpClient{
		XClient: kithttp.NewXClient(),
		wrapper: kithttp.NewWrapper(clt),
	}

	{{- if ne .DomainsLen 0}}
	domains := map[string]string{ {{ range $env, $domain := .Domains }}
		"{{ $env }}": "{{ $domain }}",
	{{- end }}
	}
	if d, ok := domains[string(clt.Env())]; ok {
		xclt.WithDomain(d)
	}
	{{- end}}
	return
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}
{{- if eq .Type 1 }}
func (s *{{ .Service }}XHttpClient) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty {{ else }}*{{ .Request }}{{ end }}) (resp{{ if eq .Reply $s1 }} *emptypb.Empty{{ else if eq .RespTpl "" }} *{{ .Reply }}{{ else }} *{{ .Reply }}Wraper{{ end }}, err error) {
	{{ if ne .Function "" }}req = {{ .Function }}(ctx, req){{ end }}

	opts := make([]xhttp.Opt, 0, 6)
	opts = append(opts, xhttp.WithHeader(ctx))
	opts = append(opts, xhttp.WithUrl(s.Domain+"{{ .Path }}"))
	opts = append(opts, xhttp.WithMethod("{{ .Method }}")) {{ if ne .TimeLimit "" }}
	opts = append(opts, xhttp.WithTimeLimit({{ .TimeLimit }}*time.Second)){{ end }} {{ if ne .RetryTimes "" }}
	opts = append(opts, xhttp.WithRetryTimes({{ .RetryTimes }})){{ end }} {{ if ne .RetryDuration "" }}
	opts = append(opts, xhttp.WithRetryDuration({{ .RetryDuration }}*time.Second)){{ end }} {{ if ne .RetryMaxDuration "" }}
	opts = append(opts, xhttp.WithRetryMaxDuration({{ .RetryMaxDuration }}*time.Second)){{ end }} {{ if eq .Method "GET" }}
	if ctx, err = xhttp.AppendUrlByStruct(ctx, req); err!=nil {
		return
	}
	{{ else }}

	var b []byte
	if b, err = s.XClient.Codec("{{ .ContentTypeRequest }}").Marshal(req); err != nil {
		return
	}
	opts = append(opts, xhttp.WithBody(b))
	{{ end }}

	{{ if ne .Reply $s1 }}
	respObj := s.wrapper.Call(ctx, s.Domain, opts...)
	var body []byte
	if body, err = respObj.Body(ctx); err!=nil {
		return
	}
	{{ if eq .RespTpl "" }}resp = &{{ .Reply }}{}{{ else }}resp = &{{ .Reply }}Wraper{}{{ end }}
	err = s.XClient.Codec("{{ .ContentTypeResponse }}").Unmarshal(body, resp)
	{{ else }}
	s.wrapper.Call(ctx, s.Domain, opts...)
	{{ end }}
	return 
}

{{- else if eq .Type 2 }}
func (s *{{ .Service }}XHttpClient) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Client) error {
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
func (s *{{ .Service }}XHttpClient) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Client) error {
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
func (s *{{ .Service }}XHttpClient) {{ .Name }}(req {{ if eq .Request $s1 }}*emptypb.Empty
{{ else }}*pb.{{ .Request }}{{ end }}, conn pb.{{ .Service }}_{{ .Name }}Client) error {
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
	Package       string
	PackageName   string
	Service       string
	Methods       []*Method
	GoogleEmpty   bool
	Domains       map[string]string
	DomainsLen    int
	TargetDir     string
	ProtoFileName string
	HasImportTime string

	ServiceLower string

	UseIO      bool
	UseContext bool
}

// Method is a proto method.
type Method struct {
	Service  string
	MService string
	Name     string
	Request  string
	Reply    string

	// type: unary or stream
	Type MethodType

	//xhttp opts
	Path                string
	TimeLimit           string
	Method              string
	RetryTimes          string
	RetryDuration       string
	RetryMaxDuration    string
	Function            string
	ContentTypeRequest  string
	ContentTypeResponse string
	RespTpl             string
	ReqOmitEmpty        string
}

func (s *Service) execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	sPackage := strings.Split(s.Package, "/")
	mService := sPackage[len(sPackage)-2]

	for _, method := range s.Methods {
		if method.ContentTypeRequest == "" {
			method.ContentTypeRequest = "json"
		}
		if method.ContentTypeResponse == "" {
			method.ContentTypeResponse = "json"
		}

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
		method.MService = mService
	}

	s.ServiceLower = strings.ToLower(s.Service)

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
