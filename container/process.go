/*
	基本上对 docker 初始化要做的事情都放在了这个文件中
	主要是启动一个容器，然后对该容器做一些资源限制
	这里需要关注的是 container.NewParentProcess(tty)，它会给我们返回一个被 namespace 隔离的进程
	写日志，将输出到控制台的内容输出到指定文件
*/

package container

import (
	"docker-go/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
)

// NewParentProcess 创建一个会隔离namespace进程的Comand
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, _ := os.Pipe()
	// 调用自身，传入 init 参数， 也就是执行initComand
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		// 指定交换终端，日志直接输出到控制台
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 创建日志存放目录
		logDir := path.Join(common.DefaultContainerInfoPath, "container")
		if _, err := os.Stat(logDir); err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				logrus.Errorf("mkdir container log, err: %v", err)
			}
		}
		// 创建日志文件
		logFileName := path.Join(logDir, common.ContainerLogFileName)
		file, err := os.Create(logFileName)
		if err != nil {
			logrus.Errorf("create log file, err: %v", err)
		}
		// 将cmd的输出流改到日志文件中
		cmd.Stdout = file
	}
	cmd.ExtraFiles = []*os.File{
		readPipe,
	}

	err := NewWorkSpace(common.RootPath, common.MntPath, volume)
	if err != nil {
		logrus.Errorf("new work space, err: %v", err)
	}

	// 指定容器初始化后的工作目录
	cmd.Dir = common.MntPath
	return cmd, writePipe
}

// LookContainerLog 用于查看容器内日志信息
func LookContainerLog(containerName string) {
	// 打开文件
	logFileName := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerLogFileName)
	file, err := os.Open(logFileName)
	if err != nil {
		logrus.Errorf("open log file, path: %s, err: %v", logFileName, err)
	}
	// 读文件
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read log file, err: %v", err)
	}
	_, _ = fmt.Fprint(os.Stdout, string(bs))
}
