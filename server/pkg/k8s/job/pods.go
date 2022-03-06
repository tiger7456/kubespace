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

package job

import (
	"context"
	"go.uber.org/zap"
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"kubespace/server/common"
	"kubespace/server/models/k8s"
	k8scommon "kubespace/server/pkg/k8s/common"
	"kubespace/server/pkg/k8s/event"
	"kubespace/server/pkg/k8s/pods"
)

type PodList struct {
	ListMeta k8s.ListMeta `json:"listMeta"`

	// Basic information about resources status on the list.
	Status k8scommon.ResourceStatus `json:"status"`

	// Unordered list of Pods.
	Pods []pods.Pod `json:"pods"`
}

// Returns simple info about pods(running, desired, failing, etc.) related to given job.
func getJobPodInfo(client *kubernetes.Clientset, job *batch.Job) (*k8scommon.PodInfo, error) {
	labelSelector := labels.SelectorFromSet(job.Spec.Selector.MatchLabels)
	channels := &k8scommon.ResourceChannels{
		PodList: k8scommon.GetPodListChannelWithOptions(client, k8scommon.NewSameNamespaceQuery(
			job.Namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	podList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	podInfo := k8scommon.GetPodInfo(job.Status.Active, job.Spec.Completions, podList.Items)

	// This pod info for jobs should be get from job status, similar to kubectl describe logic.
	podInfo.Running = job.Status.Active
	podInfo.Succeeded = job.Status.Succeeded
	podInfo.Failed = job.Status.Failed
	return &podInfo, nil
}

func getJobToPod(client *kubernetes.Clientset, job *batch.Job) (po *PodList) {

	selector, err := metaV1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil
	}
	options := metaV1.ListOptions{LabelSelector: selector.String()}

	podData, err := client.CoreV1().Pods(job.Namespace).List(context.TODO(), options)
	if err != nil {
		common.LOG.Error("Get a pod exception from the job", zap.Any("err", err))
	}
	podList := PodList{
		Pods: make([]pods.Pod, 0),
	}
	podList.ListMeta = k8s.ListMeta{TotalItems: len(podData.Items)}
	for _, pod := range podData.Items {
		warnings := event.GetPodsEventWarnings(nil, []v1.Pod{pod})
		podDetail := pods.ToPod(&pod, warnings)
		podList.Pods = append(podList.Pods, podDetail)
	}
	return &podList
}
