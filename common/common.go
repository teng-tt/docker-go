package common

// 定义目录
const (
	RootPath   = "/root/"
	MntPath    = "/root/mnt/"
	WriteLayer = "writeLayer"
)

const (
	Running = "running"
	Stop    = "stopped"
	Exit    = "exited"
)

const (
	DefaultContainerInfoPath = "/var/run/docker-go/"
	ContainerInfoFileName    = "config.json"
	ContainerLogFileName     = "container.log"
)

const (
	EnvExecPid = "docker_pid"
	EnvExecCmd = "docker_cmd"
)

const (
	DefaultNetworkPath   = "/var/run/docker-go/network/network/"
	DefaultAllocatorPath = "/var/run/docker-go/network/ipam/subnet.json"
)
