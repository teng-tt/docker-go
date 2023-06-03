// 内存限制实例

package subsystem

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
	apply bool
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}

// Set 设置内存限制
func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", cgroupPath, err)
		return err
	}
	if res.MemoryLimit != "" {
		// 设置cgroup内存限制
		// 将这个限制写入到cgroup对应目录的 memory.limit_in_bytes 文件即可
		err := ioutil.WriteFile(path.Join(subsystemCgroupPath, "memoey.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// Remove 移除内存限制
func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}

	return os.RemoveAll(subsystemCgroupPath)
}

// Apply 将进程加入内存限制
func (m *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}

	// 将进程id 写入task文件进行继承父内存限制
	taskPath := path.Join(subsystemCgroupPath, "tasks")
	err = ioutil.WriteFile(taskPath, []byte(strconv.Itoa(pid)), os.ModePerm)
	if err != nil {
		logrus.Errorf("write pid to task, path %s, pid: %d, err: %v", taskPath, pid, err)
		return err
	}

	return nil
}
