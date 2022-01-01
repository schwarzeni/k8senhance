package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/schwarzeni/k8senhance/pkg/model"
)

var client = &http.Client{}

func PutServiceInfo(addr string, svc *model.ServiceInfo) error {
	bytedata, err := json.Marshal(svc)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, addr+"/api/v1/service", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("get request %v", err)
	}
	_ = resp.Body.Close()
	return nil
}

func GetServiceInfo(addr string, serviceName string) (*model.ServiceInfo, error) {
	req, err := http.NewRequest(http.MethodGet, addr+"/api/v1/service?name="+serviceName, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytedata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp body: %v", err)
	}
	serviceInfo := &model.ServiceInfo{}
	if err := json.Unmarshal(bytedata, serviceInfo); err != nil {
		return nil, fmt.Errorf("unmarshal json %v", err)
	}
	return serviceInfo, nil
}

func DeleteServiceInfo(addr string, serviceName string) error {
	req, err := http.NewRequest(http.MethodDelete, addr+"/api/v1/service?name="+serviceName, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
