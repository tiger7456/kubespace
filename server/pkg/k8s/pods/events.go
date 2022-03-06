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

package pods

import (
	"k8s.io/client-go/kubernetes"
	"kubespace/server/pkg/k8s/common"
	"kubespace/server/pkg/k8s/dataselect"
	"kubespace/server/pkg/k8s/event"
)

// GetEventsForPod gets events that are associated with this pod.
func GetEventsForPod(client *kubernetes.Clientset, dsQuery *dataselect.DataSelectQuery, namespace, podName string) (*common.EventList, error) {
	return event.GetResourceEvents(client, dsQuery, namespace, podName)
}
