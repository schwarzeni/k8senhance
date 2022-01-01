package db

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	podclient typedcorev1.PodInterface
	k8sclient *kubernetes.Clientset
)

func InitK8SClient(configPath string) {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err)
	}
	if k8sclient, err = kubernetes.NewForConfig(config); err != nil {
		panic(err)
	}
	// TODO: 这里默认全部 Pod 都位于 default 命名空间中
	podclient = k8sclient.CoreV1().Pods(corev1.NamespaceDefault)
}
