$(function() {



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
				values: [ {time: now, y: 50.0} ]
		}
	];
	var rldata = [
		{
			label: 'Regelleistung', 
				values: [ {time: now, y: 0.0} ]
		}
	];
	var freqChart = $('#freqchart').epoch({
				type: 'time.line',
				label: "Frequenz",
				data: netfreqdata,
				width: 460,
				height: 240,
				tickFormats: {
					bottom: function(d) {
						return date2string(new Date(d*1000));
					},
					right: function(d) {
						return d + "mHz";
					}
				},
			axes: ['left', 'bottom', 'right'],
			windowSize: 100,
			historySize: 20,
			queueSize: 60
		});
	var rlChart = $('#rlchart').epoch({
				type: 'time.line',
				label: "Regelleistung",
				data: rldata,
				width: 460,
				height: 240,
				tickFormats: {
					bottom: function(d) {
						return date2string(new Date(d*1000));
					},
					right: function(d) {
						return d + "MW";
					}
				},
			axes: ['left', 'bottom', 'right'],
			windowSize: 100,
			historySize: 20,
			queueSize: 60
		});
		var freqgauge = new JustGage({
			id: "freqgauge",
				value: "n/a",
				min: 49.90,
				max: 50.1,
				levelColors: [ "#2A4026", "#B6D96C", "#2A4026" ], 
				levelColorsGradient: true,
				title: "Frequenz",
				label: "Hz"
		});
		var rlgauge = new JustGage({
			id: "rlgauge",
				value: "n/a",
				min: -1000,
				max: 1000,
				levelColors: [ "#2A4026", "#B6D96C", "#2A4026" ], 
				levelColorsGradient: true,
				title: "Regelleistung",
				label: "MW"
		}); 

	if (window["WebSocket"]) {
		// get data from the websocket.
		var conn = new WebSocket(ws_endpoint);
		conn.onclose = function(evt) {
			console.log("Connection closed.");
			freqgauge.refresh("n/a");
			rlgauge.refresh("n/a");
		}
		conn.onmessage = function(evt) {
			data = JSON && JSON.parse(evt.data) || $.parseJSON(evt.data);
			freqgauge.refresh(Number(data.Value).toFixed(3));
			var regelleistung=0.0;
			if (data.Value < 50-0.01) {
				regelleistung = 17200*(50 - data.Value);
			}
			if (data.Value > 50+0.01) {
				regelleistung = 17200*(50 - data.Value);
			}
			rlgauge.refresh(Number(regelleistung).toFixed(1));
			var ts = new Date(Date.parse(data.Timestamp));
			var unixtime = Math.round(ts.getTime()/1000);
			$("#timevalue").text(date2string(ts));
			freqChart.push([{time: unixtime, y: (data.Value - 50)*1000}])
			rlChart.push([{time: unixtime, y: regelleistung}])
		}
	} else {
		$("#warnings").html("Ihr Browser unterstützt keine Websockets. Daher können Sie leider keine Frequenzdaten empfangen.");
		$("#warnings").addClass("alert alert-danger");
	}
});

