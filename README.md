# Instruction

## Install
```
go install github.com/neo532/apitool@master
cd apitool
make init
```

## Template
save this in {filePath}.tpl,and write {filePath} into rpc.option.RespTpl's Value.
```
message {{ .ReplyName }} { 
    int32 code = 1;
    string message = 2;
    {{ .ReplyType }} data = 3;
}
```

## File define
{path}/{filePath}.api.proto

### Init a proto file
```
apitool add api/{example1}/{pkg1}/{pkg1}.api.proto
vim api/{example1}/{pkg1}/{pkg1}.api.proto
```

### Generate a httpclient's structs by a proto file.
```
apitool pbstruct api/{example1}/{pkg1}/{pkg1}.api.proto
```

### Generate a httpclient by a proto file.
```
apitool httpclient api/{example1}/{pkg1}/{pkg1}.api.proto
```

### Generate a service by a proto file.
```
apitool service api/{demo}/{pkg}/{pkg}.api.proto -t api/{demo}/{pkg}
```
