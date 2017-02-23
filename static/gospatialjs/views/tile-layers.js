
	var TileLayerView = Backbone.View.extend({
	    el: "#tileLayerView",

	    initialize: function(apikey){
	        // var lp = this.$(".rightpanel");
	        _.bindAll(this, 'render', 'createTileLayer');
	        this.apikey = apikey;
	        this.gospatial = new GoSpatialApi(apikey);
	        this.customer;
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
				text: "Are you sure you want to create a new tile layer",
				type: "info",
				showCancelButton: true,
				confirmButtonColor: "#DD6B55",
				confirmButtonText: "Yes, pls!",
				cancelButtonText: "No, cancel pls!",
				//closeOnConfirm: false,
				//closeOnCancel: true,
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
	    	var self = this;
	    	$("#tilelayers_list").html("");
			this.gospatial.getCustomer(function(error,result) {
				if (error) {
					swal("Error", error, "error");
					self.customer = undefined;
					return;
				}
				self.customer = result;

				if (!self.customer.hasOwnProperty("tilelayers")) {
					swal("Error", "Invalid customer object: " + JSON.stringify(self.customer),"error");
					return;
				}

				if (!self.customer.tilelayers) {return;}

				for (var i=0; i < self.customer.tilelayers.length; i++) {
					var elem = $("<tr><td>" + i + "</td><td>" + self.customer.tilelayers[i].name + "</td><td>" + self.customer.tilelayers[i].url + "</td></tr>");
					elem.id = self.customer.tilelayers[i];
					$("#tilelayers_list").append(elem);
				}

			});

	    	// Create charts container
			// var viewHtml = ''; 
			// Populate elements
			// this.$el.empty().append(viewHtml);
			// $(".app_container").append(viewHtml);
		}
	});
