package kube

import corev1 "k8s.io/api/core/v1"

func GetPodReadyContainers(containerStatuses []corev1.ContainerStatus) (ready, total int) {
	total = len(containerStatuses)
	for _, status := range containerStatuses {
		if status.Ready {
			ready++
		}
	}
	return ready, total
}
