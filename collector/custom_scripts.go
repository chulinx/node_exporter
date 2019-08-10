package collector

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const coustomScript = "coustom_script"

type CoustomScriptCollector struct {
}

var (
	scriptPath = kingpin.Flag(
		"spath",
		"coustom script path",
	).Default("/opt/node_export/").String()
)

func init() {
	registerCollector("coustomscript", defaultEnabled, NewCoustomScriptCollector)
}

func NewCoustomScriptCollector() (Collector, error) {
	return &CoustomScriptCollector{}, nil
}

func (c *CoustomScriptCollector) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType
	files, err := filepath.Glob(*scriptPath+"/*")
	log.Debugf("scripts path %s", *scriptPath)
	if err != nil {
		log.Fatalf("get scripts file list faile %s", err.Error())
	}
	log.Debugf("scripts file %v",files)
	for index := range files {
		metricType = prometheus.CounterValue

		file := files[index]
		if err = os.Chmod(file, 0755); err != nil {
			log.Fatalf("chmod +x %s faile err:%s", file,err.Error())
		}
		result, err := c.RunCommand(file)
		log.Debugf("result %s command %s",result,file)
		if err != nil {
			log.Infof("command run fail %s",err.Error())
		}
		key, value := strings.Replace(strings.Split(result, "=")[0]," ","_",-1), strings.Split(result, "=")[1]
		log.Debugf("key %s value %s",key,value)
		values, err := strconv.ParseFloat(strings.Replace(value,"\n","",-1), 64)
		if err != nil {
			log.Fatalf("ParseFloat value %s faile,err:%s", value,err.Error())
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, coustomScript, key),
				fmt.Sprintf("CoustomScript information field %s.", key),
				nil, nil,
			),
			metricType, values,
		)
	}

	return nil
}

func (runc *CoustomScriptCollector) RunCommand(command string) (result string, err error) {
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err1 := cmd.Run()
	if err1 != nil {
		result, err = "", err1
		return
	}
	result, err = out.String(), nil
	return
}
