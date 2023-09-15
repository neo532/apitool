package transport

import (
	"context"
	"fmt"
)

/*
* @abstract 传输协议http的一些通用方法
* @mail neo532@126.com
* @date 2022-05-30
 */

type Env string

const (
	EnvDev  Env = "dev"
	EnvTest Env = "test"
	EnvGray Env = "gray"
	EnvProd Env = "prod"
)

func String2Env(str string) (e Env, err error) {
	switch str {
	case string(EnvDev):
	case string(EnvTest):
	case string(EnvProd):
	case string(EnvGray):
	default:
		err = fmt.Errorf("Wrong env[%s], should be dev|test|prod", str)
		return
	}
	e = Env(str)
	return
}

type Logger interface {
	Info(c context.Context, msg string)
	Error(c context.Context, msg string)
}

type LoggerDefault struct {
}

func (l *LoggerDefault) Info(c context.Context, msg string) {
	fmt.Println(fmt.Sprintf("default:%s", msg))
}
func (l *LoggerDefault) Error(c context.Context, msg string) {
	fmt.Println(fmt.Sprintf("default:%s", msg))
}
