# Snap collector plugin - schedstat

## Collected Metrics

This plugin has the ability to gather the following metrics:


Namespace | Data Type | Description
----------|-----------|-------------------------------------
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/schedYield| int|# of times sched_yield() was called
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/scheduleCalled| int|# of times schedule() was called
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/scheduleLeftIdle| int|# of times schedule() left the processor idle
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/tryToWakeUp| int|# of times try_to_wake_up() was called
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/tryToWakeUpLocalCPU| int|# of times try_to_wake_up() was called to wake up the local cpu
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/jiffiesRunning| int|sum of all time spent running by tasks on this processor (in jiffies)
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/jiffiesWaiting| int|sum of all time spent waiting to run by tasks on this processor (in jiffies)
/intel/proc/schedstat/cpu/{cpu_id}/{socket_id}/timeslices| int|# of timeslices run on this cpu


