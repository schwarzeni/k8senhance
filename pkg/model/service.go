package model

type ServiceInfo struct {
	Name     string            `json:"name"` // 不包含后缀
	Port     string            `json:"port"`
	Selector map[string]string `json:"selector"`
	Backends []*Backend        `json:"backends"`
}

type Backend struct {
	NodeName    string `json:"node_name"`
	ContainerIP string `json:"container_ip"` // TODO: 目前仅支持单容器 Pod
}
