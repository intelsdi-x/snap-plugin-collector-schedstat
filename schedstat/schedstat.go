/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//Package schedstat gathers information from kernel scheduling: /proc/schedstat
//At the moment it works for schedstat version 15 and collects only CPU
//related information.

package schedstat

import (
	"bufio" //ReadLine, ReadString
	"fmt"
	"io/ioutil"
	"os"     //cd
	"regexp" //version
	"strconv"
	"strings"
	"time" //time.Duration

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/core" //namespace
	//mandatory packages that any plugin must use
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	Name            = "schedstat"
	Vendor          = "Intel"
	Class           = "proc"
	Version         = 1
	Type            = plugin.CollectorPluginType
	Exclusive       = true
	Unsecure        = true
	CacheTTL        = time.Millisecond * 5
	RoutingStrategy = plugin.DefaultRouting
)

var (
	AcceptedContentTypes = []string{plugin.SnapGOBContentType}
	ReturnedContentTypes = []string{plugin.SnapGOBContentType}
	ConcurrencyCount     = 1
	schedstatPath        = "/proc/schedstat"
	schedstatVersion     string
	schedstatStructure   []string
	cpuInfoPath          = "/proc/cpuinfo"
)

type Schedstat struct {
}

//getSchedStatVersion() method reads /proc/schedstat/'s first line and
//parses \d+$ to get the current schedstat version, which is returned as a
//string.
func (c *Schedstat) getSchedstatVersion(path string) (string, error) {
	//head -n1 $schedstat
	if path == "" {
		path = schedstatPath
	}
	fmt.Println(path)
	fmt.Println(schedstatPath)
	fh, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	//Read first line of the file, get schedstat version information
	scanner.Scan()
	line := scanner.Text()
	re := regexp.MustCompile(`\d+$`)
	version := re.FindString(line)
	schedstatVersion = version
	schedstatPath = path

	return version, nil
}

/*
//sched-stat.txt version 15 for 4.6 kernel
http://lxr.free-electrons.com/source/Documentation/scheduler/sched-stats.txt

//First field is a sched_yield() statistics
1. # sched_yield() => sched_yield

//schedule() statistics
2. Legacy array expiration count, always zero
3. # times schedule () called
4. # times schedule() left the CPU iddle

//try_to_wake_up() statistics
5. # try_to_wake_up() called => wake_up_remote
6. # try_to_wake_up() called A local cpu  => wake_up_local

//Statistics describing scheduling latency
7. #jiffies A time spent running task => running
8. #jiffies A time spent waiting to run	=> waiting
9. #timeslices	=>tasks

*/

//GetConfigPolicy method to set configurable items for plugin
func (c *Schedstat) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

func (c *Schedstat) createDynamicNamespace() ([]plugin.MetricType, error) {
	pluginMetric := []plugin.MetricType{}
	for _, e := range schedstatStructure {
		//Namespace for cpu information
		metricType := plugin.MetricType{
			Namespace_: core.NewNamespace("intel",
				"proc", "schedstat", "cpu").
				AddDynamicElement("n", "number of cpu").
				AddStaticElement(e),
		}
		pluginMetric = append(pluginMetric, metricType)
		//Create namespace for dynamic socket information
		metricType = plugin.MetricType{
			Namespace_: core.NewNamespace("intel",
				"proc", "schedstat", "socket").
				AddDynamicElement("n", "number of cpu sockets").
				AddStaticElement(e),
		}
		pluginMetric = append(pluginMetric, metricType)

		//Namespace for all gathered information
		metricType = plugin.MetricType{
			Namespace_: core.NewNamespace("intel",
				"proc", "schedstat", "system").
				AddStaticElement(e),
		}
		pluginMetric = append(pluginMetric, metricType)
	}

	return pluginMetric, nil
}

