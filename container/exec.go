package container

import (
	"docker-go/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// ExecContainer 从新进入容器
// 通过设置环境变量的方式，让C语言写的程序真正执行
// 通过 setns 的系统调用，重新进入到指定的 PID 的 namespace 中
func ExecContainer(containerName string, cmdArray []string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		logrus.Errorf("get container info, err: %v", err)
	}
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = os.Setenv(common.EnvExecPid, info.Pid)
	err = os.Setenv(common.EnvExecCmd, strings.Join(cmdArray, " "))
	// 设置环境变量
	envs := getEnvsByPid(info.Pid)
	cmd.Env = append(os.Environ(), envs...)

	if err = cmd.Run(); err != nil {
		logrus.Errorf("exec cmd run, err: %v", err)
	}
}

// 通过 pid 获取进程环境变量
func getEnvsByPid(pid string) []string {
	envFilePath := fmt.Sprintf("/proc/%s/environ", pid)
	file, err := os.Open(envFilePath)
	if err != nil {
		logrus.Errorf("open env file, path: %s, err: %v", envFilePath, err)
		return nil
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read env file, err: %v", err)
	}

	return strings.Split(string(bs), "\u0000")
}
