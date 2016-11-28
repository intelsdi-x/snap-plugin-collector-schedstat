/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2016 Intel Corporation

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

package schedstatcollector

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	//Name of the plugin
	Name = "schedstat"
	// Vendor name
	Vendor = "intel"
	//Class of the collector
	Class = "proc"
	//Version of the collector
	Version = 1

	nsMetricPostion    = 3
	nsCPUPosition      = 4
	nsSocketPosition   = 5
	nsSubMetricPostion = 6
	//Schedstat structure
	schedYield          = 0
	scheduleCalled      = 1
	scheduleLeftIdle    = 2
	tryToWakeUp         = 3
	tryToWakeUpLocalCPU = 4
	jiffiesRunning      = 5
	jiffiesWaiting      = 6
	timeslices          = 7
)

var schedstatPath = "schedstat"

/*
//sched-stat.txt version 15 for 4.6 kernel
http://lxr.free-electrons.com/source/Documentation/scheduler/sched-stats.txt
*/
var schedstatStructure = []string{
	"schedYield",
	"scheduleCalled",
	"scheduleLeftIdle",
	"tryToWakeUp",
	"tryToWakeUpLocalCPU",
	"jiffiesRunning",
	"jiffiesWaiting",
	"timeslices",
}

//SchedstatCollector name
type SchedstatCollector struct {
}

//CPU struct for cpu measurments
type CPU struct {
	//CPUID id of the Thread
	CPUID string
	//SocketID id of the Socket
	SocketID string
	//Data contains measurment data
	Data []int
}

//GetConfigPolicy method to set configurable items for plugin
func (SchedstatCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	configKey := []string{Vendor, Class, Name}
	policy.AddNewStringRule(configKey, "procPath", false, plugin.SetDefaultString(schedstatPath))

	return *policy, nil

}

//GetMetricTypes method makes values from Global config available
//returns namespaes for my metrics.
func (SchedstatCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	var metrics []plugin.Metric
	ns := plugin.NewNamespace(Vendor, Class, Name)
	for _, value := range schedstatStructure {
		metrics = append(metrics, createMetric(ns.AddStaticElement("cpu").AddDynamicElement("cpu_id", "id of cpu").AddDynamicElement("socket_id", "id of socket").AddStaticElement(value)))
	}
	return metrics, nil
}

/*
https://www.kernel.org/doc/Documentation/scheduler/sched-stats.txt
Version 15 of schedstats dropped counters for some sched_yield:
yld_exp_empty, yld_act_empty and yld_both_empty. Otherwise, it is
identical to version 14.
*/

//CollectMetrics method gathers schedstat data from /proc/schedstat
func (SchedstatCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	procPath, err := mts[0].Config.GetString("procPath")
	if err != nil {
		return metrics, err
	}
	filename := filepath.Join(procPath, schedstatPath)

	schedStatData, err := readSchedstat(filename)
	if err != nil {
		return metrics, err
	}

	for _, mt := range mts {
		ns := mt.Namespace
		switch ns.Strings()[nsSubMetricPostion] {
		case "schedYield":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[schedYield], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}

		case "scheduleCalled":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[scheduleCalled], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		case "scheduleLeftIdle":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[scheduleLeftIdle], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		case "tryToWakeUp":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[tryToWakeUp], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		case "tryToWakeUpLocalCPU":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[tryToWakeUpLocalCPU], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		case "jiffiesRunning":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[jiffiesRunning], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		case "jiffiesWaiting":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[jiffiesWaiting], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		case "timeslices":
			for _, m := range schedStatData {
				newNs := copyNamespace(mt)
				metric := createSchedstatMeasurement(mt, m.Data[timeslices], newNs, m.SocketID, m.CPUID)
				metrics = append(metrics, metric)
			}
		}

	}
	return metrics, nil
}

func readSchedstat(fileName string) ([]CPU, error) {
	returnValue := []CPU{}
	var socketID int
	byteData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return returnValue, err
	}
	lines := strings.Split(string(byteData), "\n")
	cpure := regexp.MustCompile("cpu\\d?")
	for i, x := range lines[:len(lines)-1] {
		fieldsLine := strings.Fields(x)
		if cpure.MatchString(fieldsLine[0]) {
			cpuMask := strings.Fields(lines[i+1])[1][0]
			if cpuMask == byte(49) {
				socketID++
			}
			cpuData, err := stringsToInts(fieldsLine[1:])
			if err != nil {
				return returnValue, err
			}
			cpu := CPU{
				CPUID:    fieldsLine[0],
				SocketID: strconv.Itoa(socketID),
				Data:     cpuData,
			}
			returnValue = append(returnValue, cpu)
		}
	}
	return returnValue, nil

}
