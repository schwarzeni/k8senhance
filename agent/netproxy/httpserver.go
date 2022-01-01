package netproxy

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/pb"
	"github.com/schwarzeni/k8senhance/pkg/util"
)

type ProxyRequestCallBack func(w *pb.Response, err error)

type ProxyRequestWrapper struct {
	Data     *pb.Request
	CallBack ProxyRequestCallBack
}

type HTTPProxyServer struct {
	dataChan           chan *ProxyRequestWrapper
	RequestCallBackMap map[string]ProxyRequestCallBack
	NodeID             string
}

func (hps *HTTPProxyServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	svcName, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		log.Println("[err] split host and port from ", req.Host, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
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

	hps.dataChan <- &ProxyRequestWrapper{
		Data: &pb.Request{
			Id:         util.GenID(svcName),
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

func NewHTTPProxyServer(addr string, c pb.ProxyHttpServiceClient, nodeid string) *http.Server {
	ch := make(chan *ProxyRequestWrapper)
	reqCBMap := make(map[string]ProxyRequestCallBack)
	rwLock := sync.RWMutex{}
	server := &HTTPProxyServer{
		dataChan:           ch,
		RequestCallBackMap: reqCBMap,
		NodeID:             nodeid,
	}
	stream, err := c.ProxyEdge2Cloud(context.TODO())
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("[err] stream.Recv ", err)
				break
			}
			rwLock.RLock()
			fn, ok := server.RequestCallBackMap[resp.Id]
			rwLock.RUnlock()
			if !ok {
				log.Println("[warn] no callback for ", resp.Id)
				continue
			}
			fn(resp, nil)
			rwLock.Lock()
			delete(server.RequestCallBackMap, resp.Id)
			rwLock.Unlock()
		}
	}()
	go func() {
		for {
			select {
			case prw := <-server.dataChan:
				rwLock.Lock()
				server.RequestCallBackMap[prw.Data.Id] = prw.CallBack
				rwLock.Unlock()
				if err := stream.Send(prw.Data); err != nil {
					log.Printf("[err] send data to %s failed: %v", prw.Data.Id, err)
				}
			}
		}
	}()

	return &http.Server{
		Addr:    addr,
		Handler: server,
	}
}
