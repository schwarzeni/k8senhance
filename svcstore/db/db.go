package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var inmemoryDB = make(map[string]*model.ServiceInfo)
var lock sync.RWMutex

func GetService(serviceName string) (result *model.ServiceInfo, err error) {
	// get basic info
	lock.RLock()
	result, ok := inmemoryDB[serviceName]
	if !ok {
		lock.RUnlock()
		return nil, fmt.Errorf("no service named: %s", serviceName)
	}
	lock.RUnlock()

	// query k8s for backendInfo
	podList, err := podclient.List(context.TODO(), metav1.ListOptions{LabelSelector: labels.SelectorFromSet(result.Selector).String()})
	if err != nil {
		return nil, fmt.Errorf("access k8s: %v", err)
	}
	for _, pod := range podList.Items {
		result.Backends = append(result.Backends, &model.Backend{
			NodeName:    pod.Spec.NodeName,
			ContainerIP: pod.Status.PodIP, // TODO: 这里仅支持单容器 Pod ，所以取第一个容器的 IP
		})
	}
	return result, nil
}

func PutService(service *model.ServiceInfo) error {
	lock.Lock()
	defer lock.Unlock()
	inmemoryDB[service.Name] = service
	return nil
}

func DelService(serviceName string) error {
	lock.Lock()
	defer lock.Unlock()
	delete(inmemoryDB, serviceName)
	return nil
}
