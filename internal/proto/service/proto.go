package service

import (
	"log"
	"os"

	"github.com/emicklei/proto"
	"github.com/neo532/apitool/internal/base"
)

type Proto struct {
	FilePath string
	//PackageName string
	////Package        string
	//MessageNameMap map[string]struct{} // exists message name
	//Services       []*Service

	PackageDomainList base.PackageDomain
}

func (pb *Proto) ReadProtoFile() (definition *proto.Proto, err error) {
	reader, err := os.Open(pb.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err = parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	return
}
