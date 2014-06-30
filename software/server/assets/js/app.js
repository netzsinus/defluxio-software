$(function() {

  var conn;
  var msg = $("#msg");
  var log = $("#log");
  function appendLog(msg) {
    
  }

 // function random(name) {
 //   var value = 0,
 //       values = [],
 //       i = 0,
 //       last;
 //   return context.metric(function(start, stop, step, callback) {
 //     start = +start, stop = +stop;
 //     if (isNaN(last)) last = start;
 //     while (last < stop) {
 //       last += step;
 //       value = Number.NaN 
 //       values.push(value);
 //     }
 //     callback(null, values = values.slice((start - stop) / step));
 //   }, name);
 // }
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
          data: netfreqdata,
          width: 600,
          height: 200,
          tickFormats: { 
            time: function(d) { 
              console.log(time);
              console.log(d);
              var date=new Date(d*1000);
              return date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds();
            },
            y: function(v) {
              return v+" mHz";
            }
          },
          axes: ['left', 'bottom', 'right'],
          ticks: { time: 6, right: 8, left: 8 }
        }
      );

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
      $("#timevalue").text(ts.toLocaleTimeString());
      areaChartInstance.push([{time: unixtime, y: (data.Value - 50)*1000}])
    }

    



  } else {
    // TODO: Make this a popup.
    appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
  }
});

