package util

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/schwarzeni/k8senhance/pkg/pb"
)

func ExtractHeader(header http.Header) *pb.HTTPHeader {
	pbHeader := &pb.HTTPHeader{}
	for k, vs := range header {
		vvv := &pb.HTTPHeaderValue{}
		vvv.Value = append(vvv.Value, k)
		for _, v := range vs {
			vvv.Value = append(vvv.Value, v)
		}
		pbHeader.Item = append(pbHeader.Item, vvv)
	}
	return pbHeader
}

func init() {
	rand.Seed(time.Now().Unix())
}

func GenID(svcName string) string {
	// TODO: 完善随机数生成机制
	return fmt.Sprintf("%s-%d-%d", svcName, time.Now().Unix(), rand.Intn(10000))
}

func MustParseServiceNameFromHost(svcName string) string {
	tmp := strings.Split(svcName, ".")
	return strings.Join(tmp[:len(tmp)-2], ".")
}

func MustParseServiceNameFromID(id string) string {
	tmp := strings.Split(id, "-")
	return strings.Join(tmp[:len(tmp)-2], "-")
}
