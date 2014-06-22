$(function() {

  var conn;
  var msg = $("#msg");
  var log = $("#log");
  var seriesData = [ 50.00 ];
  function appendLog(msg) {
    
  }

  if (window["WebSocket"]) {
    // initialize gauge
    var g = new JustGage({
      id: "gauge",
        value: "n/a",
        min: 49.80,
        max: 50.2,
        levelColors: [ "#CC0000", "#008000", "#CC0000" ], 
        levelColorsGradient: true,
        title: "Frequenz"
    }); 
    // initialize graph
    // TODO
    // get data from the websocket.
    conn = new WebSocket(ws_endpoint);
    conn.onclose = function(evt) {
      console.log("Connection closed.");
      g.refresh("n/a");
    }
    conn.onmessage = function(evt) {
      data = JSON && JSON.parse(evt.data) || $.parseJSON(evt.data);
      g.refresh(data.Value);
      ts = new Date(Date.parse(data.Timestamp));
      $("#timevalue").text(ts.toLocaleTimeString());
      // TODO: Update Graph
    }

    



  } else {
    // TODO: Make this a popup.
    appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
  }
});

