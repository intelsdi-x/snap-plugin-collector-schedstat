# Snap collector plugin - schedstat

This plugin gather linux scheduler statistics from /proc/schedstat (Linux 2.6+) for the [Snap telemetry framework](http://github.com/intelsdi-x/snap).


1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Global Config](#global-config)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.7+](https://golang.org/dl/)  - needed only for building
* [Linux] (kernel 2.6+)

### Operating systems
* Linux/amd64

### Installation


#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-schedstat/releasess) page. Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for Snap packages).


#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-schedstat
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-schedstat.git
```

Build the Snap schedstat plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Configure `procPath`: path to procfs filesystem if you are running Snap in container

* Load the plugin and create a task, see example in [Examples](#examples).

### Collected Metrics

List of collected metrics is described in [METRICS.md](METRICS.md).

***Plugin does not support collection of metrics for a single core or socket***
### Examples

Example of running Snap schedstat collector and writing data to file.

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `snapteld -l 1 -t 0 &`

Download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-schedstat/latest/linux/x86_64/snap-plugin-collector-schedstat
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snaptel plugin load snap-plugin-collector-schedstat
$ snaptel plugin load snap-plugin-publisher-file

Create a task manifest file  (exemplary files in [examples/tasks/] (examples/tasks/)):
```yaml
---
  version: 1
  schedule:
    type: "simple"
    interval: "1s"
  max-failures: 10
  workflow:
    collect:
      metrics:
        /intel/proc/schedstat/cpu/*/*/schedYield: {}
        /intel/proc/schedstat/cpu/*/*/scheduleCalled: {}
        /intel/proc/schedstat/cpu/*/*/scheduleLeftIdle: {}
        /intel/proc/schedstat/cpu/*/*/tryToWakeUp: {}
        /intel/proc/schedstat/cpu/*/*/tryToWakeUpLocalCPU: {}
        /intel/proc/schedstat/cpu/*/*/jiffiesRunning: {}
        /intel/proc/schedstat/cpu/*/*/jiffiesWaiting: {}
        /intel/proc/schedstat/cpu/*/*/timeslices: {}
      publish:
        - plugin_name: "file"
          config:
            file: "/tmp/schedstat_metrics.log"
```
Download an [example task file](https://github.com/intelsdi-x/snap-plugin-collector-schedstat/blob/master/examples/tasks/) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-schedstat/master/examples/tasks/schedstat-file.yaml
$ snaptel task create -t schedstat-file.yaml
Using task manifest to create task
Task created
ID: 480323af-15b0-4af8-a526-eb2ca6d8ae67
Name: Task-480323af-15b0-4af8-a526-eb2ca6d8ae67
State: Running
```

See realtime output from `snaptel task watch <task_id>` (CTRL+C to exit)
```
$ snaptel task watch 480323af-15b0-4af8-a526-eb2ca6d8ae67
```

This data is published to a file `/tmp/schedstat_metrics` per task specification

Stop task:
```
$ snaptel task stop 480323af-15b0-4af8-a526-eb2ca6d8ae67
Task stopped:
ID: 480323af-15b0-4af8-a526-eb2ca6d8ae67
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. 

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-schedstat/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-schedstat/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Marcin Spoczynski](https://github.com/sandlbn/)

This software has been contributed by MIKELANGELO & Superfluidity, Horizon 2020 projects co-funded by the European Union. https://www.mikelangelo-project.eu/ http://superfluidity.eu/
## Thank You
And **thank you!** Your contribution, through code and participation, is incredibly important to us.