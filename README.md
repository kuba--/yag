YAG (Yet Another Graphite)
==========================


<img alt="yak" src="http://www.bluebison.net/sketchbook/2010/0110/monkey-riding-a-yellow-yak.png" />
<sup>image source: http://www.bluebison.net/sketchbook/2010/0110/monkey-riding-a-yellow-yak.png</sup>


"Graphite is actually a bit of a niche application. Specifically, it is designed to handle numeric time-series data. For example, Graphite would be good at graphing stock prices because they are numbers that change over time."
[http://graphite.wikidot.com/faq#toc0]

YAG is yet another graphite. I started this project mostly because I wanted to learn something...
Current version is a minimalistic implementation of graphite's engine (yes, just engine - for rendering charts I recommand 3rd party tool - giraffe).
Right now supports only basic API functions like: sumSeries, divSeries and diffSeries.

YAG contains 3 independent components which follow Unix philosophy - "do one thing and do it well":
* listener - stores datapoints received from "StatsD server"
* webserver - serves datapoints for dashboard clients (implements graphite json rest api)
* ttl - daemon which removes datapoints older than set number of seconds

For metrics management I recommand *phpRedisAdmin*. Here is my fork a little bit modified:

	https://github.com/kuba--/phpRedisAdmin/tree/yag


For metrics dashboard, I recommand *giraffe*. Here is my fork:

	https://github.com/kuba--/giraffe
<img alt="giraffe" src="https://raw.github.com/kuba--/giraffe/master/img/snapshot.png" />
<sup>snapshot taken from giraffe demo webpage</sup>



# Installing

## Dependencies
	
- YAG was implemented in "Go", so install golang first [http://golang.org/doc/install].

- YAG uses redis database (check redis.conf for configuration details) to store datapoints, so you need to install redis database [http://redis.io/download] on your DB server. Redis version >= `2.6.0` required in order to load lua scripts

## Compiling
	
- Set up $GOPATH (e.g.: $GOPATH=$HOME/workspace).
	
- Now you can compile YAG's code. Redis' driver (redix) and glog library will be installed automatically.

		cd $GOPATH

		go get github.com/kuba--/yag/listener
		go get github.com/kuba--/yag/webserver
		go get github.com/kuba--/yag/ttl
		

- If you already downloaded/cloned a code, you would be able to use make command-line tool.
		
		# just compile
		make

		# check deploy target if you want to compile and deploy
		make deploy
		
		# if you installed golang with cross-compile flag, you would be able to compile yag for any platform (e.g. linux)
		make -e GOOS=linux GOARCH=amd64 
		
		
	

- Executable files are in $GOPATH/bin directory.

# Usage

## Run

-  You can run listener, webserver and ttl with following flags:
		
		-f="config.json": path to the config file

		-alsologtostderr=false: log to standard error as well as files
		-log_dir="": If non-empty, write log files in this directory
		-logtostderr=false: log to standard error instead of files
		-stderrthreshold=0: logs at or above this threshold go to stderr
		
		

- Example:

		$ ./listener -f=./config.json -log_dir=./logs -stderrthreshold=INFO



## Profiling

- You can also profiling listener and webserver adding following flags:

		-cpuprofile="": write cpu profile to file
		-memprofile="": write memory profile to this file
		
More about profiling go programs you can find here: 
	http://blog.golang.org/profiling-go-programs


## Configuring

- Configuration file (e.g. config.json):

**Note 1:** remove comments before using this config

	{
		"DB":{                          // Database section
			"Addr":"localhost:6379",     // address and port of Redis DB
			"Timeout":30,                // timeout per connection (in seconds)
			"MaxClients":30              // maximum number of clients in DB connection pool
		},
		"Metrics":{                     // Metrics section
			"GetScript":"get.lua",       // relative path to get script
			"AddScript":"add.lua",       // relative path to add script					
			"TtlScript":"ttl.lua",       // relative path to ttl script
			"TTL":86600,                 // time to live per metric (in seconds)
			"ConsolidationStep":60,      // consolidate datapoints per step (in seconds).
			"ConsolidationFunc":"avg"    // function used to consolidate datapoints (valid function names: sum, avg, min, max)
		},
		"Listener":{                    // Listener server section 
			"Addr":":2003"               // local address and port
		},
		"Webserver":{                   // Webserver section
			"Addr":":8080",              // local address and port
			"Timeout":30                 // timeout per connection (in seconds)
		},
		"TTL":{                         // TTL daemon section
			"Tick":12                    // timers tick (in seconds)
		}
	}


**Note 2 (Webserver):** if you remove properties: "ConsolidationStep", "ConsolidationFunc" from config file, webserver will not consolidate datapoints.

**Note 3 (Webserver):** if you add "maxDataPoints" parameter > 0 then "ConsolidationStep" can be changed by webserver to return _around_ "maxDataPoints". For instance, for following data points received from StatsD daemon:


	[0,1384613389],[0,1384613399],[0,1384613409],[0.5,1384613419],[0.75,1384614209]


 with ConsolidationFunc: "avg" and maxDataPoints=7, webserver will output following data points:


	[0.12, 1384613389]
	[null, 1384613509]
	[null, 1384613629]
	[null, 1384613749]
	[null, 1384613869]
	[null, 1384613989]
	[0.75, 1384614109]



 
# Copyright and licensing

Copyright 2013 *Kuba Podgorski*

Unless otherwise noted, the source files are distributed under the 
*MIT License* found in the LICENSE file.
