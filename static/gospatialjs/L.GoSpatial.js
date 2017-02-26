
	// Source: Stacked Overflow
	// http://stackoverflow.com/questions/1960473/unique-values-in-an-array
	Array.prototype.getUnique = function(){
	   var u = {}, a = [];
	   for(var i = 0, l = this.length; i < l; ++i){
	      if(u.hasOwnProperty(this[i])) {
	         continue;
	      }
	      a.push(this[i]);
	      u[this[i]] = 1;
	   }
	   return a;
	}

	Array.prototype.getMin = function() {
		var n = null;
		for (var i=0; i<this.length; i++) {
			if ("number" == typeof(this[i])) {
				if (null == n) {
					n = this[i];
				}
				if (n > this[i]) {
					n = this[i];
				}
			}
		}
		return n;
	}

	Array.prototype.getMax = function() {
		var n = null;
		for (var i=0; i<this.length; i++) {
			if ("number" == typeof(this[i])) {
				if (null == n) {
					n = this[i];
				}
				if (n < this[i]) {
					n = this[i];
				}
			}
		}
		return n;
	}



	var MapView = Backbone.View.extend({
		
		el: "#map",

		initialize: function(div, apikey) {
			
			_.bindAll(this, 
				'render',
				'_preventPropogation',
				'_addMouseControl',
				'changeLayer',
				'addProperty',
				'zoomToLayer',
				'_addDrawEventHandlers',
				'getProperties',
				'sendFeature',
				'_renderTemplate'
			);

			this.utils = new Utils();
			this.uuid = this.utils.uuid();

			this.api = new GoSpatialApi(apikey);

			//this._map = L.map(div, {maxZoom: 23});
			this._map = L.drawMap(apikey, div, {maxZoom: 23});

			//this._addLayerControl();

            osm = L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png',{ 
                attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
            });
            osm.addTo(this._map);

			this.featureGroup = new L.choroplethLayer({});
			this.featureGroup.addTo(this._map);

			// this.ws = null;
			this.drawnItems = null;
			this._editFeatures = {};
			this._addDrawEventHandlers();

			this._addLayerControl();

			// Add base layer control
            L.control.layers(
	            {
	                "OpenStreetMap": osm,
	                "Topographic": L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22}),
	                "Streets": L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Street_Map/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22}),
	                "Imagery": L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22}),
	                "Hydda_Full": L.tileLayer('http://{s}.tile.openstreetmap.se/hydda/full/{z}/{x}/{y}.png', {
	                    attribution: 'Tiles courtesy of <a href="http://openstreetmap.se/" target="_blank">OpenStreetMap Sweden</a> &mdash; Map data &copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
	                    reuseTiles: true
	                }),
	                "OpenStreetMap_HOT": L.tileLayer('http://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png', {
	                    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>, Tiles courtesy of <a href="http://hot.openstreetmap.org/" target="_blank">Humanitarian OpenStreetMap Team</a>',
	                    reuseTiles: true
	                }),
	                "Esri_DarkGrey": L.tileLayer('https://services.arcgisonline.com/ArcGIS/rest/services/Canvas/World_Dark_Gray_Base/MapServer/tile/{z}/{y}/{x}', {
	                    attribution: 'Esri, HERE, DeLorme, MapmyIndia, © OpenStreetMap contributors, and the GIS user community',
	                    reuseTiles: true,
	                    maxZoom: 16
	                })
	            }, 
            	{ 
            		Buildings: (function() {
			            return new OSMBuildings(this._map)
			               .date(new Date(2015, 5, 15, 17, 30))
			               .load()
			               .click(function(id) {
			                    console.log('feature id clicked:', id);
			               });
            		})()
            	}, 
            	{
            		position: 'topright'
            	}
            ).addTo(this._map);

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
			this.api.getLayer($('#layers').val(), function(error, result){
				if (error) {
					swal("Error!", error, "error");
				}
				self.featureGroup.setGeoJSON(result);
				try {
					self._map.fitBounds(self.featureGroup.getBounds());
				}
				catch (err) {
					self._map.fitWorld();
				}
			});
	    },

	    addProperty: function() {
			$("#properties").append("<input type='text' class='field' placeholder='field'><input type='text' class='attr' placeholder='attribute'><br>");
	    },

	    zoomToLayer: function() {
			this._map.fitBounds(
				this.featureGroup.getBounds()
			);
	    },

		_addLayerControl: function() {
			var self = this;
			
			geojsonLayerControl = L.control({position: 'topright'});
			geojsonLayerControl.onAdd = function () {
				var div = L.DomUtil.create('div', 'info legend leaflet-bar');
				//div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="geojson" id="layers"></select>';
				div.innerHTML += '<button id="zoom" type="button" class="btn btn-xs btn-default"> <i class="fa fa-search-plus" style="padding-left:5px; margin-right:0px;"></i> </button> <select name="geojson" id="layers"></select>';
				return div;
			};
			geojsonLayerControl.addTo(this._map);

			this.api.getCustomer(function(error,result){
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
					self.changeLayer();
				}
			});
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
				var div = L.DomUtil.create('div', 'panel panel-default leaflet-bar');
				div.id = "filtersControl"
				div.innerHTML += '<div class="panel-heading"><label>Filters</label></div>'
							  +  '<div class="panel-body" id="filters"></div>';
				return div;
			};
			featureAttributesControl.addTo(this._map);
			$( "#filtersControl" ).draggable({
				containment: "#map"
			});
			$( "#filtersControl" ).resizable({
				minHeight: 139,
		  		maxHeight: 560,
				maxWidth: 280
		    });
			this._preventPropogation(featureAttributesControl);

			// Feature properties
			featurePropertiesControl = L.control({position: 'bottomleft'});
			featurePropertiesControl.onAdd = function () {
				var div = L.DomUtil.create('div', 'panel panel-default properties_form leaflet-bar');
				div.innerHTML = "<div class='panel-heading'>"
							  + 	"<strong>Feature Properties </strong>"
							  + 	"<button type='button' class='btn btn-xs btn-default' id='add_property'>Add Field</button>"
				 			  + "</div>";
				div.innerHTML += "<div class='panel-body' id='properties'></div>";
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

		_addDrawEventHandlers: function() {
			var self = this;
			function onMapClick(e) {
				var popup = L.popup();
				if (e.target.editing._enabled) { 
					console.log('editing enabled')  
			 	}
				else {
					popup
						.setLatLng(e.latlng)
						.setContent("<button class='btn btn-sm btn-default' value='Submit Feature' onClick='App.sendFeature(" + e.target._leaflet_id + ")'>Submit Feature</button>")
						.openOn(self._map);
				}
			}
			this._map.on('draw:created', function(event) {
				var layer = event.layer;
				layer.on('click', onMapClick);
				layer.options.color='blue';
				layer.layerType = event.layerType;
				self._map.drawnItems.addLayer(layer);
			});
			this._map.on("draw:drawstop", function(event) {
				var key = Object.keys(self._map.drawnItems._layers).pop();
				var feature = self._map.drawnItems._layers[key];
				var payload = {
					feature: feature.toGeoJSON(),
					key: key,
					client: self.uuid
				}
			});
			this._map.on("draw:editstop", function(event) {
				var key = Object.keys(self.drawnItems._layers).pop();
				var feature = self._map.drawnItems._layers[key];
				var payload = {
					feature: feature.toGeoJSON(),
					key: key,
					client: self.uuid
				}
			});
		},

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
			var feature = this._map.drawnItems._layers[id];
			var payload = feature.toGeoJSON();
			payload.properties = this.getProperties();

			console.log( new Date().toISOString(), "[DEBUG]:", JSON.stringify(payload) );

			swal({
				title: "Create layer",
				text: "Are you sure you want submit feature?",
				type: "info",
				showCancelButton: true,
				confirmButtonColor: "#DD6B55",
				confirmButtonText: "Yes, pls!",
				cancelButtonText: "No, cancel pls!",
				showLoaderOnConfirm: true,
			}).then(function(){
				self.api.submitFeature(
					$('#layers').val(),
					JSON.stringify(payload),
					function(error, results) {
						if (error) {
							swal("Error!", error, "error");
						} else {
							//swal("Success", "Feature has been submitted.", "success");
							swal("Success", JSON.stringify(results), "success");
							self._map.removeLayer(self._map.drawnItems._layers[id]);
							$("#properties .attr").val("");
							self.changeLayer();
						}
					}
				);
			});

		},

		render: function() {
			var self = this;
			this._renderTemplate();
		}

	});



