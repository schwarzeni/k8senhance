package model

type ImageCacheAgentHealthReqDTO struct {
	NodeID               string `json:"node_id"`
	IP                   string `json:"ip"`
	Region               string `json:"region"`
	ImageCacheServerPort string `json:"image_cache_server_port"`
}

type ImageCacheAgentHealthRespDTO struct {
	Region string                         `json:"region"`
	Nodes  []*ImageCacheAgentHealthReqDTO `json:"nodes"`
}
