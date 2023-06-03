package container

import (
	"docker-go/common"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// ContainerInfo 容器信息
type ContainerInfo struct {
	Pid         string   `json:"pid"`     // 容器的init进程在宿主机上的PID
	Id          string   `json:"id"`      // 容器ID
	Command     string   `json:"command"` // 容器内init进程运行的命令
	Name        string   `json:"name"`
	CreateTime  string   `json:"createTime"`
	Status      string   `json:"status"`
	Volume      string   `json:"volume"`      // 容器的数据卷
	PortMapping []string `json:"portmapping"` // 端口映射
}

// RecordContainerInfo 记录容器信息
// 1. 创建以容器名或 ID 命名的文件夹
// 2. 在该文件下创建 config.json
// 3. 将容器信息保存到 config.json 中
func RecordContainerInfo(containerPID int, cmdArray []string, containerName, containerID string) error {
	// 生成容器基础信息
	info := &ContainerInfo{
		Pid:        strconv.Itoa(containerPID),
		Id:         containerID,
		Command:    strings.Join(cmdArray, " "),
		Name:       containerName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     common.Running,
	}
	// 创建容器目录
	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir container dir: %s, err: %v", dir, err)
			return err
		}
	}
	// 创建容器信息文件
	fileName := fmt.Sprintf("%s/%s", dir, common.ContainerInfoFileName)
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("create config.json, fileName: %s, err: %v", fileName, err)
		return err
	}
	// 往文件写入内容
	bs, _ := json.Marshal(info)
	_, err = file.WriteString(string((bs)))
	if err != nil {
		logrus.Errorf("write config.json, fileName: %s, err: %v", fileName, err)
		return err
	}

	return nil
}

// GetContainerID 获取容器ID
func GetContainerID(n int) string {
	letterBytes := "0123456789"
	// 使用随机种子，选取随机值生成容器id
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// DeleteContainerInfo 删除容器信息
func DeleteContainerInfo(containerName string) {
	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	err := os.RemoveAll(dir)
	if err != nil {
		logrus.Errorf("remove container info, err: %v", err)
	}
}

// ListContainerInfo 列出当前的容器如 docker ps命令
// 1. 遍历 docker-go 文件夹
// 2. 读取每个容器内的 config.json 文件
// 3. 格式化打印
func ListContainerInfo() {
	files, err := ioutil.ReadDir(common.DefaultContainerInfoPath)
	if err != nil {
		logrus.Errorf("read info dir, err: %v", err)
	}
	var infos []*ContainerInfo
	// 1. 遍历 docker-go 文件夹
	for _, file := range files {
		// 2. 读取每个容器内的 config.json 文件
		info, err := getContainerInfo(file.Name())
		if err != nil {
			logrus.Errorf("get container info, name: %s, err: %v", file.Name(), err)
			continue
		}
		infos = append(infos, info)
	}

	// 3. 格式化打印
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 2, ' ', 0)
	_, _ = fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Id, info.Name, info.Pid, info.Status, info.Command, info.CreateTime)
	}
	// 刷新标准输出流缓存区，将容器列表打印出来
	if err := w.Flush(); err != nil {
		logrus.Errorf("flush info, err:%v", err)
	}
}

// 获取容器详细信息
func getContainerInfo(containerName string) (*ContainerInfo, error) {
	filePath := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerInfoFileName)
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("read file, path: %s, err: %v", filePath, err)
		return nil, err
	}
	info := &ContainerInfo{}
	err = json.Unmarshal(bs, info)

	return info, err
}
