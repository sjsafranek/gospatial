
	var VectorLayerView = Backbone.View.extend({
		el: "#vectorLayerView",

		initialize: function(apikey) {
			console.log(apikey);
			_.bindAll(this, 'render', 'createLayer', 'deleteLayer', 'toggleVectorLayerOptions');
			this.vectorlayers = new VectorLayerCollection;
			this.apikey = apikey;
			this.gospatial = new GoSpatialApi(apikey);
			this.customer;
			this.render();
			return this;
		},

	    events: {
	        "click button.deleteLayer" : "deleteLayer",
	        "click button.createLayer": "createLayer",
	        "click button.toggleVectorLayerOptions": "toggleVectorLayerOptions"
	    },

		createLayer: function() {
			var self = this;

			new SwalConfirm( 
				"Create layer?", 
				"Are you sure you want to create a new layer?", 
				"info",
				function(){
					$.ajax({
						url: '/api/v1/layer?apikey=' + self.apikey,
						type: 'POST',
						success: function(result) {
							new SwalPpSuccess("Created!", result);
							$("#layers_list").html("");
							self.render();
						},
						failure: function(result) {
							new SwalPpError("ApiError", result);
							throw new Error(result);
						},
						error: function(result) {
							new SwalPpError("ApiError", result);
							throw new Error(result);
						}
					}); 
				}
			);

			// swal({
			// 	title: "Create layer",
			// 	text: "Are you sure you want to create a new layer",
			// 	type: "info",
			// 	showCancelButton: true,
			// 	confirmButtonColor: "#337ab7",
			// 	confirmButtonText: "Yes, pls!",
			// 	cancelButtonText: "No, cancel pls!",
			// 	showLoaderOnConfirm: true,
			// }).then(function(){
			// 	$.ajax({
			// 		url: '/api/v1/layer?apikey=' + self.apikey,
			// 		type: 'POST',
			// 		success: function(result) {
			// 			new SwalPpSuccess("Created!", result);
			// 			$("#layers_list").html("");
			// 			self.render();
			// 		},
			// 		failure: function(result) {
			// 			new SwalPpError("ApiError", result);
			// 			throw new Error(result);
			// 		},
			// 		error: function(result) {
			// 			new SwalPpError("ApiError", result);
			// 			throw new Error(result);
			// 		}
			// 	}); 
			// });

		},

		deleteLayer: function(event) {
			var self = this;
			var datasource_id = $(event.target).attr("ds_id");
			// HANDLE IF CLICK ON ICON!!!

			new SwalConfirm( 
				"Delete layer?", 
				"Are you sure you want to delete " + datasource_id,
				"warning",
				function(){
					$.ajax({
						url: '/api/v1/layer/'+ datasource_id +'?apikey=' + self.apikey,
						type: 'DELETE',
						success: function(result) {
							new SwalPpSuccess("Deleted!", result);
							var model = self.vectorlayers.get(datasource_id);
							self.vectorlayers.remove(model);
							self.render();
						},
						failure: function(result) {
							new SwalPpError("ApiError", result);
							throw new Error(result);
						},
						error: function(result) {
							new SwalPpError("ApiError", result);
							throw new Error(result);
						}
					});
				}, 
				function(dismiss){
					if ('cancel' == dismiss) {
						swal("Cancelled", "Your data is safe :)", "error");
					}
				}
			);

			// swal({
			// 	title: "Delete layer",
			// 	text: "Are you sure you want to delete " + datasource_id,
			// 	type: "warning",	 
			// 	showCancelButton: true,	 
			// 	confirmButtonColor: "#DD6B55",	 
			// 	confirmButtonText: "Yes, delete it!",	 
			// 	cancelButtonText: "No, cancel pls!",	 
			// 	showLoaderOnConfirm: true
			// }).then(
			// 	function(){
			// 		$.ajax({
			// 			url: '/api/v1/layer/'+ datasource_id +'?apikey=' + self.apikey,
			// 			type: 'DELETE',
			// 			success: function(result) {
			// 				new SwalPpSuccess("Deleted!", result);
			// 				var model = self.vectorlayers.get(datasource_id);
			// 				self.vectorlayers.remove(model);
			// 				self.render();
			// 			},
			// 			failure: function(result) {
			// 				new SwalPpError("ApiError", result);
			// 				throw new Error(result);
			// 			},
			// 			error: function(result) {
			// 				new SwalPpError("ApiError", result);
			// 				throw new Error(result);
			// 			}
			// 		});
			// 	}, 
			// 	function(dismiss){
			// 		if ('cancel' == dismiss) {
			// 			swal("Cancelled", "Your data is safe :)", "error");
			// 		}
			// 	}
			// );
		},

		viewLayer: function() {
			var self = this;
			var datasource_id = $(event.target).attr("ds_id");
			self.gospatial.getLayer(datasource_id, function(error, data) {
				console.log(data);
				if (error) {
					throw new Error(error);
				} else {
					$(".raw-json").html(JSON.stringify(data, null, 2));
				}
			});
		},

		toggleVectorLayerOptions: function(event) {
			var self = this;
			var datasource_id = $(event.target).attr("ds_id");
			$(".vectorlayer").hide();
			$(".vectorlayer").each(function() {
				if ($(this).attr("ds_id") == datasource_id) {
					$(this).show();

					// download json file
					var downloadFile = null;
					var makeGeoJSONFile = function (text) {
						// var data = new Blob([text], {type: 'text/plain'});
						var data = new Blob([text], {type: 'json'});
						// If we are replacing a previously generated file we need to
						// manually revoke the object URL to avoid memory leaks.
						if (downloadFile !== null) {
							window.URL.revokeObjectURL(downloadFile);
						}
						downloadFile = window.URL.createObjectURL(data);
						return downloadFile;
					};

					var container = this; 
					// if code.raw-json is empty request for geojson
					if ($(container).find(".raw-json").html() == "") {
						self.gospatial.getLayer(datasource_id, function(error, data) {
							if (error) {
								SwalPpError("ApiError", error)
								return;
							}
							// $(container).find(".raw-json").html(JSON.stringify(data, null, 2));
							var metadata = {
								size: JSON.stringify(data).length,
								features: data.features.length
							};
							$(container).find(".raw-json").html(JSON.stringify(metadata, null, 2));
							// download link for geojson file
							var link = $(container).find("a");
							link.href = makeGeoJSONFile(data);
						});
					}
				}
			});
		},

		render: function() {
			var self = this;
			$("#layers_list").html("");
			
			this.gospatial.getCustomer(function(error,result) {
				if (error) {
					SwalPpError("ApiError", error);
					self.customer = undefined;
					return;
				}
				self.customer = result;

				if (!self.customer.hasOwnProperty("datasources")) {
					SwalPpError("ApiResponseError", self.customer);
					return;
				}

				for (var i=0; i < self.customer.datasources.length; i++) {
					var lyr = new VectorLayer({id: self.customer.datasources[i]});
					self.vectorlayers.add(lyr);
				}

				self.vectorlayers.each(function(model) {
					var ds = model.get("id");
					var html =  '<div class="panel panel-default">' +
									'<div class="panel-heading">' + ds +
										'<div class="panel_controls">' +
											'<button type="button" title="options" ds_id=' + ds + ' class="btn btn-default btn-sm toggleVectorLayerOptions">' + 
												'<i class="fa fa-cog" aria-hidden="true"></i>' + 
											'</button>' +
										'</div>' +
									'</div>' +
									'<div class="panel-body vectorlayer" ds_id=' + ds + '>' + 
										'<div class="col-md-11 column">' +
											'<div class="well">' +
												'<code class="raw-json"></code>' +
											'</div>' +
										'</div>' +
										'<div class="col-md-1 column">' +
											'<button class="btn btn-sm btn-info downloadLayer" title="view" ds_id=' + ds + '>' + 
												'<a href="/api/v1/layer/'+ ds +'?apikey=' + self.apikey + '" download="' + ds + '.geojson">' +
													// '<i class="fa fa-cloud-download" aria-hidden="true"></i>' + 
													'<i class="fa fa-download" aria-hidden="true"></i>' + 
												'</a>' + 
											'</button>' +
											'<button class="btn btn-sm btn-danger deleteLayer" title="delete" ds_id=' + ds + '>' + 
												'<i class="fa fa-trash"></i>' +
											'</button>' +
										'</div>' +
									'</div>' +
								'</div>';
					var elem = $(html);
					elem.id = ds;
					$("#layers_list").append(elem);
				});
				$(".vectorlayer").hide();

			});

		}

	});

