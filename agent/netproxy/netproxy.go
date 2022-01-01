package netproxy

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/schwarzeni/k8senhance/pkg/pb"
	"github.com/schwarzeni/k8senhance/pkg/util"
	"google.golang.org/grpc"

	"github.com/schwarzeni/k8senhance/config"
)

type NetProxy struct {
	gc *config.Config
}

func (ic *NetProxy) Run() error {
	log.Println("start net proxy service")
	conn, err := grpc.Dial(ic.gc.Agent.Netproxy.CloudGrpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewProxyHttpServiceClient(conn)
	go func() {
		if err := NewHTTPProxyServer(ic.gc.Agent.Netproxy.Addr, c, ic.gc.NodeName).ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	RunEdgeCloudServer(c, ic.gc.NodeName)
	return nil
}

func NewNetProxy(gc *config.Config) *NetProxy {
	return &NetProxy{gc: gc}
}

func RunEdgeCloudServer(c pb.ProxyHttpServiceClient, nodeid string) {
	stream, err := c.ProxyCloud2Edge(context.TODO())
	if err != nil {
		panic(err)
	}

	err = stream.Send(&pb.Response{
		Hello:  true,
		Nodeid: nodeid,
	})
	if err != nil {
		log.Fatalf("failed to connect to cloud %v", err)
	}
	handshakeresp, err := stream.Recv()
	if err != nil || !handshakeresp.Hello {
		log.Fatalf("invalid handshake %v", err)
	}
	for {
		reqPB, err := stream.Recv()
		if err != nil {
			log.Println("[err] stream recv", err)
			return
		}
		callbackFunc := func(resp *http.Response, err error) {
			if err != nil {
				log.Printf("[err] start callback %v", err)
				return
			}
			respBodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				log.Printf("[err] read resp body %v", err)
				return
			}
			respPB := &pb.Response{
				Id:     reqPB.Id,
				Nodeid: nodeid,
				Header: util.ExtractHeader(resp.Header),
				Body:   respBodyBytes,
			}
			if err := stream.Send(respPB); err != nil {
				log.Printf("[err] failed to send, %v", err)
				return
			}
		}
		realReq, err := http.NewRequest(reqPB.HttpMethod, "http://"+reqPB.TargetAddr+reqPB.Url, bytes.NewReader(reqPB.Body))
		if err != nil {
			log.Printf("[err] failed build url: %v", err)
			continue
		}
		if reqPB.Header != nil {
			for _, vv := range reqPB.Header.Item {
				if len(vv.Value) > 1 {
					for _, v := range vv.Value[1:] {
						realReq.Header.Add(vv.Value[0], v)
					}
				}
			}
		}
		// TODO: maybe use goroutine pool later
		go func() {
			callbackFunc(http.DefaultClient.Do(realReq))
		}()
	}
}
