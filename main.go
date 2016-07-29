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
package main

import (
	"github.com/intelsdi-x/snap-plugin-collector-schedstat/schedstat"
	"github.com/intelsdi-x/snap/control/plugin"
	"os"
)

func main() {
	// Start starts a plugin where:
	// PluginMeta - base information about plugin
	// Plugin - CollectorPlugin, ProcessorPlugin or PublisherPlugin
	// requestString - plugins arguments (marshaled json of control/plugin Arg struct)
	// returns an error and exitCode (exitCode from SessionState initilization or plugin termination code)
	//func Start(m *PluginMeta, c Plugin, requestString string) (error, int) {

	//PluginMeta, Plugin:Collector, requestString
	//this.meta, this.new()
	plugin.Start(schedstat.Meta(),
		new(schedstat.Schedstat), os.Args[1])
}
