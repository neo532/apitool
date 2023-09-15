package upgrade

import (
	"fmt"

	"github.com/neo532/apitool/internal/base"

	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the apitool tools",
	Long:  "Upgrade the apitool tools. Example: apitool upgrade",
	Run:   Run,
}

// Run upgrade the apitool tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
