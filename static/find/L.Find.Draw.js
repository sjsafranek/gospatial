/**
 * FIND Draw Client
 * Author: Stefan Safranek
 * Email:  sjsafranek@gmail.com
 */

L.Find.Draw = L.Class.extend({
// L.Find.Draw = L.Find.extend({
	
	options: {},

	initialize: function(datasources, options) {
		L.setOptions(this, options || {});
		this.find = L.find(datasources);
		// this.find = L.Find.prototype.initialize(datasources);
		// L.Find.prototype.initialize.call(datasources);
		this._featureTypes = {
			"marker": "Point",
			"polygon": "Polygon",
			"rectangle": "Polygon",
			"polyline": "LineString"
		};
		this.drawnItems = null;
		this._map = null;
	},

	addTo: function(map) {
		this._map = map;
		this.find.addTo(map);
		this._addMeasureControl();
		this._addLocateControl();
		this._addGeosearchControl();
		this._addDrawingControl();
		this._addDrawEventHandlers();
		this._addFeaturePropertiesControl();
	},

	_addMeasureControl: function() {
		var measureControl = new L.Control.Measure();
		measureControl.addTo(this._map);
	},

	_addLocateControl: function() {
		L.control.locate().addTo(this._map);
	},

	_addGeosearchControl: function() {
		new L.Control.GeoSearch({
			provider: new L.GeoSearch.Provider.OpenStreetMap()
		}).addTo(this._map);
	},

	_addDrawingControl: function() {
		this.drawnItems = L.featureGroup().addTo(this._map);
		this._map.addControl(new L.Control.Draw({
			draw: { circle: false },
			edit: { featureGroup: this.drawnItems }
		}));
	},

	_addFeaturePropertiesControl: function() {
		// Ui Control element
		featurePropertiesControl = L.control({position: 'bottomleft'});
		featurePropertiesControl.onAdd = function () {
			var div = L.DomUtil.create('div', 'info legend properties_form');
			div.innerHTML = "<div>";
			div.innerHTML += "<strong>Feature Properties </strong>";
			div.innerHTML += "<a href='#' id='add_property'>[Add Field]</a>";
			div.innerHTML += "</div>";
			div.innerHTML += "<div id='properties'>";
			div.innerHTML += "</div>";
			return div;
		};
		featurePropertiesControl.addTo(this._map);
		this.preventPropogation(featurePropertiesControl);
		// Ui Event Handlers
		$("#add_property").on("click", function() {
			$("#properties").append("<input type='text' class='field' placeholder='field'><input type='text' class='attr' placeholder='attribute'><br>");
		});
	},

	_addDrawEventHandlers: function() {
		function onMapClick(e) {
			var popup = L.popup();
			if (e.target.editing._enabled) { 
				console.log('editing enabled')  
		 	}
			else {
				popup
					.setLatLng(e.latlng)
					.setContent("<div class='button' value='Submit Feature' onClick='FindGeo.sendFeature(" + e.target._leaflet_id + ")'><h4>Submit Feature</h4><div>")
					.openOn(map);
			}
		}
		find_draw = this;
		this._map.on('draw:created', function(event) {
			var layer = event.layer;
			layer.on('click', onMapClick);
			layer.options.color='blue';
			layer.layerType = event.layerType;
			find_draw.drawnItems.addLayer(layer);
		});
		// testing 1234
		this._map.on("draw:drawstop", function(event) {
			// console.log(event);
			var key = Object.keys(find_draw.drawnItems._layers).pop();
			console.log(find_draw.drawnItems._layers[key]);
			// Send through websocket
			find_draw.find.ws.send(find_draw.drawnItems._layers[key])
		});
		this._map.on("draw:editstop", function(event) {
			// console.log(event);
			var key = Object.keys(find_draw.drawnItems._layers).pop();
			console.log(find_draw.drawnItems._layers[key]);
			// Send through websocket
			find_draw.find.ws.send(find_draw.drawnItems._layers[key])
		});
	},

	preventPropogation: function(obj) {
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

	getProperties: function() {
		var properties = {};
		var fields = $("#properties .field");
		var attrs = $("#properties .attr");
		for (var _i=0; _i < fields.length; _i++) {
			properties[fields[_i].value] = attrs[_i].value;
		}
		return properties;
	},

	sendFeature: function(id) {
		var results;
		var feature = this.drawnItems._layers[id];
		var payload = {
			"geometry": {
				"type": this._featureTypes[feature.layerType],
				"coordinates": []
			},
			"properties": this.getProperties()
		}
		if (payload.geometry.type == "Point") {
			payload.geometry.coordinates.push(feature._latlng.lng);
			payload.geometry.coordinates.push(feature._latlng.lat);
		} else if (payload.geometry.type == "LineString") {
			for (var i = 0; i < feature._latlngs.length; i++) {
				payload.geometry.coordinates.push([feature._latlngs[i].lng,feature._latlngs[i].lat])
			}
		} else if (payload.geometry.type == "Polygon") {
			payload.geometry.coordinates.push([]);
			for (var i = 0; i < feature._latlngs.length; i++) {
				payload.geometry.coordinates[0].push([feature._latlngs[i].lng,feature._latlngs[i].lat])
			}
		} else {
			alert("Unknown feature type!")
			return results
		}
		results = this.postRequest(
			'/api/v1/layer/' + $('#layers').val() + '/feature',
			JSON.stringify(payload)
		);
		console.log(this.$super);
		this.find.getLayer($('#layers').val());
		map.removeLayer(this.drawnItems._layers[id]);
		$("#properties .attr").val("");
		return results;
	},

	postRequest: function(route, data) {
		var results;
		find = this;
		$.ajax({
			crossDomain: true,
			type: "POST",
			async: false,
			data: data,
			url: route,
			dataType: 'JSON',
			success: function (data) {
				try {
					results = data;
					// find.find.ws.send("update");
				}
				catch(err){  console.log('Error:', err);  }
			},
			error: function(xhr,errmsg,err) {
				console.log(xhr.status,xhr.responseText,errmsg,err);
				result = null;
				var message = "status: " + xhr.status + "<br>";
				message += "responseText: " + xhr.responseText + "<br>";
				message += "errmsg: " + errmsg + "<br>";
				message += "Error:" + err;
				find.errorMessage(message);
			}
		});
		return results;
	},

	getRequest: function(route, data) {
		var results;
		find = this;
		$.ajax({
			crossDomain: true,
			type: "GET",
			async: false,
			data: data,
			url: route,
			dataType: 'JSON',
			success: function (data) {
				try {
					results = data;
				}
				catch(err){  console.log('Error:', err);  }
			},
			error: function(xhr,errmsg,err) {
				console.log(xhr.status,xhr.responseText,errmsg,err);
				result = null;
				var message = "status: " + xhr.status + "<br>";
				message += "responseText: " + xhr.responseText + "<br>";
				message += "errmsg: " + errmsg + "<br>";
				message += "Error:" + err;
				find.errorMessage(message);
			}
		});
		return results;
	}


});

L.find.draw = function(datasources, options) {
	return new L.Find.Draw(datasources, options);
};

