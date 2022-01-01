package netproxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/schwarzeni/k8senhance/pkg/pb"
	"github.com/schwarzeni/k8senhance/pkg/util"
)

type ProxyRequestCallBack func(w *pb.Response, err error)

type ProxyRequestWrapper struct {
	Data     *pb.Request
	CallBack ProxyRequestCallBack
}

type HTTPProxyServer struct {
	proxyChannelManager *ProxyChannelManager
}

func (hps *HTTPProxyServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	svcName, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			log.Println("[err] split host and port from ", req.Host, err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		svcName = req.Host
	}
	log.Println("[debug] proxy for svc", svcName)
	// TODO: 优化：流式读取
	bodyContent, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("[err] read req body", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	waitChan := make(chan struct{})

	targetNode, targetAddr, err := SelectTarget(util.MustParseServiceNameFromHost(svcName))
	if err != nil {
		log.Println("[err] failed to select target", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	ch, ok := hps.proxyChannelManager.Get(targetNode)
	if !ok {
		log.Println("[err] target node channel not found", targetNode)
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	ch <- &ProxyRequestWrapper{
		Data: &pb.Request{
			Id:         util.GenID(svcName),
			TargetAddr: targetAddr,
			TargetNode: targetNode,
			Header:     util.ExtractHeader(req.Header),
			Body:       bodyContent,
			HttpMethod: req.Method,
			Url:        req.URL.String(),
		},
		CallBack: func(proxyResponse *pb.Response, err error) {
			defer func() {
				waitChan <- struct{}{}
			}()
			if err != nil {
				log.Println("[err] start callback", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			if proxyResponse.Header != nil {
				for _, item := range proxyResponse.Header.Item {
					if len(item.Value) > 1 {
						for _, v := range item.Value[1:] {
							resp.Header().Add(item.Value[0], v)
						}
					}
				}
			}
			// TODO: 优化：流式读写
			if _, err := io.Copy(resp, bytes.NewReader(proxyResponse.Body)); err != nil {
				log.Println("[err] copy body", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
		},
	}
	<-waitChan
	log.Println("[debug] session finish")
}

func RunHTTPProxyServer(addr string, proxyChannelManager *ProxyChannelManager) {
	server := http.Server{Addr: addr, Handler: &HTTPProxyServer{proxyChannelManager: proxyChannelManager}}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
