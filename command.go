/*
	定义了两个运行命令runCommand，initCommand
*/

package main

import (
	"docker-go/cgroups/subsystem"
	"docker-go/container"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// 创建namespace隔离的容器进程
// 启动容器
var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "docker volume",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container args")
		}
		tty := context.Bool("ti")
		volume := context.String("v")

		res := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		// cmdArray 为容器运行后，执行的第一个命令信息
		// cmdArray[0] 为命令内容，后面的为命令参数
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		Run(cmdArray, tty, res, volume)
		return nil
	},
}

// 初始化容器内容，挂载proc文件系统，运行用户执行程序
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		logrus.Info("init come on")
		return container.RunContainerInitProcess()
	},
}

// 导出容器内容
var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "docker commit a container into image",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Usage: "export image path",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		imageName := context.Args().Get(0)
		imagePath := context.String("c")
		return container.CommitContainer(imageName, imagePath)
	},
}

// logs 命令参数用于查看日志
var logCommand = cli.Command{
	Name:  "logs",
	Usage: "look container log",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		containerName := ctx.Args().Get(0)
		container.LookContainerLog(containerName)

		return nil
	},
}
