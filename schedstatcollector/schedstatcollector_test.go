//
// +build small

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

package schedstatcollector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedStatPlugin(t *testing.T) {
	pwd, _ := os.Getwd()
	procPath := filepath.Join(pwd, "proc")

	config := plugin.Config{
		"procPath": procPath,
	}
	Convey("Create SchedStat Collector", t, func() {
		schedstatCol := SchedstatCollector{}
		Convey("So SchedStat should not be nil", func() {
			So(schedstatCol, ShouldNotBeNil)
		})
		Convey("So SchedStat should be of MongoDB type", func() {
			So(schedstatCol, ShouldHaveSameTypeAs, SchedstatCollector{})
		})
		Convey("SchedStat.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := schedstatCol.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a plugin.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})
	Convey("Get Metric SchedStat Types", t, func() {
		schedstatCol := SchedstatCollector{}
		var cfg = plugin.Config{}
		metrics, err := schedstatCol.GetMetricTypes(cfg)
		So(err, ShouldBeNil)
		So(len(metrics), ShouldResemble, 8)
	})
	Convey("Collect Metrics", t, func() {
		schedstatCol := SchedstatCollector{}

		mts := []plugin.Metric{}

		for _, value := range schedstatStructure {
			mts = append(mts, plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Class, Name).AddStaticElement("cpu").AddDynamicElement("socket_id", "id of socket").AddDynamicElement("cpu_id", "id of cpu").AddStaticElement(value), Config: config})
		}
		metrics, err := schedstatCol.CollectMetrics(mts)
		So(err, ShouldBeNil)
		So(len(metrics), ShouldResemble, 64)
		So(metrics[0].Data, ShouldNotBeNil)
		So(metrics[0].Namespace.Strings()[nsCPUPosition], ShouldResemble, "cpu0")
		So(metrics[0].Namespace.Strings()[nsSocketPosition], ShouldResemble, "1")
		So(metrics[45].Namespace.Strings()[nsCPUPosition], ShouldResemble, "cpu5")
		So(metrics[45].Namespace.Strings()[nsSocketPosition], ShouldResemble, "2")
	})
	Convey("Strings to Ints if slice contains non valid strings", t, func() {
		inputSlice := []string{"one", "1", "2"}
		outputSlice, err := stringsToInts(inputSlice)
		So(err, ShouldNotBeNil)
		So(outputSlice, ShouldResemble, []int{})
	})
}
