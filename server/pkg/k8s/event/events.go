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

package event

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"kubespace/server/pkg/k8s/common"
	"strings"
)

func GetClusterNodeEvent(client *kubernetes.Clientset, namespace string, field string) (*v1.EventList, error) {

	events, err := client.CoreV1().Events(namespace).List(context.TODO(),
		metav1.ListOptions{
			FieldSelector: field,
		},
	)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetPodsEventWarnings returns warning pod events by filtering out events targeting only given pods
func GetPodsEventWarnings(events []v1.Event, pods []v1.Pod) []common.Event {
	result := make([]common.Event, 0)

	// Filter out only warning events
	events = getWarningEvents(events)
	failedPods := make([]v1.Pod, 0)

	// Filter out ready and successful pods
	for _, pod := range pods {
		if !isReadyOrSucceeded(pod) {
			failedPods = append(failedPods, pod)
		}
	}

	// Filter events by failed pods UID
	events = filterEventsByPodsUID(events, failedPods)
	events = removeDuplicates(events)

	for _, event := range events {
		result = append(result, common.Event{
			Message: event.Message,
			Reason:  event.Reason,
			Type:    event.Type,
		})
	}

	return result
}

// Returns filtered list of event objects.
// Event list object is filtered to get only warning events.
func getWarningEvents(events []v1.Event) []v1.Event {
	return filterEventsByType(FillEventsType(events), v1.EventTypeWarning)
}

// Returns true if given pod is in state ready or succeeded, false otherwise
func isReadyOrSucceeded(pod v1.Pod) bool {
	if pod.Status.Phase == v1.PodSucceeded {
		return true
	}
	if pod.Status.Phase == v1.PodRunning {
		for _, c := range pod.Status.Conditions {
			if c.Type == v1.PodReady {
				if c.Status == v1.ConditionFalse {
					return false
				}
			}
		}

		return true
	}

	return false
}

// Returns filtered list of event objects. Events list is filtered to get only events targeting
// pods on the list.
func filterEventsByPodsUID(events []v1.Event, pods []v1.Pod) []v1.Event {
	result := make([]v1.Event, 0)
	podEventMap := make(map[types.UID]bool, 0)

	if len(pods) == 0 || len(events) == 0 {
		return result
	}

	for _, pod := range pods {
		podEventMap[pod.UID] = true
	}

	for _, event := range events {
		if _, exists := podEventMap[event.InvolvedObject.UID]; exists {
			result = append(result, event)
		}
	}

	return result
}

// Removes duplicate strings from the slice
func removeDuplicates(slice []v1.Event) []v1.Event {
	visited := make(map[string]bool, 0)
	result := make([]v1.Event, 0)

	for _, elem := range slice {
		if !visited[elem.Reason] {
			visited[elem.Reason] = true
			result = append(result, elem)
		}
	}

	return result
}

// Filters kubernetes API event objects based on event type.
// Empty string will return all events.
func filterEventsByType(events []v1.Event, eventType string) []v1.Event {
	if len(eventType) == 0 || len(events) == 0 {
		return events
	}

	result := make([]v1.Event, 0)
	for _, event := range events {
		if event.Type == eventType {
			result = append(result, event)
		}
	}

	return result
}

// Returns true if reason string contains any partial string indicating that this may be a
// warning, false otherwise
func isFailedReason(reason string, partials ...string) bool {
	for _, partial := range partials {
		if strings.Contains(strings.ToLower(reason), partial) {
			return true
		}
	}

	return false
}
