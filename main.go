/*
	程序的入口，主要是接收命令行参数
	接收命令行参数处理使用的第三方工具包为 cli
	日志打印采用的 logrus
*/

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `docker-go`

func main() {
	app := cli.NewApp()
	app.Name = "docker-go"
	app.Usage = usage

	app.Commands = []cli.Command{
		runCommand,
		initCommand,
		commitCommand,
	}

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
