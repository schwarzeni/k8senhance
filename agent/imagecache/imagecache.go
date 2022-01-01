package imagecache

import (
	"log"

	"github.com/schwarzeni/k8senhance/config"
)

type ImageCache struct {
	gc *config.Config
}

func (ic *ImageCache) Run() error {
	log.Println("start image cache service")
	return nil
}

func NewImageCache(gc *config.Config) *ImageCache {
	return &ImageCache{gc: gc}
}
