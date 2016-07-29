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

package schedstat

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConfigPolicy(t *testing.T) {
	Convey("Testing GetConfigPolicy", t, func() {
		c := new(Schedstat)
		So(c, ShouldNotBeNil)
		cfg, err := c.GetConfigPolicy()
		So(err, ShouldEqual, nil)
		So(cfg, ShouldNotBeNil)
	})
}

func TestMeta(t *testing.T) {
	Convey("Testing Meta", t, func() {
		c := new(Schedstat)
		So(c, ShouldNotBeNil)
		meta := Meta()
		So(meta, ShouldNotBeNil)
		So(meta.Name, ShouldEqual, "schedstat")
		So(meta.Version, ShouldEqual, 1)
		So(meta.Type, ShouldEqual, plugin.CollectorPluginType)
	})
}

func TestGetSchedstatVersion(t *testing.T) {
	c := new(Schedstat)
	Convey("Get current schedstat version without path.", t, func() {
		version, err := c.getSchedstatVersion("")
		So(version, ShouldEqual, "15")
		So(err, ShouldEqual, nil)
	})
	Convey("Get schedstat version from /totalRandom", t, func() {
		version, err := c.getSchedstatVersion("/totalRandom")
		So(err, ShouldNotBeNil)
		So(version, ShouldEqual, "")
	})
	schedstatPath = "proc/schedstat"
	Convey("Get schedstat version from mock data", t, func() {
		version, err := c.getSchedstatVersion(schedstatPath)
		So(err, ShouldEqual, nil)
		So(version, ShouldEqual, "15")
	})
	schedstatPath = "proc/schedstat14"
	Convey("Get schedstat version from proc/schedstat14", t, func() {
		version, err := c.getSchedstatVersion(schedstatPath)
		So(err, ShouldEqual, nil)
		So(version, ShouldEqual, "14")
	})
	schedstatPath = "proc/negSchedstat"
	Convey("Get schedstat version from proc/negSchedstat", t, func() {
		version, err := c.getSchedstatVersion(schedstatPath)
		So(err, ShouldEqual, nil)
		So(version, ShouldEqual, "")
	})
}

func TestGetMetricTypes(t *testing.T) {
	var version string
	cfg := plugin.NewPluginConfigType()
	c := new(Schedstat)
	Convey("Get namespace metrics for version 14", t, func() {
		version = "14"
		vers, err := c.getSchedstatVersion("proc/schedstat14")
		So(vers, ShouldEqual, version)
		So(err, ShouldEqual, nil)
		mts, err := c.GetMetricTypes(cfg)
		metricNamespace := []string{}
		//	//get namespace elements from the metrics
		for _, e := range mts {
			metricNamespace = append(metricNamespace, e.Namespace().String())
		}

		So(err, ShouldEqual, nil)
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/schedYield")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/jiffiesRunning")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/jiffiesWaiting")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/tryToWakeUpLocalCpu")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/scheduleCalled")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/scheduleLeftIdle")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/timeslices")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/tryToWakeUp")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/emptyBothQueue")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/emptyActiveQueue")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/emptyExpiredQueue")
		So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/switchedExpiredQueue")

	})
	Convey("Get namespace metrics for version 15", t, func() {
		version = "15"
		//be sure to set the schedstatPath
		vers, err := c.getSchedstatVersion("proc/schedstat")
		So(vers, ShouldEqual, version)
		So(err, ShouldEqual, nil)
		mts, err := c.GetMetricTypes(cfg)
		metricNamespace := []string{}
		//get namespace elements from the metrics
		for _, e := range mts {
			metricNamespace = append(metricNamespace, e.Namespace().String())
		}
		Convey("Test namespace elements positive", func() {
			So(err, ShouldEqual, nil)
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/schedYield")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/jiffiesRunning")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/jiffiesWaiting")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/tryToWakeUpLocalCpu")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/scheduleCalled")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/scheduleLeftIdle")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/timeslices")
			So(metricNamespace, ShouldContain, "/intel/proc/schedstat/cpu/*/tryToWakeUp")

			Convey("Test namespace elements negative", func() {
				So(metricNamespace, ShouldNotContain, "/intel/proc/schedstat/cpu/*/emptyBothQueue")
				So(metricNamespace, ShouldNotContain, "/intel/proc/schedstat/cpu/*/emptyActiveQueue")
				So(metricNamespace, ShouldNotContain, "/intel/proc/schedstat/cpu/*/emptyExpiredQueue")
				So(metricNamespace, ShouldNotContain, "/intel/proc/schedstat/cpu/*/switchedExpiredQueue")
			})
		})

	})
	Convey("Collect metrics for version 15", t, func() {
		version = "15"
		cfgNode := cdata.NewNode()
		metrics := []plugin.MetricType{{
			Namespace_: core.NewNamespace("intel",
				"proc", "schedstat", "cpu").
				AddDynamicElement("n", "number of cpu").
				AddStaticElement("schedYield"),
			Config_: cfgNode,
		}}
		cpuInfoPath = "proc/cpuinfo"
		mts, err := c.CollectMetrics(metrics)

		Convey("Should Return 8 Metrics from Schedstat", func() {
			So(err, ShouldEqual, nil)
			var expectedType int64
			So(len(mts), ShouldResemble, 8)
			So(mts[0].Data_, ShouldNotBeEmpty)
			So(mts[0].Data_, ShouldHaveSameTypeAs, expectedType)

		})

	})
	Convey("Collect metrics for version 14", t, func() {
		version = "14"
		cfgNode := cdata.NewNode()
		metrics := []plugin.MetricType{{
			Namespace_: core.NewNamespace("intel",
				"proc", "schedstat", "cpu").
				AddDynamicElement("n", "number of cpu").
				AddStaticElement("jiffiesRunning"),
			Config_: cfgNode,
		}}
		vers, err := c.getSchedstatVersion("proc/schedstat14")
		So(vers, ShouldEqual, version)
		So(err, ShouldEqual, nil)
		mts, err := c.CollectMetrics(metrics)
		So(err, ShouldEqual, nil)
		Convey("Should Return 1 Metrics from Schedstat", func() {
			So(err, ShouldEqual, nil)
			var expectedType int64
			So(len(mts), ShouldResemble, 1)
			So(mts[0].Data_, ShouldNotBeEmpty)
			So(mts[0].Data_, ShouldHaveSameTypeAs, expectedType)

		})

	})
}
