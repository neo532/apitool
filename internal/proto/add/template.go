package add

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};

import "google/api/annotations.proto";
import "validate/validate.proto";
import "google/api/domain.proto";
import "google/protobuf/empty.proto";

option go_package = "{{.GoPackage}}";

service {{.Service}} {

    option (google.api.domain) = {
		dev: "http://127.0.0.1:8500/demo"
        test: "http://127.0.0.1:8500/demo"
        gray: "http://127.0.0.1:8500/demo"
        prod: "http://127.0.0.1:8500/demo"
    };

	rpc Post (PostRequest) returns (google.protobuf.Empty){
		option (google.api.http)={
            post: "/resource" 
            body: "*"
			needClient: "true"
        };
	};

	rpc Put (PutRequest) returns (google.protobuf.Empty){
		option (google.api.http)={
            put: "/resource" 
            body: "*"
			needClient: "true"
        };
	};

	rpc Get (GetRequest) returns (GetReply){
		option (google.api.http)={
            get: "/resource" 
			needClient: "true"
        };
	};

	rpc GetById (GetByIdRequest) returns (GetByIdReply){
		option (google.api.http)={
            get: "/v1/resource/{id}" 
			needClient: "true"
        };
	};
}

// Post
message PostRequest {
    int32 type = 1 [(validate.rules).int32 = {in: 1}];
    string name = 2 [(validate.rules).string.pattern = "^[a-z0-9]{41}$"];
    string summary = 3 [(validate.rules).string = {min_len: 1}];;
}

// Put
message PutRequest {
    int64 id = 2 [(validate.rules).int64 = {lt: 100000000}];
}

// Get
message GetRequest {
}
message GetReply {
}

// GetById
message GetByIdRequest {
    uint64 id = 1 [(validate.rules).uint64 = {lt: 100000000}];
}
message GetByIdReply {
	string url = 1;
    int32 status = 2;
}

`

func (p *Proto) execute() (b []byte, err error) {

	var tmpl *template.Template
	if tmpl, err = template.
		New("proto").
		Parse(strings.TrimSpace(protoTemplate)); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, p); err != nil {
		return
	}
	b = buf.Bytes()
	return
}
