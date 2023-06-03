package container

import (
	"docker-go/common"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

// RemoveContainer 删除容器
func RemoveContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		logrus.Errorf("get container info, err: %v", err)
	}

	// 只能删除停止状态的容器
	if info.Status != common.Stop {
		logrus.Errorf("can't remove running container")
		return
	}
	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	err = os.RemoveAll(dir)
	if err != nil {
		logrus.Errorf("remove container dir: %s, err: %v", dir, err)
		return
	}
}
