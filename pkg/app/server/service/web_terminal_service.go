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
 * Date: 2022-04-20 16:48:55
 * LastEditTime: 2022-04-20 17:48:55
 * Description: web socket连接逻辑代码
 ******************************************************************************/

package service

import (
	"encoding/base64"
	"encoding/json"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
)

func DecodedMsgToSSHClient(msg string) (model.SSHClient, error) {
	client := model.SSHClient{}
	decoded, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return client, err
	}
	err = json.Unmarshal(decoded, &client)
	if err != nil {
		return client, err
	}
	return client, nil
}
