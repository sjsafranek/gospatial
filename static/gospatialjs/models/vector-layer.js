
	var VectorLayer = Backbone.Model.extend({
		initialize: function(){
			console.log("VectorLayer created");
		}
	});

	var VectorLayerCollection = Backbone.Collection.extend({
		model: VectorLayer
	});
