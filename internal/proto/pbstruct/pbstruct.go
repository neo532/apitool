package pbstruct

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neo532/apitool/internal/base"

	"github.com/spf13/cobra"
)

// CmdStruct represents the source command.
var (
	CmdStruct = &cobra.Command{
		Use:   "pbstruct",
		Short: "Generate the proto server code",
		Long:  "Generate the proto server code. Example: apitool struct helloworld.proto",
		Run:   run,
	}
	verboseValue string
	verboseKey   = "verbose"
)

var protoPath string

func init() {
	if protoPath = os.Getenv("APIPATH_PROTO_PATH"); protoPath == "" {
		protoPath = "./third_party"
	}
	CmdStruct.Flags().StringVarP(&protoPath, "proto_path", "p", protoPath, "proto path")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	for _, v := range args {
		if v == verboseKey {
			verboseValue = verboseKey
		}
	}
	var (
		err   error
		proto = strings.TrimSpace(args[0])
	)
	if err = look(
		"protoc-gen-go",
		// "protoc-gen-go-grpc",
		// "protoc-gen-go-http",
		// "protoc-gen-go-errors",
		//"protoc-gen-openapi",
	); err != nil {
		// update the apitool plugins
		if _, err = base.Run("apitool", "upgrade"); err != nil {
			fmt.Println(err)
			return
		}
	}
	if strings.HasSuffix(proto, ".proto") {
		err = generate(proto, args)
	} else {
		err = walk(proto, args)
	}
	if err != nil {
		fmt.Println(runtime.Caller(0))
		log.Println(err)
		fmt.Println(fmt.Sprintf("err:\t%+v", err))
		fmt.Println(err)
		fmt.Println(runtime.Caller(0))
		return
	}
	pbGoPath := strings.Replace(proto, ".proto", ".pb.go", 1)
	if _, err = base.Run("protoc-go-inject-tag", "-input="+pbGoPath, verboseValue); err != nil {
		fmt.Println(err)
	}
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" || strings.HasPrefix(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

// generate is used to execute the generate command for the specified proto file
func generate(proto string, args []string) (err error) {
	input := []string{
		"--proto_path=.",
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}
	inputExt := []string{
		"--proto_path=" + base.ModPath(),
		"--proto_path=" + filepath.Join(base.ModPath(), "third_party"),
		//"--gofast_out=paths=source_relative:.",
		"--go_out=paths=source_relative:.",
		// "--go-grpc_out=paths=source_relative:.",
		// "--go-http_out=paths=source_relative:.",
		// "--go-errors_out=paths=source_relative:.",
		// "--openapi_out=paths=source_relative:.",
	}
	input = append(input, inputExt...)
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}

	if _, err = base.Run("protoc", input...); err != nil {
		return
	}
	if verboseValue != "" {
		fmt.Printf("proto: %s\n", proto)
	}
	return
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
