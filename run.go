/*
	run 命令主要就是启动一个容器，然后对该进程设置隔离，
	init 是 run 命令中调用的，不是我们自身通过命令行调用的，
	这里我们主要关注 Run(cmdArray, tty, res)函数即可，它接收我们传递过来的参数，
	tty 表示是否前台运行，对应docker 的 -ti命令
*/

package main

import (
	"docker-go/cgroups"
	"docker-go/cgroups/subsystem"
	"docker-go/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(cmdArray []string, tty bool, res *subsystem.ResourceConfig, containerName, imageName, volume, net string, envs, ports []string) {
	containerID := container.GetContainerID(10)
	if containerName == "" {
		containerName = containerID
	}
	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envs)
	if parent == nil {
		logrus.Errorf("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("parent start failed, err %v", err)
		return
	}
	// 记录容器信息
	err := container.RecordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerID)
	if err != nil {
		logrus.Errorf("record container info, err: %v", err)
	}
	// 添加资源限制
	cgroupManager := cgroups.NewCGroupManager("docker-go")
	// 删除资源限制
	defer cgroupManager.Destroy()
	// 设置资源限制
	cgroupManager.Set(res)
	// 将容器进程，加入到各个subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)
	// 设置初始化命令
	sendInitCommand(cmdArray, writePipe)
	// 等待父进程结束
	if tty {
		// 等待父进程结束
		err := parent.Wait()
		if err != nil {
			logrus.Errorf("parent wait, err: %v", err)
		}
		// 删除容器工作空间
		err = container.DeleteWorkSpace(containerName, volume)
		if err != nil {
			logrus.Errorf("delete work space, err: %v", err)
		}
		// 删除容器信息
		container.DeleteContainerInfo(containerName)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