/**
 * Choropleth GeoJSON Layer
 */
L.ChoroplethLayer = L.GeoJSON.extend({

	options: {
		
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
					e.target._map.setView(e.target.getBounds().getCenter());
				},
				
				click: function(e) {
					// center map on feature
					e.target._map.setView(e.target.getBounds().getCenter());
				},
				
				contextmenu: function(e) { 
					e.target._map.fitBounds(e.target.getBounds(), {maxZoom:12});
				}

			});
		}
	},

	initialize: function(geojson, options) {
		this._layers = {};
		this.datasource_id = "";
	},

	addTo: function(map) {
		this._map = map;
		return this;
	},

	// @method 			setGeoJSON
	// @description 	clears layers and creates new layers from supplied geojson
	// @params 			geojson{object} 
	setGeoJSON: function(geojson) {
	    this.clearLayers();
	    this.addData(geojson);
	    this._buildColorTree();
	},

	// @method 			_buildColorTree
	// @description 	Scans layer features and build meta data for datasource columns
	_buildColorTree: function() {
		var fields = {};
		
		// loop through features features
		this.eachLayer(function(layer) { 
			var properties = layer.feature.properties;
			for (var j in properties) {
				if (!fields.hasOwnProperty(j)) {
					fields[j] = {
						attrs:[], 
						type:"", 
						color:null, 
						name: j
					};
				}
				// get field values from feature properties
				fields[j].attrs.push(properties[j]);
			}
		});

		// get unique field values
		for (var i in fields) {
			fields[i].attrs = fields[i].attrs.getUnique();
			for (var j in fields[i].attrs) {
				var item = fields[i].attrs[j];
				if ("string" == typeof(item)) {
					fields[i].type = "string";
					break;
				} else if ("boolean" == typeof(item)) {
					if ( -1 != ["", "boolean"].indexOf(fields[i].type) ) {
						fields[i].type = "boolean";
					} else {
						fields[i].type = "string";
						break;
					}
				} else if ("number" == typeof(item)) {
					if ( -1 != ["","number"].indexOf(fields[i].type) ) {
						fields[i].type = "number";
					} else {
						fields[i].type = "string";
						break;	
					}
				}
			}

			// sort field values
			if ("number" == fields[i].type) {
				fields[i].attrs.sort(function(a, b){return a-b});
			} else {
				fields[i].attrs.sort();
			}

			// create color based on data type
			switch(fields[i].type) {
				case "number":
					// color range
					if (0 != fields[i].attrs.length) {
						fields[i].color = d3.scaleLinear()
											.domain([
												fields[i].attrs.getMin(),
												fields[i].attrs.getMax()
											])
											.range(["#330F53", "#FFDC00"]);
					}
					break;
				case "boolean":
					fields[i].color = d3.schemeCategory10;
					break;
				case "string":
					fields[i].color = d3.schemeCategory20;
					break;
				default:
					console.log("[DEBUG]: Uncaught data type", fields[i]);
			}

		}

		this.columns = fields;

	},

	// @method 			getColumnNames
	// @description 	Returns array of column names
	// @returns 		{array}
	getColumnNames: function() {
		return Object.keys(this.columns);
	},

	// @method 			hasColumn
	// @description 	Checks if datasource has column
	// @params 			column_name{string} 
	// @returns 		{boolean}
	hasColumn: function(column_name) {
		return this.columns.hasOwnProperty(column_name);
	},

	// @method 			getColumn
	// @description 	Returns datasource column object 
	// @params 			column_name{string} 	
	// @returns 		{object}
	getColumn: function(column_name) {
		if (!this.hasColumn(column_name)) {
			throw new Error("Column not found");
		}
		return this.columns[column_name];
	},

	//
	updateColorLegend: function() {

	},

	// @method 			choropleth
	// @description 	colors layer features by selected column
	// @params 			column_name{string} 
	choropleth: function(column_name) {
		var self = this;
		var column = this.getColumn(column_name);
		this.eachLayer(function(layer) {
			var feature = layer.feature; 
			if (column.type == "number" ) {
				layer.setStyle({ 
					weight: 2, 
					color: column.color(feature.properties[column.name]), 
					fillOpacity: 0.8,
					fillColor: column.color(feature.properties[column.name])
				});
			} else {
				var index = column.attrs.indexOf(feature.properties[column.name]);
				layer.setStyle({
					weight: 2, 
					color: column.color(index),
					fillOpacity: 0.8,
					fillColor: column.color(index)
				});
			}
		});
		this.updateColorLegend();
	}

});

L.choroplethLayer = function(geojson, options) {
	return new L.ChoroplethLayer(geojson, options);
};





/**
 * GoSpatial Map Client
 * Author: Stefan Safranek
 * Email:  sjsafranek@gmail.com
 */
L.DrawMap = L.Map.extend({

	initialize: function(apikey, id, options) {
		var self = this;
		L.Map.prototype.initialize.call(this, id, options);
		this.drawnItems = {};
		this._addDrawingControl();
	},

	_addDrawingControl: function() {
		this.drawnItems = L.featureGroup().addTo(this);
		this.addControl(
			new L.Control.Draw({
				draw: { circle: false },
				edit: { featureGroup: this.drawnItems }
			})
		);
	}

});

L.drawMap = function(apikey, container, options) {
	return new L.DrawMap(apikey, container, options);
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
	// 				} else {apiClient
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
