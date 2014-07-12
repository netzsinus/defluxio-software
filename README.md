defluxio
========

lat. defluxio - "die Abweichung"


Web server components
---------------------

This is the web server part. It accepts new measurements and pushes
realtime updates to all attached browsers. It uses the 
[Gorilla websocket](https://github.com/gorilla/websocket) package for
the websocket part, the browser javascript uses JQuery and
[JustGage](http://justgage.com/) as well as
[Epoch](http://fastly.github.io/epoch/)for the eye candy.

## Running 

Once you have Go up and running, you can download, build and run the example
using the following commands.

    $ go get github.com/gorilla/websocket
    $ go get github.com/gorilla/mux
    $ go run *.go

## Installing InfluxDB

You need to install the InfluxDB as outlined [on the influxdb
website](http://influxdb.com/docs/v0.7/introduction/installation.html).
to install the commandline interface.



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


