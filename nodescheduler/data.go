package nodescheduler

import (
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/metrics"
)

var dataMap = map[string]*metrics.NodeInfoRecord{}
var lock sync.Mutex

func getAllNodeLatestMetrics() []*metrics.NodeFullMetric {
	lock.Lock()
	defer lock.Unlock()
	var res []*metrics.NodeFullMetric
	for _, v := range dataMap {
		res = append(res, &(v.Metrics[len(v.Metrics)-1]))
	}
	return res
}

func getRecord(nodeid string) (*metrics.NodeInfoRecord, bool, error) {
	lock.Lock()
	defer lock.Unlock()
	record, ok := dataMap[nodeid]
	if !ok {
		record = &metrics.NodeInfoRecord{ID: nodeid}
	}
	return record, ok, nil
}

func saveRecord(nodeid string, record *metrics.NodeInfoRecord) error {
	lock.Lock()
	defer lock.Unlock()
	dataMap[nodeid] = record
	return nil
}
