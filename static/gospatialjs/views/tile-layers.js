
	var TileLayerView = Backbone.View.extend({
	    el: "#tileLayerView",

	    initialize: function(apikey){
	        // var lp = this.$(".rightpanel");
	        _.bindAll(this, 'render', 'createTileLayer');
	        this.apikey = apikey;
	        this.render();
	    },

	    events: {
	        "click button.createTileLayer": "createTileLayer",
	    },

		createTileLayer: function() {
			var self = this;
			var datasource_id = $(this).attr("ds_id");
			swal({
				title: "Create tile layer",
				text: "Are you sure you want to create a new tile lyaer",
				type: "info",
				showCancelButton: true,
				confirmButtonColor: "#DD6B55",
				confirmButtonText: "Yes, pls!",
				cancelButtonText: "No, cancel pls!",
				closeOnConfirm: false,
				closeOnCancel: true,
				showLoaderOnConfirm: true
			},
			function(isConfirm){
				if (isConfirm) {
					$.ajax({
						url: '/api/v1/tilelayer?apikey=' + self.apikey,
						data: {
							tilelayer_name: $("#tilelayer_name").val(),
							tilelayer_url: $("#tilelayer_url").val()
						},
						type: 'POST',
						success: function(result) {
							swal("Created!", result, "success");
							// $("#tilelayers_list").html("");
							// refreshTileLayers();
							self.render();
						},
						failure: function(result) {
							console.log(result);
							swal("Error", JSON.stringify(result), "error");
							throw new Error(result);
						},
						error: function(result) {
							console.log(result);
							swal("Error", JSON.stringify(result), "error");
							throw new Error(result);
						}
					});
				} 
			});
		},

	    render: function(e) {
	    	$("#tilelayers_list").html("");
			var customer = GoSpatial.getCustomer();
			var tilelayers = customer.tilelayers;
			if (tilelayers) {
				for (var i=0; i < tilelayers.length; i++) {
					var elem = $("<tr><td>" + i + "</td><td>" + tilelayers[i].name + "</td><td>" + tilelayers[i].url + "</td></tr>");
					elem.id = tilelayers[i];
					console.log(tilelayers[i]);
					$("#tilelayers_list").append(elem);
				}
			}
	    	// Create charts container
			// var viewHtml = ''; 
			// Populate elements
			// this.$el.empty().append(viewHtml);
			// $(".app_container").append(viewHtml);
		}
	});
