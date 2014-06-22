$(function() {

  var conn;
  var msg = $("#msg");
  var log = $("#log");

  function appendLog(msg) {
    
  }

  if (window["WebSocket"]) {
    var g = new JustGage({
      id: "gauge",
        value: 67,
        min: 49.90,
        max: 50.1,
        title: "Frequenz"
    }); 
    conn = new WebSocket(ws_endpoint);
    conn.onclose = function(evt) {
      appendLog($("<div><b>Connection closed.</b></div>"));
    }
    conn.onmessage = function(evt) {
      appendLog($("<div/>").text(evt.data));
      data = JSON && JSON.parse(evt.data) || $.parseJSON(evt.data);
      console.log("Frequenz: " + data.Value)
      g.refresh(data.Value);
    }
  } else {
    appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
  }
});

