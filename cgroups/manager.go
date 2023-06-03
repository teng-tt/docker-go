/*
	资源限制管理器
*/

package cgroups

import (
	"docker-go/cgroups/subsystem"
	"github.com/sirupsen/logrus"
)

type CGroupManager struct {
	Path string
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{Path: path}
}

func (c *CGroupManager) Set(res *subsystem.ResourceConfig) {
	for _, subystem := range subsystem.Subsystems {
		err := subystem.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("set %s err: %v", subystem.Name(), err)
		}
	}
}

func (c *CGroupManager) Apply(pid int) {
	for _, subystem := range subsystem.Subsystems {
		err := subystem.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply task, err: %v", err)
		}
	}
}

func (c *CGroupManager) Destroy() {
	for _, subystem := range subsystem.Subsystems {
		err := subystem.Remove(c.Path)
		if err != nil {
			logrus.Errorf("remove %s err: %v", subystem.Name(), err)
		}

	}
}
