Netzsinus software components
=============================

[![Build Status](https://travis-ci.org/netzsinus/defluxio-software.svg?branch=master)](https://travis-ci.org/netzsinus/defluxio-software)

Please note: This is an alpha release. Kittens might get harmed.

Prerequisites
-------------

Install Go > 1.3 and, if you want to run the server, InfluxDB.

## Installing InfluxDB (only needed for the server)

You need to install the InfluxDB version 0.7 as outlined [on the influxdb
website](http://influxdb.com/docs/v0.7/introduction/installation.html).
to install the commandline interface.

## Installing the defluxio software

If you're new to golang: an important concept is the "workspace" which contains all the libraries, 
project files etc. On a plain machine, I usually do the following:

	$ mkdir -p ~/go/{src|bin|pkg}

Now, add the following two lines to your .profile:

	export GOPATH="$HOME/go"
	export PATH="$GOPATH/bin:$PATH"
	
Reload the profile:

	$ . ~/.profile

Your workspace is now set up. The dependencies of the software are managed using [godep](https://github.com/tools/godep). The installation is a good test whether your environment is set up correctly. Simply run

	go get github.com/tools/godep

You should now have a working godep binary in your path. Continue by checking out the defluxio source:
	
	$ cd ~/go/src
	$ mkdir github.com/netzinus
	$ cd github.com/netzsinus
	$ git clone https://github.com/netzsinus/defluxio-software.git
	$ cd defluxio-software
	
Now, restore the libraries used by defluxio:

	$ godep restore
	
You're all set. You can build the project using 

	$ go install ./...
	
This puts some binaries in the ```~/go/bin``` subdirectory:

	$ ls ~/go/bin
	defluxio-exporter  defluxio-provider  defluxiod  godep

If you run into errors regarding other packages, please make sure you
have *all* of the following version control systems installed: git,
mercurial, bazaar and subversion.

## Configuration 

All binaries can generate a default configuration file which serves as a template for your configuration. In order to set up a new frequency sensor, do the following:

	$ cd ~
	$ mkdir etc
	$ cd etc
	$ defluxio-provider -genconfig=true

Now, edit the file ```defluxio-provider.conf```. If you want to use the "official" netzsin.us server (that would make me happy!) please contact ```md@gonium.net``` for access credentials. On the typical raspberry pi installation, you can start the frequency provider daemon like this:

	$ sudo /home/pi/go/bin/defluxio-provider -config=/home/pi/etc/defluxio-provider.conf

A typical startup sequence looks like this:

````
2015/02/20 11:55:17 Received unknown data: 
2015/02/20 11:55:17 Startup: Ignoring measurement 49.99652099609375
2015/02/20 11:55:17 Received unknown data: 
2015/02/20 11:55:17 Info message: I;Frequency: 49.995566 Hz, delta:  -4 mHz
2015/02/20 11:55:17 Received unknown data: 
2015/02/20 11:55:17 Startup: Ignoring measurement 49.99557113647461
2015/02/20 11:55:17 Received unknown data: 
2015/02/20 11:55:17 Info message: I;Frequency: 49.995868 Hz, delta:  -4 mHz
2015/02/20 11:55:17 Received unknown data: 
2015/02/20 11:55:17 Startup: Ignoring measurement 49.99586868286133
2015/02/20 11:55:17 Received unknown data: 
2015/02/20 11:55:18 Info message: I;Frequency: 49.996924 Hz, delta:  -3 mHz
2015/02/20 11:55:18 Startup: Ignoring measurement 49.99692153930664
2015/02/20 11:55:19 Info message: I;Freque995667 Hz, delta:  -4 mHz
2015/02/20 11:55:19 Startup: Ignoring measurement 49.995670318603516
2015/02/20 11:55:20 Info message: I;Freque996321 Hz, delta:  -4 mHz
2015/02/20 11:55:20 Frequency: 49.99632
````

Notes
-----

* [Cooper Hewitt Font](http://www.cooperhewitt.org/colophon/cooper-hewitt-the-typeface-by-chester-jenkins/)
* [cubism.js data source](https://stackoverflow.com/questions/18069409/are-there-any-tutorials-or-examples-for-cubism-js-websocket)
* [Epoch - Graph library by Fastly](http://fastly.github.io/epoch/)
* [Schmitt-Trigger](http://www.mikrocontroller.net/articles/Schmitt-Trigger)
* [PowerBox: The Safe AC Power Meter](https://instruct1.cit.cornell.edu/Courses/ee476/FinalProjects/s2008/cj72_xg37/cj72_xg37/)
* [Komparatorschaltung mit dem LM393](http://www.ne555.at/schaltungstechnik/390-komparator-mit-lm393-und-einfacher-spannungsversorgung.html)
* [Jaschinsky, Markus: Untersuchung des Zusammenhangs zwischen gemessener Netzfrequenz und Regelenergieeinsatz als Basis eines Reglerentwurfs zum Intraday Lastmanagement](http://edoc.sub.uni-hamburg.de/haw/frontdoor.php?source_opus=2067&la=de)
* [High speed capture mit ATMega Timer](http://www.mikrocontroller.net/articles/High-Speed_capture_mit_ATmega_Timer) - für eine exakte Frequenzmessung


