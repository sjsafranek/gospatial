
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

	    getDatasource: function() {
	    	return $('#layers').val();
	    },

	    changeLayer: function() {
			var self = this;
			this.api.getLayer($('#layers').val(), function(error, result){
				if (error) {
					new SwalPpError("ApiError", error);
					return;
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
					new SwaPplError("ApiError", error);
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


			new SwalConfirm( 
				"Create feature?", 
				"Are you sure you want submit feature?", 
				"info",
				function(){
					self.api.submitFeature(
						$('#layers').val(),
						payload,
						function(error, results) {
							if (error) {
								new SwalPpError("ApiError", error);
								return;
							}
							new SwalPpSuccess("Success", results);
							self._map.removeLayer(self._map.drawnItems._layers[id]);
							$("#properties .attr").val("");
							self.changeLayer();
						}
					);
				}
			);


			// swal({
			// 	title: "Create layer",
			// 	text: "Are you sure you want submit feature?",
			// 	type: "info",
			// 	showCancelButton: true,
			// 	confirmButtonColor: "#337ab7",
			// 	confirmButtonText: "Yes, pls!",
			// 	cancelButtonText: "No, cancel pls!",
			// 	showLoaderOnConfirm: true,
			// }).then(function(){
			// 	self.api.submitFeature(
			// 		$('#layers').val(),
			// 		payload,
			// 		function(error, results) {
			// 			if (error) {
			// 				new SwalPpError("ApiError", error);
			// 				return;
			// 			}
			// 			new SwalPpSuccess("Success", results);
			// 			self._map.removeLayer(self._map.drawnItems._layers[id]);
			// 			$("#properties .attr").val("");
			// 			self.changeLayer();
			// 		}
			// 	);
			// });

		},

		render: function() {
			var self = this;
			this._renderTemplate();
		}

	});