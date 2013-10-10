YAG (Yet Another Graphite)
==========================

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


<img alt="yak" src="http://www.bluebison.net/sketchbook/2010/0110/monkey-riding-a-yellow-yak.png" />
image source: http://www.bluebison.net/sketchbook/2010/0110/monkey-riding-a-yellow-yak.png


## Installing

# Dependencies
	
- YAG was implemented in "Go", so install golang first [http://golang.org/doc/install].

- YAG uses redis database (check redis.conf for configuration details) to store datapoints, so you need to install redis database [http://redis.io/download] on your DB server. Redis version >= `2.6.0` required in order to load lua scripts

# Compiling
	
- Set up $GOPATH (e.g.: $GOPATH=$HOME/workspace).
	
- Now you can compile YAG's code with redis' driver - redix. 

		cd $GOPATH

		go get github.com/fzzy/radix/redis

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

# Running

-  Start listener, webserver, ttl

		$GOPATH/bin/listener [-f Specify a path to the config file]
	
		$GOPATH/bin/webserver [-f Specify a path to the config file]
	
		$GOPATH/bin/ttl [-f Specify a path to the config file]



## Configuring

- Configuration file (e.g. config.json):

**Note:** remove comments before using this config

	{
		"DB":{                           // Database section
			"Addr":"localhost:6379", // address and port of Redis DB
			"Timeout":30,            // timeout per connection (in seconds)
			"MaxClients":30          // maximum number of clients in DB connection pool
		},
		"Metrics":{                      // Metrics section
			"GetScript":"get.lua",   // relative path to get script
			"AddScript":"add.lua",   // relative path to add script					
			"TtlScript":"ttl.lua",   // relative path to ttl script
			"TTL":86600              // time to live per metric (in seconds)
		},
		"Listener":{                     // Listener server section 
			"Addr":":2003"           // local address and port
		},
		"Webserver":{                    // Webserver section
			"Addr":":8080",          // local address and port
			"Timeout":30             // timeout per connection (in seconds)
		},
		"TTL":{                          // TTL daemon section
			"Tick":12                // timers tick (in seconds)
		}
	}

## TODO
* Consolidate datapoints to improve rendering performance (support for maxDataPoints).
* Add more functions
* Uncaught stacked series cannot have differing numbers of points: 10 vs 251; see Rickshaw.Series.fill() 
[http://beecy.net/post/2009/04/15/fixing-data-series-for-chart-cannot-have-different-number-of-data-points.aspx]


 
## Copyright and licensing

Copyright 2013 *Kuba Podgorski*

Unless otherwise noted, the source files are distributed under the 
*MIT License* found in the LICENSE file.
