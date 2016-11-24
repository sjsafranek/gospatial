

	var MapView = Backbone.View.extend({
		el: "#map",

		initialize: function(gospatial) {
			_.bindAll(this, 
				'render',
				'_preventPropogation',
				'_addMouseControl',
				'changeLayer',
				'addProperty',
				'zoomToLayer',
				'_renderTemplate'
			);
			this._map = gospatial._map;
			this.gospatial = gospatial;
			this.render();
			return this;
		},

	    events: {
	    	"change #layers": 'changeLayer',
	    	'click #zoom': 'zoomToLayer',
	    	'click #add_property': 'addProperty'
	    },

	    changeLayer: function() {
			var self = this;
			this.gospatial.apiClient.getLayer($('#layers').val(), function(error, result){
				if (error) {
					swal("Error!", error, "error");
				} else {
					self.gospatial.updateFeatureLayers(result);
				}
			});
	    },

	    addProperty: function() {
			$("#properties").append("<input type='text' class='field' placeholder='field'><input type='text' class='attr' placeholder='attribute'><br>");
	    },

	    zoomToLayer: function() {
			this._map.fitBounds(
				this.gospatial.vectorLayers[$('#layers').val()].getBounds()
			);
	    },

	    _renderTemplate: function() {
	    	// Add logo
			var logo = L.control({position : 'topleft'});
			logo.onAdd = function () {
				this._div = L.DomUtil.create('div', 'logo-hypercube');
				this._div.innerHTML = "<img class='img-logo-hypercube' src='/images/HyperCube2.png' alt='logo'>"
				return this._div;
			};
			logo.addTo(this._map);

			// Measurement control
			var measureControl = new L.Control.Measure();
			measureControl.addTo(this._map);

			// Geo-locate control
			L.control.locate().addTo(this._map);

			// GeoSearch control
			new L.Control.GeoSearch({
				provider: new L.GeoSearch.Provider.OpenStreetMap()
			}).addTo(this._map);

			// Choropleth options
			featureAttributesControl = L.control({position: 'bottomright'});
			featureAttributesControl.onAdd = function () {
				var div = L.DomUtil.create('div', 'info legend');
				div.innerHTML += "<div id='filters'></div>";
				return div;
			};
			featureAttributesControl.addTo(this._map);
			this._preventPropogation(featureAttributesControl);

			// Feature properties
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
			this._preventPropogation(featurePropertiesControl);
	    },

		/** 
		 * method:     _addMouseControl()
		 * desciption: Creates L.control for displaying cursor location
		 */
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

		/** 
		 * method:     _preventPropogation()
		 * source:     http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
		 * desciption: disables mouseover map events from leaflet control objected
		 * @param obj {L.control} Leaflet control object
		 */
		_preventPropogation: function(obj) {
			var self = this;
			// http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
			// Disable dragging when user's cursor enters the element
			obj.getContainer().addEventListener('mouseover', function () {
				self._map.dragging.disable();
				self._map.scrollWheelZoom.disable();
				self._map.doubleClickZoom.disable();
			});
			// Re-enable dragging when user's cursor leaves the element
			obj.getContainer().addEventListener('mouseout', function () {
				self._map.dragging.enable();
				self._map.scrollWheelZoom.enable();
				self._map.doubleClickZoom.enable();
			});
		},

		render: function() {
			var self = this;
			this._renderTemplate();
		}

	});

















/**
 * GoSpatial Map Client
 * Author: Stefan Safranek
 * Email:  sjsafranek@gmail.com
 */

