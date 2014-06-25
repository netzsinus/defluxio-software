$(function() {

  var conn;
  var msg = $("#msg");
  var log = $("#log");
  var seriesData = [ 50.00 ];
  function appendLog(msg) {
    
  }

  function random(name) {
    var value = 0,
        values = [],
        i = 0,
        last;
    return context.metric(function(start, stop, step, callback) {
      start = +start, stop = +stop;
      if (isNaN(last)) last = start;
      while (last < stop) {
        last += step;
        value = Math.max(-10, Math.min(10, value + .8 * Math.random() - .4 + .2 * Math.cos(i += .2)));
        values.push(value);
      }
      callback(null, values = values.slice((start - stop) / step));
    }, name);
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
    var context = cubism.context()
        .serverDelay(0)
        .clientDelay(0)
        .step(2e3)
        .size(960);

    // TODO: Replace with ringbuffer thingie
    var foo = random("foo"),
        bar = random("bar");
    // The graph updates itself asynchronously
    d3.select("#chart").call(function(div) {
      div.append("div")
          .attr("class", "axis")
          .call(context.axis().orient("top"));

      div.selectAll(".horizon")
          .data([foo, bar])
        .enter().append("div")
          .attr("class", "horizon")
          .call(context.horizon().extent([-20, 20]));

      div.append("div")
          .attr("class", "rule")
          .call(context.rule());
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
      ts = new Date(Date.parse(data.Timestamp));
      $("#timevalue").text(ts.toLocaleTimeString());
      // TODO: Update ringbuffer
    }

    



  } else {
    // TODO: Make this a popup.
    appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
  }
});

