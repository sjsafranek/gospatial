/**
 * FIND Client
 * Author: Stefan Safranek
 * Email:  sjsafranek@gmail.com
 */


// TODO
// Server argument
// Get datsources request on initialize
// utils class
// communication class
//

L.Find = L.Class.extend({

	options: {},

	initialize: function(apikey, datasources, options) {
		L.setOptions(this, options || {});
		this._map = null;
		this.apikey = apikey;
		this.datasources = datasources;
		this.featureLayers = {};
		this.ws = null;
		this._editFeatures = {};
	},

	addTo: function(map) {
		self = this;
		this._map = map;
		this._addUiControls();
		this.getLayer($('#layers').val());
		try {
			this._map.fitBounds(
				self.featureLayers[$('#layers').val()].getBounds()
			);
		}
		catch (err) {
			this._map.fitWorld();
		}
		this.ws = this.getWebSocket();
		return this;
	},

	parseURL: function(url) {
		// http://www.abeautifulsite.net/parsing-urls-in-javascript/
	    var parser = document.createElement('a'),
	        searchObject = {},
	        queries, split, i;
	    // Let the browser do the work
	    parser.href = url;
	    // Convert query string to object
	    queries = parser.search.replace(/^\?/, '').split('&');
	    for( i = 0; i < queries.length; i++ ) {
	        split = queries[i].split('=');
	        searchObject[split[0]] = split[1];
	    }
	    return {
	        protocol: parser.protocol,
	        host: parser.host,
	        hostname: parser.hostname,
	        port: parser.port,
	        pathname: parser.pathname,
	        search: parser.search,
	        searchObject: searchObject,
	        hash: parser.hash
	    };
	},

	_preventPropogation: function(obj) {
		map = this._map;
		// http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
		// Disable dragging when user's cursor enters the element
		obj.getContainer().addEventListener('mouseover', function () {
			map.dragging.disable();
		});
		// Re-enable dragging when user's cursor leaves the element
		obj.getContainer().addEventListener('mouseout', function () {
			map.dragging.enable();
		});
	},

	_addUiControls: function() {
		this._addLogoControl();
		this._addLayerControl();
		this._addMouseControl();
		this._addFeatureAttributesControl(); 
	},

	_addLogoControl: function() {
		var logo = L.control({position : 'topleft'});
		logo.onAdd = function () {
			this._div = L.DomUtil.create('div', 'logo-hypercube');
			this._div.innerHTML = "<img class='img-logo-hypercube' src='/images/HyperCube2.png' alt='logo'>"
			return this._div;
		};
		logo.addTo(this._map);
	},

	_addLayerControl: function() {
		self = this;
		// Create UI control element
		geojsonLayerControl = L.control({position: 'topright'});
		geojsonLayerControl.onAdd = function () {
			var div = L.DomUtil.create('div', 'info legend');
			div.innerHTML = '<div><button id="submitTileLayer">Add TileLayer</button> <input type=text id="newTileLayer"></input></div>';
			div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="geojson" id="layers"></select>';
			div.innerHTML += '<br>Viewers: <span id="viewers">1</span>';
			return div;
		};
		geojsonLayerControl.addTo(this._map);
		//
		$("#submitTileLayer").on("click", function() {
			var newTiles = L.tileLayer(
				$("#newTileLayer").val(),
				{maxZoom:25});
			newTiles.addTo(self._map);
		});
		this._preventPropogation(geojsonLayerControl);
		// Fill drop down options
		for (var _i=0; _i < this.datasources.length; _i++) {
			var obj = document.createElement('option');
			obj.value = this.datasources[_i];
			obj.text = this.datasources[_i];
			$('#layers').append(obj);
		}
		// UI Events listeners
		$('#layers').on('change', function(){ 
			self.getLayer($('#layers').val());
		});
		$('#zoom').on('click', function(){ 
			self._map.fitBounds(
				self.featureLayers[$('#layers').val()].getBounds()
			);
		});
	},

	_addMouseControl: function() {
		// Create UI control element
		mouseLocationControl = L.control({position: 'bottomright'});
		mouseLocationControl.onAdd = function () {
			var div = L.DomUtil.create('div');
			div.innerHTML = "<div id='location'></div>";
			return div;
		};
		mouseLocationControl.addTo(this._map);
		this._preventPropogation(mouseLocationControl);
		// UI Event listeners
		this._map.on('mousemove', function(e) {
			$("#location")[0].innerHTML = "<strong>Lat, Lon : " + e.latlng.lat.toFixed(4) + ", " + e.latlng.lng.toFixed(4) + "</strong>";
		});
	},

	_addFeatureAttributesControl: function() {
		// Create UI control element
		featureAttributesControl = L.control({position: 'bottomright'});
		featureAttributesControl.onAdd = function () {
			var div = L.DomUtil.create('div', 'info legend');
			div.innerHTML = "<div id='attributes'>Hover over features</div>";
			return div;
		};
		featureAttributesControl.addTo(this._map);
		this._preventPropogation(featureAttributesControl);
	},

	getLayer: function(datasource) {
		self = this;
		this.getRequest("/api/v1/layer/" + datasource, function(error, result){
			if (error) {
				throw error;
				self.errorMessage(error);
			} else {
				self.updateFeatureLayers(result);
			}
		});
	},

	updateFeatureLayers: function(data) {
		for (var _i in this.featureLayers){
			if (this._map.hasLayer(this.featureLayers[_i])) {
				this._map.removeLayer(this.featureLayers[_i]);
			}
		}
		try {
			this.featureLayers[$('#layers').val()] = this.createFeatureLayer(data);
			this.featureLayers[$('#layers').val()].addTo(this._map);
		}
		catch(err) { console.log(err); }
	},

	createFeatureLayer: function(data) {
		var featureLayer = L.geoJson(data, {
			style: {
				weight: 2, 
				color: "#000", 
				fillOpacity: 0.25,
			},
			pointToLayer: function(feature, latlng) {
				return L.circleMarker(latlng, {
					radius: 4,
					weight: 1,
					fillOpacity: 0.25,
					color: '#000'
				});
			},
			onEachFeature: function (feature, layer) {

				function highlightFeature(e) {
					var layer = e.target;
					layer.setStyle({
						weight: 3,
						opacity: 1,
						color: '#000'
					});
					if (!L.Browser.ie && !L.Browser.opera) {
						layer.bringToFront();
					}
				}

				function resetHighlight(e) {
					featureLayer.resetStyle(e.target);
				}

				function zoomToFeature(e) {
					map.fitBounds(e.target.getBounds());
				}

				layer.on({
					mouseover: function(feature){
						var properties = feature.target.feature.properties;
						var results = "<table>";
						results += "<th>Field</th><th>Attribute</th>";
						for (var item in properties) {
							results += "<tr><td>" + item + "</td><td>" + properties[item] + "</td></tr>";
						}
						results += "</table>";
						$("#attributes")[0].innerHTML = results;
						highlightFeature(feature);
					},
					mouseout: function(feature){
						$("#attributes")[0].innerHTML = "Hover over features";
						resetHighlight(feature);
					},
					click: function(feature) {
						zoomToFeature(feature);
					}
				});
			}
		});
		return featureLayer;
	},

	_errorMessage: function(message) {
		console.log(message)
		$(".err span").html(message);
		$("#error").show();
		$("#map").hide();
	},

	getRequest: function(route, callback) {
		self = this;
		$.ajax({
			crossDomain: true,
			type: "GET",
			async: false,
			url: route + "?apikey=" + self.apikey,
			dataType: 'JSON',
			success: function (data) {
				return callback(null, data);
			},
			error: function(xhr,errmsg,err) {
				console.log(xhr.status,xhr.responseText,errmsg,err);
				result = null;
				var message = "status: " + xhr.status + "<br>";
				message += "responseText: " + xhr.responseText + "<br>";
				message += "errmsg: " + errmsg + "<br>";
				message += "Error:" + err;
				return callback(new Error(message));
			}
		});
	},

	randomColor: function() {
		return '#'+Math.floor(Math.random()*16777215).toString(16);
	},

	getWebSocket: function() {
		self = this;
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
				self.getLayer($('#layers').val());
			}
			$("#viewers").text(data.viewers);
			if (data.key) {
				if (!self._editFeatures.hasOwnProperty(data.client)) {
					self._editFeatures[data.client] = {
						color: self.randomColor()
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


});

L.find = function(datasources, options) {
	return new L.Find(datasources, options);
};
