# Web server component

This is the web server part. It accepts new measurements and pushes
realtime updates to all attached browsers. It uses the 
[Gorilla websocket](https://github.com/gorilla/websocket) package for
the websocket part, the browser javascript uses JQuery and
[JustGage](http://justgage.com/) for the eye candy.

## Running 

Once you have Go up and running, you can download, build and run the example
using the following commands.

    $ go get github.com/gorilla/websocket
    $ go get github.com/gorilla/mux
    $ go run *.go


