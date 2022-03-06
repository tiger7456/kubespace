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

package routers

import (
	"github.com/gin-gonic/gin"
	"kubespace/server/controller"
)

func InitCloudRouter(r *gin.RouterGroup) {
	InitCloudRouter := r.Group("cloud")
	{
		InitCloudRouter.GET("listPlatform", controller.ListPlatform)
		InitCloudRouter.POST("account", controller.CloudPlatformAccount)
	}
}
