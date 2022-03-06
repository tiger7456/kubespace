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

package service

import (
	"fmt"
	client "k8s.io/client-go/kubernetes"
	"kubespace/server/common"
	"kubespace/server/models/k8s"
	k8scommon "kubespace/server/pkg/k8s/common"
	"kubespace/server/pkg/k8s/dataselect"
	"kubespace/server/pkg/k8s/event"
)

// GetServiceEvents returns model events for a service with the given name in the given namespace.
func GetServiceEvents(client *client.Clientset, dsQuery *dataselect.DataSelectQuery, namespace, name string) (*k8scommon.EventList, error) {
	eventList := k8scommon.EventList{
		Events:   make([]k8scommon.Event, 0),
		ListMeta: k8s.ListMeta{TotalItems: 0},
	}

	serviceEvents, err := event.GetEvents(client, namespace, name)
	if err != nil {
		return &eventList, err
	}

	eventList = event.CreateEventList(event.FillEventsType(serviceEvents), dsQuery)
	common.LOG.Info(fmt.Sprintf("Found %d events related to %s service in %s namespace", len(eventList.Events), name, namespace))
	return &eventList, nil
}
