package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func DiskUsageHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_use, err := agent.DiskUsage()
	if err != nil {
		response.Fail(c, nil, "获取磁盘的使用情况失败!")
		return
	}
	response.Success(c, gin.H{"disk_use": disk_use}, "Success")
}

func DiskInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_info, err := agent.DiskInfo()
	if err != nil {
		response.Fail(c, nil, "获取磁盘的IO信息失败!")
		return
	}
	response.Success(c, gin.H{"disk_info": disk_info}, "Success")
}
func DiskCreatPathHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	mountpath := c.Query("mountpath")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_path, err := agent.DiskCreatPath(mountpath)
	if disk_path != nil || err != nil {
		response.Fail(c, gin.H{"error": disk_path}, "创建挂载目录失败!")
		return
	}
	response.Success(c, gin.H{"disk_path": disk_path}, "Success")
}
func DiskMountHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	sourceDisk := c.Query("source")
	destPath := c.Query("dest")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_mount, err := agent.DiskMount(sourceDisk, destPath)
	if disk_mount != nil || err != nil {
		response.Fail(c, gin.H{"error": disk_mount}, "挂载磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_mount": disk_mount}, "Success")
}
func DiskUMountHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	diskPath := c.Query("path")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_umount, err := agent.DiskUMount(diskPath)
	if disk_umount != nil || err != nil {
		response.Fail(c, gin.H{"error": disk_umount}, "卸载磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_umount": disk_umount}, "Success")
}
func DiskFormatHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	fileType := c.Query("type")
	diskPath := c.Query("path")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	disk_format, err := agent.DiskFormat(fileType, diskPath)
	if disk_format != nil || err != nil {
		response.Fail(c, gin.H{"error": disk_format}, "格式化磁盘失败!")
		return
	}
	response.Success(c, gin.H{"disk_format": disk_format}, "Success")
}
