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

package k8s

type JobData struct {
	Namespace string `json:"namespace"  binding:"required"`
	Name      string `json:"name"  binding:"required"`
}

type ScaleJob struct {
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
	// Job并行运行的Pod数量
	Number *int32 `json:"number" binding:"required"`
}
