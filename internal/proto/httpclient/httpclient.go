package httpclient

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/cobra"

	"github.com/neo532/apitool/internal/base"
)

// CmdClient the httpclient command.
var (
	CmdClient = &cobra.Command{
		Use:   "httpclient",
		Short: "Generate the proto Client implementations",
		Long:  "Generate the proto Client implementations. Example: apitool httpclient api/xxx.proto -target-dir=internal/service",
		Run:   run,
	}

	targetDir    string
	verboseKey   = "verbose"
	verboseValue string
)

func init() {
	//CmdClient.Flags().StringVarP(&targetDir, "target-dir", "t", "", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: apitool httpclient api/xxx.proto")
		return
	}
	for _, v := range args {
		if v == verboseKey {
			verboseValue = verboseKey
			break
		}
	}
	for _, v := range args {
		if v == verboseKey {
			continue
		}
		buildHttpClient(cmd, v)
	}
}

func buildHttpClient(cmd *cobra.Command, filePath string) {
	var err error

	pb := &Proto{
		MessageNameMap:    make(map[string]struct{}, 10),
		Services:          make([]*Service, 0, 10),
		FilePath:          filePath,
		CacheTpl:          make(map[string]string, 1),
		PackageDomainList: NewPackageDomainList(),
	}

	//targetDir = filepath.Dir(pb.FilePath)
	targetDir = filepath.Dir(filePath)

	var definition *proto.Proto
	definition, err = pb.ReadProtoFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Read proto[%s] has err[%+v]", pb.FilePath, err))
		return
	}

	var needClient bool
	proto.Walk(definition,
		proto.WithImport(func(i *proto.Import) {
			pb.PackageDomainList.Add(i.Filename)
		}),
		proto.WithOption(func(o *proto.Option) {
			switch o.Name {
			case "go_package":
				p := strings.Split(o.Constant.Source, ";")
				pb.PackageName = filepath.Base(p[len(p)-1])
			}
		}),
		proto.WithMessage(func(m *proto.Message) {
			// exist messageName
			if m.Position.Column == 1 {
				pb.MessageNameMap[m.Name] = struct{}{}
			}
		}),
		proto.WithService(func(s *proto.Service) {
			cs := &Service{
				TargetDir:     targetDir,
				ProtoFileName: pb.FilePath,
				//Package:       pb.Package,
				PackageName:       pb.PackageName,
				Service:           s.Name,
				Domains:           make(map[string]string, 2),
				ImportList:        NewImportList(),
				PackageDomainList: pb.PackageDomainList,
			}
			for _, e := range s.Elements {

				if r, ok := e.(*proto.Option); ok {
					if r.Name == "(google.api.domain)" {
						for _, v := range r.AggregatedConstants {
							cs.Domains[v.Name] = v.Source
						}
						cs.DomainsLen = len(cs.Domains)
						continue
					}
				}

				if r, ok := e.(*proto.RPC); ok {
					method := &Method{
						Service: s.Name,
						Name:    r.Name,
						Request: r.RequestType,
						Reply:   r.ReturnsType,
						Type:    getMethodType(r.StreamsRequest, r.StreamsReturns),
					}
					// http parameter
					if len(r.Elements) > 0 {
						packageHttpParameter2Method(method, r.Options, cs)
					}
					cs.Methods = append(cs.Methods, method)
					continue
				}

			}
			pb.Services = append(pb.Services, cs)
			if cs.NeedClient == true {
				needClient = true
			}
		}),
	)
	if needClient == false {
		return
	}
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exsits\n", targetDir)
		return
	}
	for _, s := range pb.Services {

		// ._http.client.go 结尾的文件
		to := path.Join(targetDir, strings.ToLower(s.Service)+".httpclient.pb.go")
		//if _, err := os.Stat(to); !os.IsNotExist(err) {
		//fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to)
		//continue
		//}
		b, err := s.execute()
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(to, b, 0o644); err != nil {
			log.Fatal(err)
		}
		if verboseValue != "" {
			fmt.Println(to)
		}

		// .proto suffix
		for _, m := range s.Methods {

			if pb.IsNeedAddWrapper(m) == true {
				var tpl string
				if tpl, err = pb.GetTpl(m); err != nil {
					fmt.Fprintln(os.Stderr, fmt.Sprintf("BuildHttpClient %s has error[%+v]", pb.FilePath, err))
					return
				}
				pb.NewWraper(m, tpl)
			}
		}
		pb.AppendWraper()

		// let wraper append to file for PB generator
		base.Run("apitool", "pbstruct", pb.FilePath, verboseValue)
	}
}

// http parameter
func packageHttpParameter2Method(method *Method, opts []*proto.Option, cs *Service) {
	for _, o := range opts {
		for _, c := range o.AggregatedConstants {
			switch c.Name {
			case "contentType":
				method.ContentType = c.Literal.Source
				cs.NeedClient = true
			case "contentTypeResponse":
				method.ContentTypeResponse = c.Literal.Source
				cs.NeedClient = true
			case "retryTimes":
				method.RetryTimes = c.Literal.Source
				cs.NeedClient = true
			case "retryDuration":
				cs.ImportList = cs.ImportList.Time()
				method.RetryDuration = c.Literal.Source
				cs.NeedClient = true
			case "retryMaxDuration":
				cs.ImportList = cs.ImportList.Time()
				method.RetryMaxDuration = c.Literal.Source
				cs.NeedClient = true
			case "timeLimit":
				cs.ImportList = cs.ImportList.Time()
				method.TimeLimit = c.Literal.Source
				cs.NeedClient = true
			case "get", "post", "put", "delete", "head", "patch", "options", "trace", "connect":
				method.Method = strings.ToUpper(c.Name)
				method.Path = c.Literal.Source
			case "respTpl":
				tmp := strings.SplitN(c.Literal.Source, ",", 2)
				method.RespTpl = tmp[0]
				method.RespTplDataName = "Data"
				if len(tmp) >= 2 {
					method.RespTplDataName = tmp[1]
				}
				cs.NeedClient = true
			case "requestEncoder":
				method.RequestEncoder = c.Literal.Source
				cs.NeedClient = true
			case "responseDecoder":
				method.ResponseDecoder = c.Literal.Source
				cs.NeedClient = true
			case "errorDecoder":
				method.ErrorDecoder = c.Literal.Source
				cs.NeedClient = true
			case "insecureSkipVerify":
				method.InsecureSkipVerify = c.Literal.Source
				cs.NeedClient = true
			case "caCertFile":
				method.CaCertFile = c.Literal.Source
				cs.NeedClient = true
			case "needClient":
				cs.NeedClient = true
			case "certFile":
				cs.NeedClient = true
				crt := strings.Split(c.Literal.Source, ",")
				method.CertFileCrt = crt[0]
				if len(crt) == 2 {
					method.CertFileKey = crt[1]
				} else {
					log.Fatal("Please input valid certFile,eg:./configs/xx.crt,./configs/xx.key")
				}
			}
		}
	}
	return
}

func getMethodType(streamsRequest, streamsReturns bool) MethodType {
	if !streamsRequest && !streamsReturns {
		return unaryType
	} else if streamsRequest && streamsReturns {
		return twoWayStreamsType
	} else if streamsRequest {
		return requestStreamsType
	} else if streamsReturns {
		return returnsStreamsType
	}
	return unaryType
}
