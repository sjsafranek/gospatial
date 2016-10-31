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

		// Get customer details
		this.apiClient.getCustomer(function(error,result){
			if (error) {
				 swal("Error!", error, "error");
			} else {
				self.customer = result;
				self.datasources = self.customer.datasources;
				for (var _i=0; _i < self.datasources.length; _i++) {
					var obj = document.createElement('option');
					obj.value = self.datasources[_i];
					obj.text = self.datasources[_i];
					$('#layers').append(obj);
				}
			}
		});

		this.color = d3.scale.category10();
		// this.color = d3.scale.category20b();
		this.vectorLayers = {};
		// this.ws = null;
		this.drawnItems = null;
		this.uuid = this.utils.uuid();
		this._editFeatures = {};
	},

	addTo: function(map) {
		this._map = map;
		// Reading
		this._addLogoControl();
		this._addLayerControl();
		this._addMouseControl();
		// Drawing
		this._addMeasureControl();
		this._addLocateControl();
		this._addGeosearchControl();
		this._addDrawingControl();
		this._addDrawEventHandlers();
		this._addFeaturePropertiesControl();
		this._addChoroplethOptions();
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
	 * method:     _addLogoControl()
	 * desciption: Creates L.control containing logo
	 */
	_addLogoControl: function() {
		var logo = L.control({position : 'topleft'});
		logo.onAdd = function () {
			this._div = L.DomUtil.create('div', 'logo-hypercube');
			this._div.innerHTML = "<img class='img-logo-hypercube' src='/images/HyperCube2.png' alt='logo'>"
			return this._div;
		};
		logo.addTo(this._map);
	},

	/** 
	 * method:     _addLayerControl()
	 * desciption: Creates L.control for selecting geojson layers
	 */
	_addLayerControl: function() {
		var self = this;
		// Create UI control element
		geojsonLayerControl = L.control({position: 'topright'});
		geojsonLayerControl.onAdd = function () {
			var div = L.DomUtil.create('div', 'info legend');
			// div.innerHTML = '<div><button id="submitTileLayer">Add TileLayer</button> <input type=text id="newTileLayer"></input></div>';
			div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="geojson" id="layers"></select>';
			// div.innerHTML += '<br>Viewers: <span id="viewers">1</span>';
			return div;
		};
		geojsonLayerControl.addTo(this._map);
		// this._preventPropogation(geojsonLayerControl);
		var attempt = 0;
		function loadFirstLayer() {
			console.log("loading...");
			var lyr = $('#layers').val();
			if (lyr) {
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
			} else {
				if (attempt < 50) {
					setTimeout(loadFirstLayer, 10);
					attempt++;
				} else {
					swal("Error!", "Failed to load datasource layer: " + lyr, "error");
				}
			}
		}
		setTimeout(loadFirstLayer, 10);

		$('#layers').on('change', function(){ 
			self.apiClient.getLayer($('#layers').val(), function(error, result){
				if (error) {
					swal("Error!", error, "error");
				} else {
					self.updateFeatureLayers(result);
				}
			});
		});

		$('#zoom').on('click', function(){ 
			self._map.fitBounds(
				self.vectorLayers[$('#layers').val()].getBounds()
			);
		});
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
	 * method:     _addMeasureControl()
	 * desciption: enables L.Control.Measure
	 */
	_addMeasureControl: function() {
		var measureControl = new L.Control.Measure();
		measureControl.addTo(this._map);
	},

	/** 
	 * method:     _addLocateControl()
	 * desciption: enables L.control.locate
	 */
	_addLocateControl: function() {
		L.control.locate().addTo(this._map);
	},

	/** 
	 * method:     _addGeosearchControl()
	 * desciption: enables L.Control.GeoSearch
	 */
	_addGeosearchControl: function() {
		new L.Control.GeoSearch({
			provider: new L.GeoSearch.Provider.OpenStreetMap()
		}).addTo(this._map);
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
	 * method:     _addFeaturePropertiesControl()
	 * desciption: creates L.control for adding properties to new features
	 */
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
		this._preventPropogation(featurePropertiesControl);
		// Ui Event Handlers
		$("#add_property").on("click", function() {
			$("#properties").append("<input type='text' class='field' placeholder='field'><input type='text' class='attr' placeholder='attribute'><br>");
		});
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
					// .setContent("<div class='button' value='Submit Feature' onClick='GoSpatial.sendFeature(" + e.target._leaflet_id + ")'><h4>Submit Feature</h4><div>")
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
			// self.ws.send(JSON.stringify(payload));
		});
		this._map.on("draw:editstop", function(event) {
			var key = Object.keys(self.drawnItems._layers).pop();
			var feature = self.drawnItems._layers[key];
			var payload = {
				feature: feature.toGeoJSON(),
				key: key,
				client: self.uuid
			}
			// self.ws.send(JSON.stringify(payload));
		});
	},

/*************************************************************************
 * FEATURE LAYERS
 *************************************************************************/
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
                // return feature.properties.BusType == "Cafe";
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

/*************************************************************************
 * CHOROPLETH
 *************************************************************************/
 	/** 
	 * method:     _addChoroplethOptions()
	 * desciption: Creates L.control for changing choropleth
	 */
	_addChoroplethOptions: function() {
		var self = this;
		// Create UI control element
		featureAttributesControl = L.control({position: 'bottomright'});
		featureAttributesControl.onAdd = function () {
			var div = L.DomUtil.create('div', 'info legend');
			div.innerHTML += "<div id='filters'></div>";
			return div;
		};
		featureAttributesControl.addTo(this._map);
		this._preventPropogation(featureAttributesControl);
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
			console.log(data[i]);
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

/*************************************************************************
 * SUBMIT FEATURES
 *************************************************************************/
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
		// this.ws.send(JSON.stringify(payload));
		// send new feature
		var results;
		var feature = this.drawnItems._layers[id];
		var payload = feature.toGeoJSON();
		payload.properties = this.getProperties();
	/*	
		// add date_created & date_modified to feature properties
		var now = new Date();
		if (!payload.properties.hasOwnProperty("date_created")) {
			// payload.properties.date_created = now.toISOString();
			payload.properties.date_created = parseInt(now.getTime()/1000);
		}
		if (!payload.properties.hasOwnProperty("date_modified")) {
			// payload.properties.date_modified = now.toISOString();
			payload.properties.date_modified = parseInt(now.getTime()/1000);
		}
	*/
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

});

L.gospatial = function(apikey, options) {
	return new L.GoSpatial(apikey, options);
};



