package netproxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/pb"
	"github.com/schwarzeni/k8senhance/pkg/util"
)

type CloudEdgeServer struct {
	pb.UnimplementedProxyHttpServiceServer
	RequestCallBackMap  map[string]ProxyRequestCallBack
	proxyChannelManager *ProxyChannelManager
}

func (s *CloudEdgeServer) ProxyEdge2Cloud(stream pb.ProxyHttpService_ProxyEdge2CloudServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("[err] ProxyCloud2Edge stream.Recv: %v\n", err)
			return err
		}

		targetNode, targetAddr, err := SelectTarget(util.MustParseServiceNameFromID(req.Id))
		if err != nil {
			log.Println("[err] failed to select target", err)
			_ = stream.Send(&pb.Response{
				Id:         req.Id,
				Statuscode: http.StatusServiceUnavailable,
			})
			continue
		}
		//targetNode := "oooo"
		//targetAddr := "10.211.55.52:8080"
		//targetNode := "xxxx"
		//targetAddr := "10.211.55.2:8080"
		ch, ok := s.proxyChannelManager.Get(targetNode)
		if !ok {
			log.Println("[err] target node channel not found", targetNode)
			_ = stream.Send(&pb.Response{
				Id:         req.Id,
				Statuscode: http.StatusServiceUnavailable,
			})
			continue
		}
		ch <- &ProxyRequestWrapper{
			Data: &pb.Request{
				Id:         req.Id,
				TargetAddr: targetAddr,
				TargetNode: targetNode,
				Header:     req.Header,
				Body:       req.Body,
				HttpMethod: req.HttpMethod,
				Url:        req.Url,
			},
			CallBack: func(proxyResponse *pb.Response, err error) {
				if err != nil {
					log.Println("[err] ProxyEdge2Cloud Callback start callback", err)
					_ = stream.Send(&pb.Response{
						Id:         req.Id,
						Statuscode: http.StatusServiceUnavailable,
					})
					return
				}
				// TODO: 优化：流式读写
				if err := stream.Send(proxyResponse); err != nil {
					log.Println("[err] ProxyEdge2Cloud Callback failed to send stream", err)
				}
			},
		}
	}
}

func (s *CloudEdgeServer) ProxyCloud2Edge(stream pb.ProxyHttpService_ProxyCloud2EdgeServer) error {
	stopChan := make(chan struct{})
	resp, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Println("[err] ProxyCloud2Edge stream.Recv ", err)
		return err
	}
	edgeNodeId := resp.Nodeid
	dataChan := make(chan *ProxyRequestWrapper)
	if resp.Hello {
		s.proxyChannelManager.Set(edgeNodeId, dataChan)
		err := stream.Send(&pb.Request{
			Hello: true,
		})
		if err != nil {
			log.Fatalf("failed to finish handshake: %v", err)
		}
		log.Println("[debug] register to stream manager", edgeNodeId)
	} else {
		log.Println("[err] invalid interaction")
		return fmt.Errorf("invalid interaction")
	}

	var rwLock sync.RWMutex

	// get response loop
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(stopChan)
				break
			}
			if err != nil {
				log.Println("[err] stream.Recv ", err)
				break
			}
			rwLock.RLock()
			fn, ok := s.RequestCallBackMap[resp.Id]
			rwLock.RUnlock()
			if !ok {
				log.Println("[warn] no callback for ", resp.Id)
				continue
			}
			fn(resp, nil)
			rwLock.Lock()
			delete(s.RequestCallBackMap, resp.Id)
			rwLock.Unlock()
		}
	}()

	// send request loop
	go func() {
		for {
			select {
			case <-stopChan:
				return
			case prw := <-dataChan:
				rwLock.Lock()
				s.RequestCallBackMap[prw.Data.Id] = prw.CallBack
				rwLock.Unlock()
				if err := stream.Send(prw.Data); err != nil {
					log.Printf("[err] send data to %s failed: %v", prw.Data.Id, err)
				}
			}
		}
	}()

	<-stopChan
	s.proxyChannelManager.Del(edgeNodeId)
	return nil
}
