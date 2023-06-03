/*
	资源限制接口，Apply将进程ID 添加到tasks 中，即将此进程加入 cgroup 中
	Set 则对某个资源进行限制，Remove 则为移除该 cgroup
	都比较简单，就是创建文件，写文件罢了，理解原理之后，写起来很轻松
*/

package subsystem

// ResourceConfig 资源配置限制
type ResourceConfig struct {
	MemoryLimit string // 内存限制
	CpuShare    string // CPU时间片权重
	CpuSet      string // CPU核数
}

// Subystem 将cgroup抽象成path, 因为在hierarchy中，cgroup便是虚拟的路径地址
type Subystem interface {
	Name() string                                     // 返回subsystem名字， 如cpu, memory
	Set(cgroupPath string, res *ResourceConfig) error // 设置cgroup在这个subsystem中的资源限制
	Remove(cgroupPath string) error                   // 移除这个cgroup资源限制
	Apply(cgroupPath string, pid int) error           // 将某个进程添加到cgroup中
}

var (
	Subsystems = []Subystem{
		&MemorySubSystem{},
		// 设置task时这两个必须设置
		&CpuSubSystem{},
		&CpuSetSubSystem{},
	}
)
