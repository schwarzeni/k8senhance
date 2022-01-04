package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type URLInfo struct {
	ImageName  string
	SourceType string
	SourceID   string
}

func MustParseURL(reqURL string) *URLInfo {
	urlItems := strings.Split(reqURL, "/")
	imageName := urlItems[2] + "/" + urlItems[3]
	sourceType := urlItems[4]
	sourceID := urlItems[5]
	return &URLInfo{
		ImageName:  imageName,
		SourceType: sourceType,
		SourceID:   sourceID,
	}
}

func copyHeader(src http.Header, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
			//log.Println("[debug] manifest header: ", k, v)
		}
	}
}

var (
	authCache     = map[string]string{}
	authCacheLock sync.RWMutex
)

func withDockerhubPullAuth(req *http.Request, imageName string) (err error) {
	// TODO: 后期缓存可以采用 cache + singleflight 优化
	authCacheLock.Lock()
	defer authCacheLock.Unlock()
	token, ok := authCache[imageName]
	if !ok {
		log.Println("[debug] request token for ", imageName)
		resp, err := http.Get(fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", imageName))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		respObj := dockerhubAuthObj{}
		respData, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(respData, &respObj); err != nil {
			return err
		}
		token = respObj.Token
		authCache[imageName] = respObj.Token
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return nil
}

type dockerhubAuthObj struct {
	Token string `json:"token"`
}

func doGetProxy(remoteAddr string, rawReq *http.Request) (*http.Response, error) {
	targetURL := rawReq.URL.String()
	// TODO: set this proxy to be configtable
	//os.Setenv("HTTP_PROXY", "http://10.211.55.2:7890")
	//os.Setenv("HTTPS_PROXY", "http://10.211.55.2:7890")
	req, err := http.NewRequest(http.MethodGet, remoteAddr+targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("newRequest for %s: %v", targetURL, err)
	}
	for k, vv := range rawReq.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	client := http.DefaultClient
	client.Timeout = time.Second * 100000
	return client.Do(req)
}
