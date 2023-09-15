package grpcclient

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/cobra"
)

// CmdClient the grpcclient command.
var CmdClient = &cobra.Command{
	Use:   "grpcclient",
	Short: "Generate the proto Client implementations",
	Long:  "Generate the proto Client implementations. Example: apitool httpclient api/xxx.proto -target-dir=internal/service",
	Run:   run,
}
var targetDir string

func init() {
	CmdClient.Flags().StringVarP(&targetDir, "target-dir", "t", "", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: apitool httpclient api/xxx.proto")
		return
	}
	protoFileName := args[0]
	reader, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	if targetDir == "" {
		targetDir = filepath.Dir(protoFileName)
	}

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	var (
		pkg string
		res []*Service
	)
	proto.Walk(definition,
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				pkg = strings.Split(o.Constant.Source, ";")[0]
			}
		}),
		proto.WithService(func(s *proto.Service) {
			cs := &Service{
				Package:       pkg,
				Service:       s.Name,
				TargetDir:     targetDir,
				ProtoFileName: protoFileName,
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if ok {
					cs.Methods = append(cs.Methods, &Method{
						Service: s.Name, Name: r.Name, Request: r.RequestType,
						Reply: r.ReturnsType, Type: getMethodType(r.StreamsRequest, r.StreamsReturns),
					})
				}
			}
			res = append(res, cs)
		}),
	)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exsits\n", targetDir)
		return
	}
	for _, s := range res {
		to := path.Join(targetDir, strings.ToLower(s.Service)+"_grpc.client.go")
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
		fmt.Println(to)
	}
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