//GetMetricTypes() method makes values from Global config available
//returns namespaes for my metrics.
func (c *Schedstat) GetMetricTypes(cfg plugin.ConfigType) (
	[]plugin.MetricType, error) {
	var err error
	schedstatVersion, err = c.getSchedstatVersion(schedstatPath)
	if err != nil || schedstatVersion == "" {
		return []plugin.MetricType{}, err
	}
	var schedStats []string
	if schedstatVersion == "15" {
		schedStats = []string{
			"schedYield",
			"scheduleCalled",
			"scheduleLeftIdle",
			"tryToWakeUp",
			"tryToWakeUpLocalCpu",
			"jiffiesRunning",
			"jiffiesWaiting",
			"timeslices",
		}
	} else if schedstatVersion == "14" {
		schedStats = []string{
			"emptyBothQueue",    //these are different
			"emptyActiveQueue",  //values compared to
			"emptyExpiredQueue", //version 15
			"schedYield",
			"switchedExpiredQueue", //
			"scheduleCalled",
			"scheduleLeftIdle",
			"tryToWakeUp",
			"tryToWakeUpLocalCpu",
			"jiffiesRunning",
			"jiffiesWaiting",
			"timeslices",
		}
	}
	schedstatStructure = schedStats
	pluginMetric, err := c.createDynamicNamespace()
	return pluginMetric, err
}

//Meta function returns plugin metainformation.
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		//string plugin.SnapGOBcontentType
		[]string{},
		[]string{plugin.SnapJSONContentType},
		plugin.ConcurrencyCount(1),
	//	AcceptedContentTypes,
	//	ReturnedContentTypes,
	)
}

/*
https://www.kernel.org/doc/Documentation/scheduler/sched-stats.txt
Version 15 of schedstats dropped counters for some sched_yield:
yld_exp_empty, yld_act_empty and yld_both_empty. Otherwise, it is
identical to version 14.
*/

//CollecgtMetrics method gathers schedstat data from /proc/schedstat
func (c *Schedstat) CollectMetrics(metric []plugin.MetricType) ([]plugin.MetricType, error) {
	return c.parseMetric(metric)
}

//createSocketNamespace creates socket namespace /intel/proc/schedstat/socket/<jiffies>
func (c *Schedstat) createSocketNamespace(data map[int]map[int][]int) ([]plugin.MetricType, error) {
	pluginMetric := []plugin.MetricType{}
	for j := range data {
		for _, i := range schedstatStructure {
			metricType := plugin.MetricType{
				Namespace_: core.NewNamespace("intel",
					"proc", "schedstat", "socket").
					AddStaticElement(strconv.Itoa(j)).
					AddStaticElement(i),
			}
			pluginMetric = append(pluginMetric, metricType)
		}
	}
	return pluginMetric, nil
}

//readCpuinfo reads /proc/cpuinfo to get cpu sockets
func (c *Schedstat) readCpuinfo() (map[int]map[int][]int, error) {
	//map[keyType]ValueType
	var data map[int]map[int][]int
	data = make(map[int]map[int][]int)
	byteData, err := ioutil.ReadFile(cpuInfoPath)
	if err != nil {
		return data, err
	}
	dt := string(byteData)
	//FIXME split by thread info, this is kind of stupid split, \n\n would be better
	fields := strings.Split(dt, "power management:")
	for _, e := range fields {
		var socketID = -1 //socket per thread
		var coreID = -1
		//socket ID as physical id
		//regexp:	physical id	: 0
		re := regexp.MustCompile("physical\\s+id\\s+:\\s*([\\d])")
		find := re.FindStringSubmatch(e)
		if len(find) == 0 {
			continue //split, doesn't contain the information
			//See above fixme
		} else {
			sid := find[len(find)-1]
			id, _ := strconv.Atoi(sid)
			socketID = id
			if _, exists := data[id]; !exists {
				data[id] = make(map[int][]int)
			}
		}
		//actual #cpu
		//regexp:	 core id		: 0
		re = regexp.MustCompile("core\\s+id\\s+:\\s*([\\d])")
		find = re.FindStringSubmatch(e)
		if len(find) == 0 {
			log.WithFields(log.Fields{
				"data": cpuInfoPath,
			}).Error("Failed to parse core id from cpuinfo.")
			return data, fmt.Errorf("Failed to parse core id from %s", cpuInfoPath)
		}
		if socketID < 0 {
			log.WithFields(log.Fields{
				"data": cpuInfoPath,
			}).Error("Socket ID error in cpuinfo.")
			return data, fmt.Errorf("Socket ID error in cpuinfo %s", cpuInfoPath)
		}
		sid := find[len(find)-1]
		id, _ := strconv.Atoi(sid)
		if _, exists := data[socketID][id]; !exists {
			data[socketID][id] = []int{}
			coreID = id
		} else {
			//core id already exists
			coreID = id
		}

		//find thread id
		//regexp: 	processor	: 0
		re = regexp.MustCompile("processor\\s+:\\s*([\\d])")
		find = re.FindStringSubmatch(e)
		if len(find) == 0 {
			log.WithFields(log.Fields{
				"data": cpuInfoPath,
			}).Error("Failed to parse processor from cpuinfo.")
			return data, fmt.Errorf("Failed to parse processor from %s", cpuInfoPath)
		} else {
			sid := find[len(find)-1] //get data from last element
			if socketID < 0 || coreID < 0 {
				log.WithFields(log.Fields{
					"data":   cpuInfoPath,
					"socket": socketID,
					"cpu":    coreID,
					"id":     sid,
				}).Error("Socket ID or cpu id error in cpuinfo.")
				return data, fmt.Errorf("Socket ID or cpu ID error in cpuinfo %s", cpuInfoPath)
			} else {
				id, _ := strconv.Atoi(sid)
				//map[int]map[int][]int
				data[socketID][coreID] = append(data[socketID][coreID], id)
			}
		}
	}
	return data, nil
}

