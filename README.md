Netzsinus software components
=============================

[![Build Status](https://travis-ci.org/netzsinus/defluxio-software.svg?branch=master)](https://travis-ci.org/netzsinus/defluxio-software)

Please note: This is an alpha release. Kittens might get harmed.

Prerequisites
-------------

Install Go > 1.3 and InfluxDB.

## Installing InfluxDB

You need to install the InfluxDB as outlined [on the influxdb
website](http://influxdb.com/docs/v0.7/introduction/installation.html).
to install the commandline interface.

Afterwards, you can install the software components like this:

	$ go install github.com/netzsinus/defluxio-software

If you run into errors regarding other packages, please make sure you
have *all* of the following version control systems installed: git,
mercurial, bazaar and subversion.


TODO: Clean up everything below.

Web server components
---------------------

This is the web server part. It accepts new measurements and pushes
realtime updates to all attached browsers. It uses the 
[Gorilla websocket](https://github.com/gorilla/websocket) package for
the websocket part, the browser javascript uses JQuery and
[JustGage](http://justgage.com/) as well as
[Epoch](http://fastly.github.io/epoch/)for the eye candy.



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


