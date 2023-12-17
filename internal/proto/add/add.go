package add

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/neo532/apitool/internal/base"
	"github.com/spf13/cobra"
)

// CmdAdd represents the add command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a proto API template",
	Long:  "Add a proto API template. Example: apitool add ./transport/http/proto/example.api.proto",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {

	input := args[0]
	n := strings.LastIndex(input, "/")
	if n == -1 {
		fmt.Println("The proto path needs to be hierarchical.")
		return
	}

	pb := &Proto{
		Path:     filepath.Dir(input),
		FileName: filepath.Base(input),
	}
	pb.Package = filepath.Base(pb.Path)
	pb.Service = base.UpperFirstChar(strings.SplitN(pb.Package, ".", 2)[0])
	pb.GoPackage = "/" + strings.TrimLeft(pb.Path, "./")
	if modName, err := base.ModuleName("go.mod"); err == nil {
		pb.GoPackage = modName + pb.GoPackage
	}

	if err := pb.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}
