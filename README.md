# 对外API

## 功能
该项目内容核心文件.proto，其它文件全为.proto文件生成.
* HTTP路由表
* 接口参数验证
* 参数结构体
* 错误码定义
* RPC客户端
* HTTP客户端



## 安装工具
```
go install github.com/neo532/apitool@master
cd apitool
make init
```



## 创建 & 生成
{微服务名}/{包名}/{文件名}.proto
{包名}={文件名}

### 创建proto模板 & 书写
```
apitool add api/{example1}/{pkg1}/{pkg1}.api.proto
vim api/{example1}/{pkg1}/{pkg1}.api.proto
```

### 生成API服务代码+pb结构文件
```
apitool pbstruct api/{example1}/{pkg1}/{pkg1}.api.proto
```

### 生成HTTP客户端文件
```
// 不填路径就默认和proto文件同包
apitool httpclient api/{example1}/{pkg1}/{pkg1}.api.proto
// -t 路径，以项目根路径开始
apitool httpclient api/{example1}/{pkg1}/{pkg1}.api.proto -t api/{example1}/{pkg1}
```


### 生成GRPC客户端文件
```
apitool grpcclient api/{example1}/{pkg1}/{pkg1}.api.proto
```


### 编写错误码proto文件后，在api根目录执行，生成错误码文件
```
make errors
```

### openapi.yaml 文件可作用swagger-ui页面展示
```
// 该命令可在根目录生成所有api的openapi.yaml文件
make api
```
