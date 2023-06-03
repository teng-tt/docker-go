package container

import (
	"docker-go/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

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
