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

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_package = "{{.JavaPackage}}";


service {{.Service}} {

    option (google.api.domain) = {
        dev: "http://api.xesv5.com/layout"
        test: "http://api.xesv5.com/layout"
        prod: "http://api.xesv5.com/layout"
        gray: "http://api.xesv5.com/layout"
    };

	rpc Create (CreateRequest) returns (CreateReply){
		option (google.api.http)={
            post: "/session/create" 
            body: "*"
        };
	};
	rpc Get (GetRequest) returns (GetReply){
		option (google.api.http)={
            get: "/work/get" 
        };
	};
	rpc Update (UpdateRequest) returns (UpdateReply){
		option (google.api.http)={
            post: "/work/update" 
            body: "*"
			//retryTimes: "2" // 状态码非200的重试次数
			//retryDuration: "2" // 重试间隔，单位为秒
			//retryMaxDuration: "3" // 重试最大时间，单位为秒
			//timeLimit: "2" // 请求超时时间，单位为秒
			//header: "true" // 是否有自定义header，"true" or "false"
			//contentType: "application/json" // httpHeader的请求类型，eg: application/json
			//function: "userDefined" // 请求执行自定义方法名（同包内）：function userDefined(c context.Context, req rxxx, header hxxx)(req rxxx, header hxxx)
			//respTpl: "origin" // 针对接口resp非标准的{"msg":"ok", "code":0,"data":{}} 的可以自定义解析模板。eg：respTpl:"origin"
        };
	};
}

// Create
message CreateRequest {
    int32 type = 1 [(validate.rules).int32 = {in: 1}];
    string name = 3 [(validate.rules).string = {min_len: 1}];;
}
message CreateReply {
}

// Get
message GetRequest {
    uint64 id = 2 [(validate.rules).uint64 = {lt: 100000000}];
}
message GetReply {
	string url = 1;
    uint32 status = 2;
}

// Update
message UpdateRequest {
    uint64 id = 2 [(validate.rules).uint64 = {lt: 100000000}];
}
message UpdateReply {
	string url = 1;
    uint32 status = 2;
}
`

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
