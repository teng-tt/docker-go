// 提交容器,打包容器镜像为tar

package container

import (
	"docker-go/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"path"
)

func CommitContainer(imageName, imagePath string) error {
	if imagePath == "" {
		imagePath = common.RootPath
	}
	imageTar := path.Join(imagePath, fmt.Sprintf("%s.tar", imageName))
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", common.MntPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar container image, file name: %s, err: %v", imageTar, err)
		return err
	}
	return nil
}