//func (c *Schedstat) pointerMethod {} // method on pointer
//func (c Schedstat) valueMethod() //method on value

//func (c *Schedstat) parseMetric(dynamicMetric []plugin.MetricType, pluginmetric []plugin.MetricType) (
func (c *Schedstat) parseMetric(pluginmetric []plugin.MetricType) (
	[]plugin.MetricType, error) {
	// /intel/proc/schedstat/socket/<dynamic>/schedYield
	//[3] is where the cpu/socket/system ns is
	//fields threadN->jiffies->239578

	var collectedMetric []plugin.MetricType

	//[cpuN][JiffiesN}[value]
	fields, err := c.collectSchedstatData()
	if err != nil {
		return collectedMetric, err
	}

	//Calculate dt before creating ns
	cpuinfoDt, err := c.readCpuinfo()
	if err != nil {
		return collectedMetric, err
	}
	socketDt := c.createAggrecateSocket(fields, cpuinfoDt)
	systemDt := c.createAggrecateCPU(fields)

	for _, metric := range pluginmetric {
		currentNs := metric.Namespace()
		domain := currentNs.Strings()[3]
		//last element /jiffies => "", "jiffies"
		NsElement := metric.Namespace()[len(metric.Namespace())-1:]
		elements := strings.Split(NsElement.String(), "/")
		lastNsElements := elements[len(elements)-1:]
		//[lastNsElements] eq jiffies, timeslices etc
		if strings.Contains(domain, "cpu") {
			for cpuN := range fields {
				for _, lastNsElem := range lastNsElements {
					lastNsElem := string(lastNsElem)
					if b, dynIndex := currentNs.IsDynamic(); b == true {
						dynNs := make(core.Namespace, len(currentNs))
						copy(dynNs, currentNs)
						for _, j := range dynIndex {
							dynNs[j].Value = strconv.Itoa(cpuN)
						}
						mets := plugin.MetricType{
							Namespace_: dynNs,
							Timestamp_: time.Now(),
							Data_:      fields[cpuN][lastNsElem],
						}
						collectedMetric = append(collectedMetric, mets)
					}
				}
			}
		} else if strings.Contains(domain, "socket") {
			//parse socket, information ready, no dynamic parts
			//information gathered in readCpuInfo method
			for _, lastNsElem := range lastNsElements {
				lastNsElem := string(lastNsElem)
				if b, dynIndex := metric.Namespace().IsDynamic(); b == true {
					dynNs := make(core.Namespace, len(metric.Namespace()))
					copy(dynNs, metric.Namespace())
					for socketN := range socketDt {
						for _, j := range dynIndex {
							dynNs[j].Value = strconv.Itoa(socketN)
						}
						mets := plugin.MetricType{
							Namespace_: dynNs,
							Timestamp_: time.Now(),
							Data_:      socketDt[socketN][lastNsElem],
						}
						collectedMetric = append(collectedMetric, mets)
					}
				}
			}
		} else if strings.Contains(domain, "system") {
			for _, lastNsElem := range lastNsElements {
				lastNsElem := string(lastNsElem)
				dt := systemDt[lastNsElem]
				mets := plugin.MetricType{
					Namespace_: metric.Namespace(),
					Timestamp_: time.Now(),
					Data_:      dt,
				}
				collectedMetric = append(collectedMetric, mets)
			}
		}
	}
	return collectedMetric, nil
}

