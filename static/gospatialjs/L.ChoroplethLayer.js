
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

	_getD3ScaleLinear: function() {
		// check d3 version
		if ("3" == d3.version[0]) {
			return d3.scale.linear()
		} else if ("4" == d3.version[0])  {
			return d3.scaleLinear()
		} else {
			throw new Error("Unsupported d3 version: " + d3.version);
		}
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
						fields[i].color = this._getD3ScaleLinear()
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
