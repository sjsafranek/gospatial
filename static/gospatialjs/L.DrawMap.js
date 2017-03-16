
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
