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
	Spath string // 脚本所在目录
}

var (
	scriptPath = kingpin.Flag(
		"spath",
		"coustom script path",
	).Default("/opt/node_exporter/").String()
)

func init() {
	registerCollector("coustomscript", defaultEnabled, NewCoustomScriptCollector)
}

func NewCoustomScriptCollector() (Collector, error) {
	return &CoustomScriptCollector{
		Spath: *scriptPath,
	}, nil
}

func (c *CoustomScriptCollector) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType
	files, err := filepath.Glob(c.Spath + "/*")
	log.Debugf("scripts path %s", c.Spath)
	if err != nil {
		log.Fatalf("get scripts file list faile %s", err.Error())
	}
	log.Debugf("scripts file %v", files)
	for index := range files {
		metricType = prometheus.CounterValue

		file := files[index]
		if err = os.Chmod(file, 0755); err != nil {
			log.Fatalf("chmod +x %s faile err:%s", file, err.Error())
		}
		result, err1 := c.RunCommand(file)
		log.Debugf("result %s command %s", result, file)
		if err1 != nil || result == "" {
			log.Errorf("command run fail,result: %s,file: %s", result, file)
			continue
		}
		results := strings.Split(result, "\n")

		for k := range results {
			if results[k] == "" {
				continue
			}
			log.Debugf("result[k]:%s", results[k])
			resultArr := strings.Split(results[k], "=")
			if len(resultArr) < 2 {
				continue
			}
			key, value := strings.Replace(strings.Replace(resultArr[0], " ", "_", -1), "-", "_", -1), resultArr[1]
			log.Debugf("key %s value %s", key, value)
			values, err := strconv.ParseFloat(strings.Replace(value, "\n", "", -1), 64)
			if err != nil {
				log.Errorf("ParseFloat value %s faile,err:%s", value, err.Error())
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
