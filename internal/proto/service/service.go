package service

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/neo532/apitool/internal/base"
	"github.com/neo532/apitool/internal/proto/entity"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CmdService the service command.
var CmdService = &cobra.Command{
	Use:   "service",
	Short: "Generate the proto Service implementations",
	Long:  "Generate the proto Service implementations. Example: apitool service api/xxx.proto -target-dir=internal/service",
	Run:   run,
}
var targetDir string

func init() {
	CmdService.Flags().StringVarP(&targetDir, "target-dir", "t", "", "generate target directory")
}

// TODO gopackage 最好把gomod的包名全路径，不然service不好使
func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: apitool service api/xxx.proto")
		return
	}
	buildHttpClient(cmd, args[0])
}

func buildHttpClient(cmd *cobra.Command, filePath string) {

	var err error

	pb := &Proto{
		// MessageNameMap:    make(map[string]struct{}, 10),
		// Services:          make([]*Service, 0, 10),
		FilePath:          filePath,
		PackageDomainList: base.NewPackageDomainList(),
	}

	var definition *proto.Proto
	definition, err = pb.ReadProtoFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Read proto[%s] has err[%+v]", pb.FilePath, err))
		return
	}

	if targetDir == "" {
		targetDir = filepath.Dir(filePath)
	}
	targetDir = strings.TrimSuffix(targetDir, "/")

	var (
		pkg string
		res []*Service
	)
	proto.Walk(definition,
		proto.WithImport(func(i *proto.Import) {
			pb.PackageDomainList.Add(i.Filename)
		}),
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				pkg = strings.Split(o.Constant.Source, ";")[0]
			}
		}),
		proto.WithService(func(s *proto.Service) {
			ts := strings.Split(targetDir, "/")
			typ := ts[len(ts)-1]
			cs := &Service{
				Package: pkg,
				Service: serviceName(s.Name),

				ServiceType: toUpperCamelCase(typ),
				PackageName: typ,

				ImportList:        base.NewImportList(),
				PackageDomainList: pb.PackageDomainList,
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if !ok {
					continue
				}
				cs.Methods = append(cs.Methods, &Method{
					Service: serviceName(s.Name), Name: serviceName(r.Name), Request: r.RequestType,
					Reply: r.ReturnsType, Type: entity.GetMethodType(r.StreamsRequest, r.StreamsReturns),

					ServiceType: toUpperCamelCase(typ),
				})
			}
			res = append(res, cs)
		}),
	)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exsit\n", targetDir)
		return
	}
	for _, s := range res {
		to := path.Join(targetDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to); !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to)
			continue
		}
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

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}
