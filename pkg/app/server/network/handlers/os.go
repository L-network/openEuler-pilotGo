package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func OSInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	os_info, err := agent.GetOSInfo()
	if err != nil {
		response.Fail(c, nil, "获取系统信息失败!")
		return
	}
	response.Success(c, gin.H{"os_info": os_info}, "Success")
}