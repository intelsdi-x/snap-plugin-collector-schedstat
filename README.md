# snap collector plugin - schedstat

Snap collector plugin schedstat collects Linux kernel scheduling information from
/proc/schedstat (schedstat version 15). 


1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)


## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/) only for building the plugin

### Operating systems
All OSs currently supported by plugin:
* Linux/amd64

### Installation
You can get snap's pre-built binaries for your OS and architecture at [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapd` to access.

### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-schedstat
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-schedstat.git
```

Build the plugin by running make within the cloned repo:
```
$ $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector/schedstat/ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

## Documentation

### Collected Metrics

The base namespace for the metrics is /intel/proc/schedstat/ for more information
see list of collected metrics in [METRICS.md](METRICS.md).
In [kernel.org](https://www.kernel.org/doc/Documentation/scheduler/sched-stats.txt) is 
the documentation for scheduler and schedstat.

### Examples

This example describes how to run schedstat plugin with 
[passthru processor plugin](https://github.com/intelsdi-x/snap/tree/master/plugin/processor/snap-plugin-processor-passthru) and
[file publisher plugin](https://github.com/intelsdi-x/snap-plugin-publisher-file). 

Verify that you have exported a proper `SNAP_PATH`: 

```
$ export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build
```
where the location of $GOPATH can be seen by `go env` command.

To start the snapdaemon:

```
$ $SNAP_PATH/bin/snapd --plugin-trust 0 --log-level 1
```

Load snap plugins from $SNAP_PATH/plugin/:

```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-plugin-processor-passthru
Plugin loaded
Name: passthru
Version: 1
Type: processor
Signed: false
Loaded Time: Fri, 02 Sep 2016 17:48:12 IST
```

Load [publisher plugin](https://github.com/intelsdi-x/snap-plugin-publisher-file).


```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-plugin-publisher-file
Plugin loaded
Name: file
Version: 2
Type: publisher
Signed: false
Loaded Time: Fri, 02 Sep 2016 17:49:19 IST
```

Finally load the schedstat collector:

```
$ $SNAP_PATH/bin/snapctl plugin load $GOPATH/bin/snap-plugin-collector-schedstat
Plugin loaded
Name: schedstat
Version: 1
Type: collector
Signed: false
Loaded Time: Fri, 02 Sep 2016 17:50:28 IST
```

By `snapctl metric list` you can see all the available metrics for collection.
You can specify which of these metrics you collect in the task file.

```
$ $SNAP_PATH/bin/snapctl metric list
NAMESPACE 								 VERSIONS
/intel/proc/schedstat/cpu/*/jiffiesRunning 		 1
/intel/proc/schedstat/cpu/*/jiffiesWaiting 		 1
/intel/proc/schedstat/cpu/*/schedYield 			 1
/intel/proc/schedstat/cpu/*/scheduleCalled 		 1
/intel/proc/schedstat/cpu/*/scheduleLeftIdle 		 1
/intel/proc/schedstat/cpu/*/timeslices 			 1
/intel/proc/schedstat/cpu/*/tryToWakeUp 		 1
/intel/proc/schedstat/cpu/*/tryToWakeUpLocalCpu 	 1
/intel/proc/schedstat/socket/*/jiffiesRunning 		 1
/intel/proc/schedstat/socket/*/jiffiesWaiting 		 1
/intel/proc/schedstat/socket/*/schedYield 		 1
/intel/proc/schedstat/socket/*/scheduleCalled 		 1
/intel/proc/schedstat/socket/*/scheduleLeftIdle 	 1
/intel/proc/schedstat/socket/*/timeslices 		 1
/intel/proc/schedstat/socket/*/tryToWakeUp 		 1
/intel/proc/schedstat/socket/*/tryToWakeUpLocalCpu 	 1
/intel/proc/schedstat/system/jiffiesRunning 		 1
/intel/proc/schedstat/system/jiffiesWaiting 		 1
/intel/proc/schedstat/system/schedYield 		 1
/intel/proc/schedstat/system/scheduleCalled 		 1
/intel/proc/schedstat/system/scheduleLeftIdle 		 1
/intel/proc/schedstat/system/timeslices 		 1
/intel/proc/schedstat/system/tryToWakeUp 		 1
/intel/proc/schedstat/system/tryToWakeUpLocalCpu 	 1

```


Create task:
```
$ $SNAP_PATH/bin/snapctl task create -t $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-schedstat/schedstat/task.json
Using task manifest to create task
Task created
ID: 81759d63-d2db-45aa-b225-7870c85e32dd
Name: Task-81759d63-d2db-45aa-b225-7870c85e32dd
State: Running
```

To see what tasks and their task IDs are:

```
$ $SNAP_PATH/bin/snapctl task list
ID 					 NAME 						 STATE 		 HIT 	 MISS  FAIL 	 CREATED 		 LAST FAILURE
81759d63-d2db-45aa-b225-7870c85e32dd 	 Task-81759d63-d2db-45aa-b225-7870c85e32dd 	 Running 	 32 	 0 	 0 	 5:55PM 9-02-2016 	
```
To stop a specific snap task you have to specify the task ID to stop.

```
$ $SNAP_PATH/bin/snapctl task stop 81759d63-d2db-45aa-b225-7870c85e32dd
Task stopped:
ID: 81759d63-d2db-45aa-b225-7870c85e32dd
```

After this you can stop the snapd (CTRL+c).

### Roadmap 
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-schedstat/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-schedstat/pulls).


## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@ssorsa](https://github.com/ssorsa/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
