package httpclient

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/neo532/apitool/transport/http/xhttp"
)

const (
	emptyPb      = "google.protobuf.Empty"
	emptyVarName = "Empty"
	emptyType    = "emptypb.Empty"

	anyPb      = "google.protobuf.Any"
	anyVarName = "Any"
	anyType    = "anypb.Any"

	wrapper = "Wrapper"
)

var serviceTemplate = `
{{- /* delete empty line */ -}}
// Code generated by tool. DO NOT EDIT.
// Command : apitool httpclient {{ .ProtoFileName }}
package {{ .PackageName }}

import (
	{{- if .UseContext }}
	"context"{{- end }}
	{{- if .UseIO }}
	"io"{{- end }}
	{{- if ne .HasImportTime "" }}
	"time"{{- end }}
	{{- if .EmptyHas }}
	"google.golang.org/protobuf/types/known/emptypb"{{- end }}
	{{- if .AnyHas }}
	"google.golang.org/protobuf/types/known/anypb"{{- end }}

	kithttp "github.com/neo532/apitool/transport/http"
	"github.com/neo532/apitool/transport/http/xhttp"
	"github.com/neo532/apitool/transport/http/xhttp/client"
)

type {{ .Service }}XHttpClient struct {
	*kithttp.XClient
}

func New{{ .Service }}XHttpClient(clt client.Client) (xclt *{{ .Service }}XHttpClient) {
	xclt = &{{ .Service }}XHttpClient{
		XClient: kithttp.NewXClient(clt),
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

{{ range .Methods }}
{{- if eq .Type 1 }}
func (s *{{ .Service }}XHttpClient) {{ .Name }}(ctx context.Context, req *{{ .RequestType }}) (resp *{{ .ReplyType }}, err error) {
	opts := make([]xhttp.Opt, 0, 6)
	opts = append(opts, xhttp.WithUrl(s.Domain+"{{ .Path }}"))
	opts = append(opts, xhttp.WithMethod("{{ .Method }}"))
	{{ if ne .TimeLimit "" }}opts = append(opts, xhttp.WithTimeLimit({{ .TimeLimit }}*time.Second)){{ end }} 
	{{ if ne .RetryTimes "" }}opts = append(opts, xhttp.WithRetryTimes({{ .RetryTimes }})){{ end }}
	{{ if ne .RetryDuration "" }}opts = append(opts, xhttp.WithRetryDuration({{ .RetryDuration }}*time.Second)){{ end }}
	{{ if ne .RetryMaxDuration "" }}opts = append(opts, xhttp.WithRetryMaxDuration({{ .RetryMaxDuration }}*time.Second)){{ end }} 
	{{ if ne .ContentType "" }}opts = append(opts, xhttp.WithContentType({{ .ContentType }})){{ end }} 
	{{ if ne .ContentTypeResponse "" }}opts = append(opts, xhttp.WithContentTypeResponse({{ .ContentTypeResponse }})){{ end }} 
	{{ if ne .RequestEncoder "" }}opts = append(opts, xhttp.WithRequestEncoder({{ .RequestEncoder }})){{ else }}
	if s.RequestEncoder != nil {
		opts = append(opts, xhttp.WithRequestEncoder(s.RequestEncoder))
	}
	{{ end }} 
	{{ if ne .ResponseDecoder "" }}opts = append(opts, xhttp.WithResponseDecoder({{ .ResponseDecoder }})){{ else }}
	if s.ResponseDecoder != nil {
		opts = append(opts, xhttp.WithResponseDecoder(s.ResponseDecoder))
	}
	{{ end }} 
	{{ if ne .ErrorDecoder "" }}opts = append(opts, xhttp.WithErrorDecoder({{ .ErrorDecoder }})){{ else }}
	if s.ErrorDecoder != nil {
		opts = append(opts, xhttp.WithErrorDecoder(s.ErrorDecoder))
	}
	{{ end }} 
	{{- if .HasQueryArgs }}
	if ctx, err = xhttp.AppendUrlByStruct(ctx, req); err != nil {
		return
	}
	req = &{{ .RequestType }}{}
	{{ end }}
	resp = &{{ .ReplyType }}{}
	err = xhttp.New(s.Client, opts...).Do(ctx, req, resp)
	return 
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
	PackageName string
	Service     string
	Methods     []*Method

	EmptyHas bool
	AnyHas   bool

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

	Request string
	Reply   string

	RequestName string
	RequestType string
	ReplyName   string
	ReplyType   string

	HasQueryArgs bool

	// type: unary or stream
	Type MethodType

	//xhttp opts
	Path                string
	TimeLimit           string
	Method              string
	RetryTimes          string
	RetryDuration       string
	RetryMaxDuration    string
	ContentType         string
	ContentTypeResponse string
	RespTpl             string

	RequestEncoder  string
	ResponseDecoder string
	ErrorDecoder    string
}

func FmtNameType(i string) (t, n, pb string) {
	if i == emptyPb {
		return emptyType, emptyVarName, emptyPb
	}
	if i == anyPb {
		return anyType, anyVarName, anyPb
	}
	return i, i, i
}

func (s *Service) execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	sPackage := strings.Split(s.Package, "/")
	mService := sPackage[len(sPackage)-2]

	for _, method := range s.Methods {
		if xhttp.HasBody(method.Method) == false {
			method.HasQueryArgs = true
		}

		method.RequestType, method.RequestName, _ = FmtNameType(method.Request)
		method.ReplyType, method.ReplyName, _ = FmtWraperName(method)

		switch method.Type {
		case unaryType:
			s.UseContext = true
			hasWrapper := IsAddWraper(method.ReplyName)
			if method.Request == anyPb ||
				(!hasWrapper && method.Reply == anyPb) {
				s.AnyHas = true
				break
			}
			if method.Request == emptyPb ||
				(!hasWrapper && method.Reply == emptyPb) {
				s.EmptyHas = true
				break
			}
		case twoWayStreamsType, requestStreamsType:
			s.UseIO = true
		case returnsStreamsType:
			if method.Request == anyPb {
				s.AnyHas = true
				break
			}
			if method.Request == emptyPb {
				s.EmptyHas = true
				break
			}
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
