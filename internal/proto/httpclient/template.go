package httpclient

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/neo532/apitool/internal/base"
	"github.com/neo532/apitool/internal/proto/entity"
	"github.com/neo532/apitool/transport/http/xhttp"
)

var serviceTemplate = `
{{- /* delete empty line */ -}}
// Code generated by tool. DO NOT EDIT.
// Command : apitool httpclient {{ .ProtoFileName }}
package {{ .PackageName }}

import (
	{{range $key,$value := .ImportList }}{{ $value }}"{{ $key }}"
	{{ end }}
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
		xclt.Domain = d
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
	{{ if and (ne .CertFileCrt "") (ne .CertFileKey "") }}
	var certFile xhttp.Opt
	if certFile, err = xhttp.WithCertFile("{{ .CertFileCrt }}", "{{ .CertFileKey }}"); err != nil {
		return
	}
	opts = append(opts, certFile)
	{{ end }} 
	{{ if ne .CaCertFile "" }}
	var caCertFile xhttp.Opt
	if caCertFile, err = xhttp.WithCaCertFile("{{ .CaCertFile }}"); err != nil {
		return
	}
	opts = append(opts, caCertFile)
	{{ end }} 
	{{- if .InsecureSkipVerify }}
	opts = append(opts, xhttp.WithInsecureSkipVerify({{ .InsecureSkipVerify }}))
	{{ end }}
	{{- if .HasQueryArgs }}
	if ctx, err = xhttp.AppendUrlByStruct(ctx, req); err != nil {
		return
	}
	req = &{{ .RequestType }}{}
	{{ end }}
	{{ if eq .ReplyType .AnyTypeKey }}
	resp = &{{ .ReplyType }}{}
	err = xhttp.New(s.Client, opts...).Do(ctx, req, resp)
	{{ else if  eq .ReplyTypeWrapper "" }}
	resp = &{{ .ReplyType }}{}
	err = xhttp.New(s.Client, opts...).Do(ctx, req, resp)
	{{ else }}
	rst := &{{ .ReplyTypeWrapper }}{}
	err = xhttp.New(s.Client, opts...).Do(ctx, req, rst)
	if rst != nil {
		resp = rst.{{ .RespTplDataName }}
	}
	{{ end }} 
	return 
}
{{- end }}
{{- end }}
`

// Service is a proto service.
type Service struct {
	//Package     string
	PackageName string
	Service     string
	Methods     []*Method

	NeedClient bool

	Domains       map[string]string
	DomainsLen    int
	TargetDir     string
	ProtoFileName string

	ServiceLower string

	EmptyHas bool
	AnyHas   bool
	//UseIO    bool
	//UseContext    bool
	//HasImportTime string

	ImportList        base.Import
	PackageDomainList base.PackageDomain
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string

	// google.protobuf.Empty | google.protobuf.Any | github.com/a.B | B
	Request string
	Reply   string

	// emptypb.Empty | anypb.Any | a.B | B
	RequestType string
	ReplyType   string

	// EmptyWrapper | AnyWrapper | BWrapper
	ReplyTypeWrapper string // anypb & wrappper at same time

	HasQueryArgs bool

	AnyTypeKey string

	// type: unary or stream
	Type entity.MethodType

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
	RespTplDataName     string

	RequestEncoder  string
	ResponseDecoder string
	ErrorDecoder    string

	CertFileCrt        string
	CertFileKey        string
	InsecureSkipVerify string
	CaCertFile         string
}

func (s *Service) Alias(m *Method) {
	var pkg string
	var ok bool

	// request
	if m.RequestType, ok = entity.SpecialMap[m.Request]; !ok {
		if m.RequestType, pkg = s.PackageDomainList.ParsePackageInParam(m.Request); pkg != "" {
			s.ImportList = s.ImportList.Special(pkg)
		}
	}

	// reply
	if m.ReplyType, ok = entity.SpecialMap[m.Reply]; !ok {
		if m.ReplyType, pkg = s.PackageDomainList.ParsePackageInParam(m.Reply); pkg != "" {
			s.ImportList = s.ImportList.Special(pkg)
		}
	}

	if m.RespTpl != "" {
		reply := strings.TrimSuffix(m.ReplyType, entity.Wrapper)
		rs := strings.Split(reply, ".")
		m.ReplyTypeWrapper = rs[len(rs)-1] + entity.Wrapper
	}
}

func (s *Service) execute() (rst []byte, err error) {

	for _, method := range s.Methods {

		method.AnyTypeKey = entity.AnyType

		if xhttp.HasBody(method.Method) == false {
			method.HasQueryArgs = true
		}

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
	}

	s.ImportList = s.ImportList.Show()

	s.ServiceLower = strings.ToLower(s.Service)

	var tmpl *template.Template
	if tmpl, err = template.
		New("service").
		Parse(serviceTemplate); err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, s); err != nil {
		return
	}
	rst = buf.Bytes()
	return
}
