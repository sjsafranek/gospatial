

	getWebSocket: function() {
		var self = this;
		console.log("Opening websocket");
		try { 
			var url = "ws://" + window.location.host + "/ws/" + self.datasources[0];
			ws = new WebSocket(url);
		}
		catch(err) {
			console.log(err);
			var url = "wss://" + window.location.host + "/ws/" + self.datasources[0];
			ws = new WebSocket(url);
		}
		ws.onopen = function(e) { 
			console.log("Websocket is open");
		};
		ws.onmessage = function(e) {
			console.log(e.data);
			var data = JSON.parse(e.data);
			console.log(data);
			if (data.update) {
				self.apiClient.getLayer($('#layers').val(), function(error, result){
					if (error) {
						throw error;
						self.errorMessage(error);
					} else {apiClient
						self.updateFeatureLayers(result);
					}
				});
			}
			$("#viewers").text(data.viewers);
			if (data.key) {
				if (!self._editFeatures.hasOwnProperty(data.client)) {
					self._editFeatures[data.client] = {
						color: self.utils.randomColor()
					};
				}
				if (self._editFeatures[data.client].hasOwnProperty(data.key)) {
					self._map.removeLayer(self._editFeatures[data.client][data.key]);
				}
				if (data.feature) {
					var featureLayer = L.geoJson(data.feature, {
						style: {
							fillOpacity: 0.5,
							color: self._editFeatures[data.client].color
						}
					});
					// featureLayer.editable.enable();
					self._editFeatures[data.client][data.key] = featureLayer;
					self._editFeatures[data.client][data.key].addTo(self._map);
				}
			}
		};
		ws.onclose = function(e) { 
			console.log("Websocket is closed"); 
		}
		ws.onerror = function(e) { console.log(e); }
		return ws;
	}
