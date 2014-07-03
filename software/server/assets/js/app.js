$(function() {

  var conn;
  var msg = $("#msg");
  var log = $("#log");
  function appendLog(msg) {
    
  }

  function date2string(date) {
    var hr = date.getHours();
    var min = date.getMinutes();
    if (min < 10) {
      min = "0" + min;
    }
    var sec = date.getSeconds();
    if (sec < 10) {
      sec = "0" + sec;
    }
    return hr + ":" + min + ":" + sec;
  }
  var now = Math.round(new Date()/1000)
  var netfreqdata = [
    {
      label: 'Netzfrequenz', 
            values: [ {time: now, y: 0} ] 
    }
  ];
  var areaChartInstance = $('#freqchart').epoch(
        { 
          type: 'time.line', 
					label: "Frequenz",
          data: netfreqdata,
          width: 460,
          height: 240,
          tickFormats: { 
            bottom: function(d) { 
              return date2string(new Date(d*1000));
            },
            y: function(d) {
              return "bar";
            }
          },
          axes: ['left', 'bottom', 'right'],
          windowSize: 100,
          historySize: 200
        }
      );

  if (window["WebSocket"]) {
    // initialize gauge
    var g = new JustGage({
      id: "freqgauge",
        value: "n/a",
        min: 49.90,
        max: 50.1,
        levelColors: [ "#2A4026", "#B6D96C", "#2A4026" ], 
        levelColorsGradient: true,
        title: "Frequenz",
        label: "Hz"
    }); 
    // get data from the websocket.
    conn = new WebSocket(ws_endpoint);
    conn.onclose = function(evt) {
      console.log("Connection closed.");
      g.refresh("n/a");
    }
    conn.onmessage = function(evt) {
      data = JSON && JSON.parse(evt.data) || $.parseJSON(evt.data);
      g.refresh(data.Value);
      var ts = new Date(Date.parse(data.Timestamp));
      var unixtime = Math.round(ts.getTime()/1000);
      $("#timevalue").text(date2string(ts));
      areaChartInstance.push([{time: unixtime, y: (data.Value - 50)*1000}])
    }
  } else {
    // TODO: Make this a popup.
    appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
  }
});