L.GoSpatial = L.Class.extend({

	options: {},

	initialize: function(apikey, options) {
		var self = this;
		L.setOptions(this, options || {});
		this._map = null;
		this.apikey = apikey;
		this.apiClient = new GoSpatialApi(apikey);
		this.color = d3.scale.category10();
		this.vectorLayers = {};
		// this.ws = null;
		this.drawnItems = null;
		this.uuid = this.utils.uuid();
		this._editFeatures = {};
	},

	addTo: function(map) {
		this._map = map;
		// Reading
		this._addLayerControl();
		// Drawing
		this._addDrawingControl();
		this._addDrawEventHandlers();
		// this.ws = this.getWebSocket();
		return this;
	},

	/** 
	 * object:     utils{}
	 * desciption: contains general methods
	 */
	utils: new Utils(),

	/** 
	 * method:     _preventPropogation()
	 * source:     http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
	 * desciption: disables mouseover map events from leaflet control objected
	 * @param obj {L.control} Leaflet control object
	 */
	_preventPropogation: function(obj) {
		map = this._map;
		// http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
		// Disable dragging when user's cursor enters the element
		obj.getContainer().addEventListener('mouseover', function () {
			map.dragging.disable();
			map.scrollWheelZoom.disable();
			map.doubleClickZoom.disable();
		});
		// Re-enable dragging when user's cursor leaves the element
		obj.getContainer().addEventListener('mouseout', function () {
			map.dragging.enable();
			map.scrollWheelZoom.enable();
			map.doubleClickZoom.enable();
		});
	},

	/** 
	 * method:     _addLayerControl()
	 * desciption: Creates L.control for selecting geojson layers
	 */
	_addLayerControl: function() {
		var self = this;
		
		geojsonLayerControl = L.control({position: 'topright'});
		geojsonLayerControl.onAdd = function () {
			var div = L.DomUtil.create('div', 'info legend');
			div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="geojson" id="layers"></select>';
			return div;
		};
		geojsonLayerControl.addTo(this._map);

		this.apiClient.getCustomer(function(error,result){
			if (error) {
				 swal("Error!", error, "error");
				 throw new Error(error);
			} else {
				self.customer = result;
				var datasources = self.customer.datasources;
				for (var _i=0; _i < datasources.length; _i++) {
					var obj = document.createElement('option');
					obj.value = datasources[_i];
					obj.text = datasources[_i];
					$('#layers').append(obj);
				}
				var lyr = $('#layers').val();
				self.apiClient.getLayer(lyr, function(error, result){
					if (error) {
						swal("Error!", error, "error");
					} else {
						self.updateFeatureLayers(result);
						try {
							self._map.fitBounds(self.vectorLayers[lyr].getBounds());
						}
						catch (err) {
							self._map.fitWorld();
						}
					}
				});
			}
		});
	},

	/** 
	 * method:     _addDrawingControl()
	 * desciption: enables L.featureGroup
	 */
	_addDrawingControl: function() {
		this.drawnItems = L.featureGroup().addTo(this._map);
		this._map.addControl(new L.Control.Draw({
			draw: { circle: false },
			edit: { featureGroup: this.drawnItems }
		}));
	},

	/** 
	 * method:     getProperties()
	 * desciption: gets feature properties from .properties_form
	 * @returns {json} json of feature properties
	 */
	getProperties: function() {
		var properties = {};
		var fields = $("#properties .field");
		var attrs = $("#properties .attr");
		for (var _i=0; _i < fields.length; _i++) {
			var word = false;
			for (var _j=0; attrs[_i].value.length > _j; _j++) {
				if ("-.0123456789".indexOf(attrs[_i].value[_j]) == -1) {
					word = true;
					break;
				};
			}
			if (word) {
				properties[fields[_i].value] = attrs[_i].value;
			}
			else {
				properties[fields[_i].value] = parseFloat(attrs[_i].value);
			}
		}
		return properties;
	},

	/** 
	 * method:     _addDrawEventHandlers()
	 * desciption: Adds draw events
	 */
	_addDrawEventHandlers: function() {
		function onMapClick(e) {
			var popup = L.popup();
			if (e.target.editing._enabled) { 
				console.log('editing enabled')  
		 	}
			else {
				popup
					.setLatLng(e.latlng)
					.setContent("<button class='btn btn-sm btn-default' value='Submit Feature' onClick='GoSpatial.sendFeature(" + e.target._leaflet_id + ")'>Submit Feature</button>")
					.openOn(map);
			}
		}
		var self = this;
		this._map.on('draw:created', function(event) {
			var layer = event.layer;
			layer.on('click', onMapClick);
			layer.options.color='blue';
			layer.layerType = event.layerType;
			self.drawnItems.addLayer(layer);
		});
		this._map.on("draw:drawstop", function(event) {
			var key = Object.keys(self.drawnItems._layers).pop();
			var feature = self.drawnItems._layers[key];
			var payload = {
				feature: feature.toGeoJSON(),
				key: key,
				client: self.uuid
			}
		});
		this._map.on("draw:editstop", function(event) {
			var key = Object.keys(self.drawnItems._layers).pop();
			var feature = self.drawnItems._layers[key];
			var payload = {
				feature: feature.toGeoJSON(),
				key: key,
				client: self.uuid
			}
		});
	},

	/** 
	 * method:     updateFeatureLayers()
	 * desciption: updates map vector layers
	 * @param data {geojson}
	 */
	updateFeatureLayers: function(data) {
		for (var _i in this.vectorLayers){
			if (this._map.hasLayer(this.vectorLayers[_i])) {
				this._map.removeLayer(this.vectorLayers[_i]);
			}
		}
		try {
			this.vectorLayers[$('#layers').val()] = this.createFeatureLayer(data);
			this.vectorLayers[$('#layers').val()].addTo(this._map);
			this.generateChoroplethColors();
		}
		catch(err) { console.log(err); }
	},

	/** 
	 * method:     createFeatureLayer()
	 * desciption: creates vector layer from geojson
	 * @param data {geojson}
	 */
	createFeatureLayer: function(data) {
		map = this._map;
		var featureLayer = L.geoJson(data, {
			filter: function(feature, layer) {
				return true;
            },
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
				if (feature.properties) {
						var results = "<table>";
						results += "<th>Field</th><th>Attribute</th>";
						for (var item in feature.properties) {
							results += "<tr><td>" + item + "</td><td>" + feature.properties[item] + "</td></tr>";
						}
						results += "</table>";
					layer.bindPopup(results);
				}
				layer.on({
					mouseover: function(e) {
						// highlight feature
						e.target.setStyle({
							weight: 4,
							radius: 6
						});
						if (!L.Browser.ie && !L.Browser.opera) {
							e.target.bringToFront();
						}
					},
					mouseout: function(e){
						// reset style
						e.target.setStyle({
							weight: 1,
							radius: 4
						});
					},
					dblclick: function(e) {
						// center map on feature
						map.setView(e.target.getBounds().getCenter());
					},
					click: function(e) {
						// center map on feature
						map.setView(e.target.getBounds().getCenter());
					},
					contextmenu: function(e) { 
						map.fitBounds(e.target.getBounds(), {maxZoom:12});
					}
				});
			}
		});
		return featureLayer;
	},

	choroplethColors: {},

	getUniqueFeatureProperties: function(){
		var data = {};
		// Get unique values
		this.vectorLayers[$('#layers').val()].eachLayer(function(layer) {
			for (var i in layer.feature.properties) {
				if (!data.hasOwnProperty(i)) {
					data[i] = [];
				}
				var value = layer.feature.properties[i];
				if (data[i].indexOf(value) == -1 && value != null) {
					data[i].push(value);
				}
			}
		});
		// Sort values
		for (var i in data) {
			if (typeof(data[i] == "number")) {
				data[i].sort(function(a, b){return a-b});
			} else {
				data[i].sort();
			}
		}
		return data;
	},

// COLOR ISSUES
	generateChoroplethColors: function() {
		var self = this;
		$("#filters").html("");
		fields = this.getUniqueFeatureProperties();
		this.choroplethColors = {};
		for (var field in fields) {
			if (!this.choroplethColors.hasOwnProperty(field)) {
				// field section
				var field_selector = $("<div>", {
					title: field
				}).append(
					$("<i>", {name: field}).addClass("fa").addClass("fa-paint-brush").on("click", function(){
						$("i.fa.fa-paint-brush").css("color", "black");
						$(this).css("color", "red");
						self.choropleth($(this).attr("name"));
					}),

					$("<strong>").addClass("toggleFieldSection").text(field).on("click", function() {
						// selector filters
						var vis = $(this).parent().find("table").is(':visible');
						if (vis) {
							$(this).parent().find("table").hide();
						} else {
							$(this).parent().find("table").show();
						}
						// range filters
						var vis = $(this).parent().find(".range-selector").is(':visible');
						if (vis) {
							$(this).parent().find(".range-selector").hide();
						} else {
							$(this).parent().find(".range-selector").show();
						}
					})
				);
				if (typeof(fields[field][0]) == "number") {
					// range filter
					this.choroplethColors[field] = { 
						type: "number",
						color: d3.scale.linear()
							.domain([fields[field][0], fields[field][fields[field].length-1]])
							.range(["yellow", "darkred"])
					};
					var rangeSelector = $("<div>").addClass("range-selector").text(fields[field][0] + " - " + fields[field][fields[field].length-1]);
					// Todo: 
					// 		Display gradient
					rangeSelector.hide();
					field_selector.append(rangeSelector);
				} else {
					// selectors filter
					var table = $("<table>").addClass("table").addClass("table-bordered").append(
						$("<thead>").append(
							$("<tr>").append(
								$("<th>").addClass("color").append(
									$("<i>").addClass("fa").addClass("fa-sort")
								),
								$("<th>").addClass("check").append(
									$("<i>").addClass("fa").addClass("fa-sort")
								),
								$("<th>").append(
									$("<i>").addClass("fa").addClass("fa-sort")
								),
								$("<th>").append(
									$("<i>").addClass("fa").addClass("fa-sort")
								)
							)
						)
					);
					var tbody = $("<tbody>");
					// fill table body with checkbox filters
					this.choroplethColors[field] = { 
						type: "string",
						colors: {}
					};
					for (var i=0; i < fields[field].length; i++) {
						this.choroplethColors[field].colors[fields[field][i]] = this.color(i);
						tbody.append(
							$("<tr>").append(
								$("<td>").addClass("cell-color").append(
									$("<i>").addClass("attr-color").css("background", this.color(i))
								),
								$("<td>").append(
									$("<input>", {
										type:"checkbox", 
										name:fields[field][i]
									}) //.addClass("inline")
								),
								$("<td>").text(fields[field][i]),
								$("<td>").text("0")
							)
						);
					}
					table.append(tbody);
					table.hide();
					field_selector.append(table);
				}
				$("#filters").append(field_selector);
			}
		}
	},

	choropleth: function(field) {
		var self = this;
		this.vectorLayers[$('#layers').val()].eachLayer(function(layer) {
			if (self.choroplethColors[field].type == "number" ) {
				layer.setStyle({ 
					weight: 2, 
					color: self.choroplethColors[field].color( layer.feature.properties[field] ), 
					fillOpacity: 0.8,
					fillColor: self.choroplethColors[field].color( layer.feature.properties[field] )
				});
			} else {
				layer.setStyle({
					weight: 2, 
					color: self.choroplethColors[field].colors[layer.feature.properties[field]],  
					fillOpacity: 0.8,
					fillColor: self.choroplethColors[field].colors[layer.feature.properties[field]]
				});
			}
		});
	},

	/** 
	 * method:     sendFeature()
	 * desciption: send feature layer to GoSpatialApi
	 * @param id {integer} integer of drawn feature layer
	 */
	sendFeature: function(id) {
		var self = this;
		// update websockets
		var payload = {
			feature: false,
			key: id,
			client: this.uuid
		}

		// send new feature
		var results;
		var feature = this.drawnItems._layers[id];
		var payload = feature.toGeoJSON();
		payload.properties = this.getProperties();

		// Send request
		this.apiClient.submitFeature(
			$('#layers').val(),
			JSON.stringify(payload),
			function(error, results) {
				if (error) {
					swal("Error!", error, "error");
				} else {
					swal("Success", "Feature has been submitted.", "success");
					self.apiClient.getLayer($('#layers').val(), function(error, result){
						if (error) {
							swal("Error!", error, "error");
						} else {
							self.updateFeatureLayers(result);
						}
					});
					self._map.removeLayer(self.drawnItems._layers[id]);
					$("#properties .attr").val("");
					return results;
				}
			}
		);
	},


});

