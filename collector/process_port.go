package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/shirou/gopsutil/process"
	"strconv"
)

const ProcessPorts = "process_port"

var ProcessPortDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, ProcessPorts, "check"),
	"Process_Port Check Information",
	[]string{"app","port"}, nil,
)

type ProcessPort struct {
}

func init() {
	registerCollector("process_port", defaultEnabled, NewProcessPort)
}

func NewProcessPort() (Collector, error) {
	return &ProcessPort{}, nil
}

func (pp *ProcessPort) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType
	var value float64
	err := pp.WriteAllListenToFile()
	if err != nil {
		log.Error("collect data faild,err:", err)
		return err
	}
	allprocessport, err1 := pp.ReadAllListenToStr()
	if err1 != nil {
		log.Error("read data faild,err:", err1)
		return err1
	}

	for k, v := range allprocessport {
		metricType = prometheus.CounterValue
		for p, cmds := range v {
			pid, err2 := strconv.ParseInt(p, 10, 32)
			if err2 != nil {
				log.Errorf("read data faild,err:%s", err2.Error())
				return err2
			}
			islive, err3 := process.PidExists(int32(pid))
			if err3 != nil {
				log.Errorf("read data faild,err:%s", err3.Error())
				return err3
			}
			if islive {
				value = 1
				log.Debug(pp.Formatlables(cmds), value)
				ch <- prometheus.MustNewConstMetric(
					ProcessPortDesc,
					metricType, value, pp.Formatlables(cmds),k,
				)
			} else {
				value = 0
				log.Debug(pp.Formatlables(cmds), value)
				ch <- prometheus.MustNewConstMetric(
					ProcessPortDesc,
					metricType, value, pp.Formatlables(append(cmds)),k,
				)
			}
		}

	}
	return nil
}
