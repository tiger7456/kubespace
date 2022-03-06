/*




Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubespace/server/common"
	"kubespace/server/controller/response"
	"kubespace/server/inner/cloud/cloudsync"
	"kubespace/server/inner/cloud/cloudvendor"
	"kubespace/server/models/cmdb"
	"kubespace/server/models/request"
	"kubespace/server/services"
)

var (
	pageInfo request.PageInfo
	account  cmdb.CloudPlatform
)

func ListPlatform(c *gin.Context) {
	_ = c.ShouldBindJSON(&pageInfo)
	err, list, total := services.ListPlatform(pageInfo)
	if err != nil {
		common.LOG.Error("获取云平台信息失败", zap.Any("err", err))
		response.FailWithMessage(500, fmt.Sprintf("获取云平台信息失败，%v", err), c)
		return

	} else {
		response.OkWithData(response.PageResult{
			Data:  list,
			Total: total,
			Page:  pageInfo.Page,
			Size:  pageInfo.PageSize,
		}, c)
		return
	}
}

// CloudPlatformAccount 云平台账号
func CloudPlatformAccount(c *gin.Context) {
	_ = c.ShouldBindJSON(&account)

	// 校验云厂商客户端
	_, err := cloudvendor.GetVendorClient(&account)
	if err != nil {
		response.FailWithMessage(500, fmt.Sprintf("AccountVerify GetVendorClient failed，%v", err), c)
		return
	}
	// TODO 校验AccessKey

	// 创建云账号
	err1 := services.CreateCloudAccount(&account)
	if err1 != nil {
		response.FailWithMessage(500, fmt.Sprintf("创建云平台账号异常，%v", err1), c)
		return
	}

	// 开启协程，后台同步ecs
	taskChan := make(chan *cmdb.CloudPlatform)
	SyncCloudResource(taskChan)

	response.OkWithDetailed("null", "添加成功, 任务正在后台同步云资源", c)
	return
}

// SyncCloudResource 同步云资源
func SyncCloudResource(taskChan chan *cmdb.CloudPlatform) {
	go func() {
		if task, ok := <-taskChan; ok {
			switch task.Type {
			case "aliyun":
				cloudsync.SyncAliYunHost(task)
			case "tencent":
				return
			case "aws":
				return
			default:
				common.LOG.Error(fmt.Sprintf("unknown resource type:%v, ignore it!", task.Type))
			}
		}
	}()
	taskChan <- &account
}
