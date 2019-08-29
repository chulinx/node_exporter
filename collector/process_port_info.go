package collector

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/tianyazc/zlxGo/zlxmodule/stringfile"
	"strconv"
	"strings"
)

func (pp *ProcessPort) WriteAllListenToFile() error {
	var allListen = map[string]map[string][]string{}
	nets, err := net.Connections("all")
	if err != nil {
		log.Error("error:", err.Error())
		return err
	}
	var cmdslice = []string{}
	for _, fd := range nets {
		if fd.Status == "LISTEN" {
			p, err := process.NewProcess(fd.Pid)
			if err != nil {
				cmdslice = []string{}
			} else {
				cmdstr,_ := p.Name()
				cmdslice = []string{cmdstr}
			}
			allListen[strconv.FormatUint(uint64(fd.Laddr.Port), 10)] = map[string][]string{strconv.FormatInt(int64(fd.Pid), 10): cmdslice}
		}
	}
	// 更新
	fallListen, err := pp.ReadAllListenToStr()
	if err == nil && len(fallListen) != 0 {
		for k, v := range allListen {
			fallListen[k] = v
		}
		data, err1 := json.Marshal(fallListen)
		if err1 != nil {
			log.Error("error1:", err1.Error())
			return err1
		}
		err2 := stringfile.WriteStringToFile(string(data), "/opt/process_port.pp", "rw")
		if err2 != nil {
			log.Error("error2:", err2.Error())
			return err2
		}
	} else {
		log.Info("/opt/process_port.pp read failed,err")
		//fmt.Println(allListen)
		data, err1 := json.Marshal(allListen)
		//fmt.Println(data)
		if err1 != nil {
			log.Error("error1:", err1.Error())
			return err1
		}
		err2 := stringfile.WriteStringToFile(string(data), "/opt/process_port.pp", "rw")
		if err2 != nil {
			log.Error("error2:", err2.Error())
			return err2
		}
	}
	return nil
}

func (pp *ProcessPort) ReadAllListenToStr() (map[string]map[string][]string, error) {
	var allListen = make(map[string]map[string][]string)
	str, err := stringfile.ReadLineToString("/opt/process_port.pp", -1)
	if err != nil {
		log.Errorf("error:%s", err.Error())
		return allListen, err
	}
	err1 := json.Unmarshal([]byte(str), &allListen)
	if err1 != nil {
		//log.Errorf("error: %s",err1.Error() )
		return allListen, err
	}
	return allListen, err
}

func (pp *ProcessPort)Formatlables(a []string) (ret []string) {
	for k := range a {
		ai := strings.Replace(strings.Replace(strings.Replace(strings.Replace(a[k]," ","_",-1),".","_",-1),"-","_",-1),".","_",-1)
		ret = append(ret,ai)
	}

	return
}