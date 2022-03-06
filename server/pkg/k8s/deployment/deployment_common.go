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

package deployment

import (
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8scommon "kubespace/server/pkg/k8s/common"
	"kubespace/server/pkg/k8s/dataselect"
	"kubespace/server/pkg/k8s/event"
)

// The code below allows to perform complex data section on Deployment

type DeploymentCell apps.Deployment

// GetProperty is used to get property of the deployment
func (self DeploymentCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		// if name is not supported then just return a constant dummy value, sort will have no effect.
		return nil
	}
}

func toCells(std []apps.Deployment) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = DeploymentCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []apps.Deployment {
	std := make([]apps.Deployment, len(cells))
	for i := range std {
		std[i] = apps.Deployment(cells[i].(DeploymentCell))
	}
	return std
}

func getStatus(list *apps.DeploymentList, rs []apps.ReplicaSet, pods []v1.Pod, events []v1.Event) k8scommon.ResourceStatus {
	info := k8scommon.ResourceStatus{}
	if list == nil {
		return info
	}

	for _, deployment := range list.Items {
		matchingPods := k8scommon.FilterDeploymentPodsByOwnerReference(deployment, rs, pods)
		podInfo := k8scommon.GetPodInfo(deployment.Status.Replicas, deployment.Spec.Replicas, matchingPods)
		warnings := event.GetPodsEventWarnings(events, matchingPods)

		if len(warnings) > 0 {
			info.Failed++
		} else if podInfo.Pending > 0 {
			info.Pending++
		} else {
			info.Running++
		}
	}

	return info
}

func getConditions(deploymentConditions []apps.DeploymentCondition) []k8scommon.Condition {
	conditions := make([]k8scommon.Condition, 0)

	for _, condition := range deploymentConditions {
		conditions = append(conditions, k8scommon.Condition{
			Type:               string(condition.Type),
			Status:             metaV1.ConditionStatus(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime,
			LastProbeTime:      condition.LastUpdateTime,
		})
	}

	return conditions
}
