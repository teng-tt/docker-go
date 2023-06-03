/*
	用于镜像空间的挂载删除，使用aufs联合文件系统
	简单来讲就三步，第一步创建只读层，第二步创建读写层，第三步将两者挂载到同一个文件夹下
	具体实现也比较简单，就是创建文件夹，删除文件夹
*/

package container

import (
	"docker-go/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

// NewWorkSpace 创建容器运行时目录
func NewWorkSpace(rootPath, mntPath, volume string) error {
	// 1.创建只读层
	err := createReadOnlyLayer(rootPath)
	if err != nil {
		logrus.Errorf("create read only layer, err: %v", err)
		return err
	}
	// 2. 创建读写层
	err = createWriteLayer(rootPath)
	if err != nil {
		logrus.Errorf("create write layer, err: %v", err)
		return err
	}
	// 3. 创建挂载点，将只读层和读写层挂载到指定的位置
	err = CreateMountPoint(rootPath, mntPath)
	if err != nil {
		logrus.Errorf("create mount point, err: %v", err)
		return err
	}

	// 4. 设置宿主机与容器文件映射
	mountVolume(rootPath, mntPath, volume)

	return nil
}

// 创建只读层
func createReadOnlyLayer(rootPath string) error {
	// 创建只读层目录
	busyBoxPath := path.Join(rootPath, common.BusyBox)
	_, err := os.Stat(busyBoxPath)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(busyBoxPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir busybox, err: %v", err)
			return err
		}
	}

	// 解压 busybox.tar
	busyBoxTarPath := path.Join(rootPath, common.BusyBoxTar)
	if _, err = exec.Command("tar", "-xvf", busyBoxTarPath, "-C", busyBoxPath).CombinedOutput(); err != nil {
		logrus.Errorf("tar busybox.tar, err: %v", err)
		return err
	}

	return nil
}

// 创建读写层
func createWriteLayer(rootPath string) error {
	writeLayerPath := path.Join(rootPath, common.WriteLayer)
	_, err := os.Stat(writeLayerPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(writeLayerPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir write layer, err: %v", err)
			return err
		}
	}

	return nil
}

// CreateMountPoint 创建挂载点
func CreateMountPoint(rootPath, mntPath string) error {
	_, err := os.Stat(mntPath)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(mntPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir mnt path, err: %v", err)
			return err
		}
	}

	dirs := fmt.Sprintf("dir=%s%s:%s%s", rootPath, common.WriteLayer, rootPath, common.BusyBox)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntPath)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mnt cmd run, err: %v", err)
		return err
	}

	return nil
}

// 创建挂载卷
func mountVolume(rootPath, mntPath, volume string) {
	if volume != "" {
		volumes := strings.Split(volume, ":")
		if len(volumes) > 1 {
			// 创建宿主机中文件路径
			parentPath := volumes[0]
			if _, err := os.Stat(parentPath); err != nil && os.IsNotExist(err) {
				if err = os.MkdirAll(parentPath, os.ModePerm); err != nil {
					logrus.Errorf("mkdir parent path: %s, err: %v", parentPath, err)
				}
			}

			// 创建容器内挂载点
			containerPath := volumes[1]
			containerVolumePath := path.Join(mntPath, containerPath)
			if _, err := os.Stat(containerVolumePath); err != nil && os.IsNotExist(err) {
				if err = os.MkdirAll(containerVolumePath, os.ModePerm); err != nil {
					logrus.Errorf("mkdir volume path path: %s, err: %v", containerVolumePath, err)
				}
			}

			// 把宿主机文件目录挂载到容器挂载点中
			dirs := fmt.Sprintf("dirs=%s", parentPath)
			cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				logrus.Errorf("mount cmd run, err: %v", err)
			}
		}
	}
}

// DeleteWorkSpace 删除容器工作空间
func DeleteWorkSpace(rootPath, mntPath, volume string) error {
	// 1. 卸载挂载点
	err := unMountPoint(mntPath)
	if err != nil {
		return err
	}

	// 2. 删除读写层
	err = deleteWriteLayer(rootPath)
	if err != nil {
		return err
	}

	// 3. 删除宿主机与文件系统映射
	deleteVolume(mntPath, volume)
	return nil
}

// 卸载挂载点
func unMountPoint(mntPath string) error {
	if _, err := exec.Command("umount", mntPath).CombinedOutput(); err != nil {
		logrus.Errorf("unmount mnt, err: %v", err)
		return err
	}
	err := os.RemoveAll(mntPath)
	if err != nil {
		logrus.Errorf("remove mnt path, err: %v", err)
		return err
	}

	return nil
}

// 删除读写层
func deleteWriteLayer(rootPath string) error {
	wirteLayerPath := path.Join(rootPath, common.WriteLayer)
	return os.RemoveAll(wirteLayerPath)
}

// 删除文件系统映射
func deleteVolume(mntPath, volume string) {
	if volume != "" {
		volumes := strings.Split(volume, ":")
		if len(volumes) > 1 {
			containerPath := path.Join(mntPath, volumes[1])
			if _, err := exec.Command("umount", containerPath).CombinedOutput(); err != nil {
				logrus.Errorf("umount container path, err: %v", err)
			}
		}
	}
}
