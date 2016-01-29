Netzsinus software components
=============================

[![Build Status](https://travis-ci.org/netzsinus/defluxio-software.svg?branch=master)](https://travis-ci.org/netzsinus/defluxio-software)


This repository contains the software components (excluding firmware) of
the netzsin.us project. These are the following:

 * ``defluxiod``, the server component (serves HTML and deals with data
	 sent by the sensors). Incoming sensor data will be stored in an
 [InfluxDB](https://influxdata.com) database, although this might change
 in the future.
 * ``defluxio-provider``: The reference implementation of the sensor
   software. It reads the measurements from the sensor and submits them
	 to the server.
 * ``defluxio-exporter``: A small software component that extracts
   measurements from the database and exports them into text files for
	 [data.netzsin.us](data.netzsin.us).

### Installing the software

You need a working [Golang installation](http://golang.org) and the [GB
build tool](http://getgb.io/) in order to compile your binary. Please
install the Go compiler first. Afterwards you can install GB like this:

    go get github.com/constabulary/gb/...

Clone this repository:

    git clone https://github.com/netzsinus/defluxio-software.git

and build it:

    cd defluxio-software
    gb build all

Now, there should be several binaries in the ````bin```` subfolder.

#### Crosscompiling e.g. for Raspberry Pi

Go has very good crosscompilation support. Typically, I develop under
Mac OS and crosscompile a binary for my RPi. It is easy:

    # clear whatever old binaries I have
    rm -rf pkg bin
    # start crosscompilation
    GOOS=linux GOARCH=arm GOARM=5 gb build all

You can then copy the binary from the ``bin`` subdirectory to the RPi
and start it. Please note that ``defluxiod`` will not be built because
it depends on libraries only available for the 64bit x86 architecture.

## Configuration 

All binaries can generate a default configuration file which serves as a template for your configuration. In order to set up a new frequency sensor, do the following:

	$ cd ~
	$ mkdir etc
	$ cd etc
	$ defluxio-provider -genconfig=true

Now, edit the file ``defluxio-provider.conf``. If you want to use the "official" netzsin.us server (that would make me happy!) please contact ``md@gonium.net`` for access credentials. On the typical raspberry pi installation, you can start the frequency provider daemon like this:

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

## Development Setup

If you want to contribute (or test your new server design) you can run
your own server. The server does not store any data by default, see
"Installing InfluxDB" below. But incoming frequency measurements will be
displayed on the webpage. In short, while in the directory
``defluxio-software``, try the following:

1. Run an instance of ``defluxiod`` on your local machine:

````
    $ ./bin/defluxiod -genconfig
		$ ./bin/defluxiod -config=defluxiod.conf
		2016/01/29 10:13:29 Configured meters are:
		2016/01/29 10:13:29 * meter1 (Meter 1)
		2016/01/29 10:13:29 * meter2 (Meter 2)
		2016/01/29 10:13:29 Starting meter surveilance routine
		2016/01/29 10:13:29 Starting server at 127.0.0.1:8080
````

	If you browse to [http://127.0.0.1:8080](http://127.0.0.1:8080), you
	should see a clone of the netzsin.us page.

2. Start a simulated frequency sensor:
````
    $ ./bin/defluxio-provider -genconfig
		$ ./bin/defluxio-provider -config=defluxio-provider.conf -sim
		2016/01/29 10:17:27 Frequency: 49.98730
		2016/01/29 10:17:27 Error posting data:  Post http://127.0.0.1:8080/api/submit/meter1: EOF
		2016/01/29 10:17:29 Frequency: 49.95660
		2016/01/29 10:17:31 Frequency: 50.01041
		2016/01/29 10:17:33 Frequency: 49.97088
		2016/01/29 10:17:35 Frequency: 49.95438
````
	The parameter ``-sim`` does enable the simulation mode - no frequency
	sensor hardware is needed. It just sends random frequency measurements
	to the server.

If you want to submit values to the server using your own client, its
rather easy: You can submit a JSON array using the POST verb. Using
curl:

    $ curl -i -X POST -H "Content-Type: application/json" \
			 -H "X-API-KEY: secretkey1" \
			 -d "{\"Timestamp\": \"`date --rfc-3339=ns | sed 's/ /T/; s/\(\....\).*-/\1-/g'`\", \"Value\":49.9999}" \
			http://127.0.0.1:8080/api/submit/meter1
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 29 Jan 2016 10:42:43 GMT
		Content-Length: 0

In plain text the body of the request looks like this:

    {"Timestamp":"2016-01-29T11:33:22.954022564+01:00","Value":49.9999}

The linux date command does not format the date correctly according to
ISO8601, so a little ``sed`` magic is applied. Please note: The server
currently has a bug leading to failing goroutines on first submission.
Subsequent calls will succeed, just call curl several times.

#### Installing InfluxDB (only needed for the server)

If you want to store frequency measurements you need to install the
InfluxDB version 0.9 as outlined [on the influxdb
website](http://influxdb.com/docs/v0.9/introduction/installation.html).

To enable the database you need to set ``enabled: true`` in your
``defluxiod.conf`` file. You can also use the ``defluxio-exporter``
commandline tool to export the measurements from the database again.

