/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-01-18 02:33:55
 * LastEditTime: 2022-04-11 16:27:35
 * Description: socket server
 ******************************************************************************/
package agentmanager

import (
	"net"
	"strings"
	"sync"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

// 用于管理server连接的agent
type AgentManager struct {
	agentMap sync.Map
}

var globalAgentManager = AgentManager{}

// func GetAgentManager() *AgentManager {
// 	return &globalAgentManager
// }

// func regitsterHandler(s *network.SocketServer) {
// 	s.BindHandler(protocol.Heartbeat, func(conn net.Conn, data []byte) error {
// 		fmt.Println("process heartbeat:", string(data))
// 		return nil
// 	})

// 	s.BindHandler(protocol.RunScript, func(conn net.Conn, data []byte) error {
// 		fmt.Println("process run script command:", string(data))
// 		return nil
// 	})
// }

func (am *AgentManager) Stop() {
	// stop server here
}

func AddAgent(a *Agent) {
	globalAgentManager.agentMap.Store(a.UUID, a)
}

func GetAgent(uuid string) *Agent {
	agent, ok := globalAgentManager.agentMap.Load(uuid)
	if ok {
		return agent.(*Agent)
	}
	return nil
}

func GetAgentList() []map[string]string {

	agentList := []map[string]string{}

	globalAgentManager.agentMap.Range(
		func(uuid interface{}, agent interface{}) bool {
			agentInfo := map[string]string{}
			agentInfo["agent_version"] = agent.(*Agent).Version
			agentInfo["IP"] = agent.(*Agent).IP
			agentInfo["agent_uuid"] = agent.(*Agent).UUID

			agentList = append(agentList, agentInfo)
			return true
		},
	)

	return agentList
}

func DeleteAgent(uuid string) {
	if _, ok := globalAgentManager.agentMap.LoadAndDelete(uuid); !ok {
		logger.Warn("delete known agent:%s", uuid)
	}
}

func AddandRunAgent(c net.Conn) {
	agent, err := NewAgent(c)
	if err != nil {
		logger.Warn("create agent from conn error, error:%s , remote addr is:%s",
			err.Error(), c.RemoteAddr().String())
	}

	AddAgent(agent)
	logger.Info("Add new agent from:%s", c.RemoteAddr().String())
	AddAgents2DB()
}

func StopAgentManager() {

}

func AddAgents2DB() {
	var agent_list model.MachineNode
	agents := GetAgentList()
	for _, agent := range agents {
		uuid := agent["agent_uuid"]
		agent_uuid := GetAgent(uuid)
		if agent_uuid == nil {
			logger.Error("获取uuid失败!")
			return
		}
		agent_OS, err := agent_uuid.GetAgentOSInfo()
		if err != nil {
			logger.Error("初始化系统信息失败!")
			return
		}
		agentOS := strings.Split(agent_OS.(string), ";")

		agent_list.MachineUUID = uuid
		if dao.IsUUIDExist(uuid) {
			logger.Warn("机器%s已经存在!", agentOS[0])
			continue
		}
		agent_list.IP = agentOS[0]
		if dao.IsIPExist(agentOS[0]) {
			// mysqlmanager.DB.Where("ip=?", agentOS[0]).Unscoped().Delete(agent_list)
			mysqlmanager.DB.Where("ip=?", agentOS[0]).Find(&agent_list)
			if agent_list.DepartId != controller.UncateloguedDepartId {
				Ma := model.MachineNode{
					MachineUUID: uuid,
					State:       model.Normal,
				}
				mysqlmanager.DB.Model(&agent_list).Where("ip=?", agentOS[0]).Update(&Ma)
			} else {
				Ma := model.MachineNode{
					MachineUUID: uuid,
					State:       model.Free,
				}
				mysqlmanager.DB.Model(&agent_list).Where("ip=?", agentOS[0]).Update(&Ma)
			}
			continue
		}
		agent_list.DepartId = controller.UncateloguedDepartId
		agent_list.Systeminfo = agentOS[1] + " " + agentOS[2]
		agent_list.CPU = agentOS[3]
		agent_list.State = model.Free
		mysqlmanager.DB.Save(&agent_list)
	}
}