//createAggrecateSocket method sums information per cpu socket.
func (c *Schedstat) createAggrecateSocket(fields map[int]map[string]int64,
	socket map[int]map[int][]int) map[int]map[string]int64 {

	var dt map[int]map[string]int64
	dt = make(map[int]map[string]int64)

	//cpuinfo map
	//fields: [threadN][jiffies][value]
	//socket:  [socket][cpu][thread1, threadN]
	//	map[0:map[0:[0 4] 1:[1 5] 2:[2 6] 3:[3 7]]]
	//[socket][jiffies][value]
	for index, core := range socket {
		dt[index] = make(map[string]int64)
		for _, threads := range core {
			//		[cpu][thread1, threadN]
			for _, j := range threads {
				for a, schedvalue := range fields[j] {
					dt[index][a] = dt[index][a] + schedvalue
				}
			}
		}
	}
	return dt
}

//createAggrecateCPU method gathers all cpuN information together.
func (c *Schedstat) createAggrecateCPU(fields map[int]map[string]int64) map[string]int64 {
	//[cpuN][<jiffies>[value]
	var dt map[string]int64
	dt = make(map[string]int64)
	for _, e := range fields {
		for i, j := range e {
			dt[i] = dt[i] + j
		}
	}
	return dt
}

//collectSchedstatData method collects information from /proc/schedstat
func (c *Schedstat) collectSchedstatData() (map[int]map[string]int64, error) {
	var fields map[int]map[string]int64
	//[cpuN][JiffiesN}[value]
	fields = make(map[int]map[string]int64)
	byteData, err := ioutil.ReadFile(schedstatPath)
	if err != nil {
		return fields, err
	}
	dt := string(byteData)
	lines := strings.Split(dt, "\n")
	for _, x := range lines {
		line := string(x)
		//Not collecting domainN information
		re := regexp.MustCompile("cpu\\d?")
		if re.MatchString(line) {
			allFields := strings.Fields(line)
			re := regexp.MustCompile(`\d?$`)
			cpun := re.FindString(allFields[0])
			n, _ := strconv.Atoi(cpun)
			fields[n] = make(map[string]int64)
			//rm cpuN part
			allFields = append(allFields[:0], allFields[1:]...)
			var fieldData []int64
			//Get schedstat data from splitted txt
			//and convert it to int64
			for _, i := range allFields {
				j, err := strconv.ParseInt(i, 10, 64)
				if err != nil {
					return fields, err
				}
				fieldData = append(fieldData, j)
			}
			//fields contains the data
			if schedstatVersion == "" {
				schedstatVersion, _ = c.getSchedstatVersion(schedstatPath)
			}
			if schedstatVersion == "15" {
				fields[n] = c.version15Data(fieldData[:])
				//dt = c.version15Data(fields, fieldData[:])
			} else if schedstatVersion == "14" {
				fields[n] = c.version14Data(fieldData[:])
			}

		}
	} //for one line
	return fields, nil
}

//version15Data method collect schedstat information for version 15
func (c *Schedstat) version15Data(fields []int64) map[string]int64 {
	//rm Legacy array expiration count, always zero
	fields = append(fields[:1], fields[2:]...)
	var data map[string]int64
	data = make(map[string]int64)
	data["schedYield"] = fields[0]
	data["scheduleCalled"] = fields[1]
	data["scheduleLeftIdle"] = fields[2]
	data["tryToWakeUp"] = fields[3]
	data["tryToWakeUpLocalCpu"] = fields[4]
	data["jiffiesRunning"] = fields[5]
	data["jiffiesWaiting"] = fields[6]
	data["timeslices"] = fields[7]
	return data
}

//version14Data method collect schedstat information for version 14
func (c *Schedstat) version14Data(fields []int64) map[string]int64 {
	var data map[string]int64
	data = make(map[string]int64)
	data["emptyBothQueue"] = fields[0]
	data["emptyActiveQueue"] = fields[1]
	data["emptyExpiredQueue"] = fields[2]
	data["schedYield"] = fields[3]
	data["switchedExpiredQueue"] = fields[4]
	data["scheduleCalled"] = fields[5]
	data["scheduleLeftIdle"] = fields[6]
	data["tryToWakeUp"] = fields[7]
	data["tryToWakeUpLocalCpu"] = fields[8]
	data["jiffiesRunning"] = fields[9]
	data["jiffiesWaiting"] = fields[10]
	data["timeslices"] = fields[11]
	return data
}
