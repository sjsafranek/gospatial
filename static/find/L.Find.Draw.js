/**
 * FIND Draw Client
 * Author: Stefan Safranek
 * Email:  sjsafranek@gmail.com
 */

L.Find.Draw = L.Class.extend({
    
    options: {},

    initialize: function(datasources, options) {
        L.setOptions(this, options || {});
        this.find = L.find(datasources);
        this._map = null;
    },

    addTo: function(map) {
        // map.addLayer(this);
        this.find.addTo(map);
        this._map = map;
    }

});

L.find.draw = function(datasources, options) {
    return new L.Find.Draw(datasources, options);
};

// draw = L.find.draw();
// draw.addTo(map);