L.gospatial = function(apikey, options) {
	return new L.GoSpatial(apikey, options);
};





	// getWebSocket: function() {
	// 	var self = this;
	// 	console.log("Opening websocket");
	// 	try { 
	// 		var url = "ws://" + window.location.host + "/ws/" + self.datasources[0];
	// 		ws = new WebSocket(url);
	// 	}
	// 	catch(err) {
	// 		console.log(err);
	// 		var url = "wss://" + window.location.host + "/ws/" + self.datasources[0];
	// 		ws = new WebSocket(url);
	// 	}
	// 	ws.onopen = function(e) { 
	// 		console.log("Websocket is open");
	// 	};
	// 	ws.onmessage = function(e) {
	// 		console.log(e.data);
	// 		var data = JSON.parse(e.data);
	// 		console.log(data);
	// 		if (data.update) {
	// 			self.apiClient.getLayer($('#layers').val(), function(error, result){
	// 				if (error) {
	// 					throw error;
	// 					self.errorMessage(error);
	// 				} else {
	// 					self.updateFeatureLayers(result);
	// 				}
	// 			});
	// 		}
	// 		$("#viewers").text(data.viewers);
	// 		if (data.key) {
	// 			if (!self._editFeatures.hasOwnProperty(data.client)) {
	// 				self._editFeatures[data.client] = {
	// 					color: self.utils.randomColor()
	// 				};
	// 			}
	// 			if (self._editFeatures[data.client].hasOwnProperty(data.key)) {
	// 				self._map.removeLayer(self._editFeatures[data.client][data.key]);
	// 			}
	// 			if (data.feature) {
	// 				var featureLayer = L.geoJson(data.feature, {
	// 					style: {
	// 						fillOpacity: 0.5,
	// 						color: self._editFeatures[data.client].color
	// 					}
	// 				});
	// 				// featureLayer.editable.enable();
	// 				self._editFeatures[data.client][data.key] = featureLayer;
	// 				self._editFeatures[data.client][data.key].addTo(self._map);
	// 			}
	// 		}
	// 	};
	// 	ws.onclose = function(e) { 
	// 		console.log("Websocket is closed"); 
	// 	}
	// 	ws.onerror = function(e) { console.log(e); }
	// 	return ws;
	// }