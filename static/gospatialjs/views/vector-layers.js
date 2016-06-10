
	var VectorLayerView = Backbone.View.extend({
		el: "#vectorLayerView",

		initialize: function(apikey) {
			console.log(apikey);
			_.bindAll(this, 'render', 'createLayer', 'deleteLayer', 'viewLayer');
			this.vectorlayers = new VectorLayerCollection;
			this.apikey = apikey;
			this.render();
			return this;
		},

	    events: {
	        "click button.deleteLayer" : "deleteLayer",
	        "click button.createLayer": "createLayer",
	        "click button.viewLayer": "viewLayer"
	    },

		createLayer: function() {
			var self = this;
			swal({
				title: "Create layer",
				text: "Are you sure you want to create a new lyaer",
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
						url: '/api/v1/layer?apikey=' + self.apikey,
						type: 'POST',
						success: function(result) {
							swal("Created!", result, "success");
							$("#layers_list").html("");
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

		deleteLayer: function(event) {
			var self = this;
			var datasource_id = $(event.target).attr("ds_id");
			// HANDLE IF CLICK ON ICON!!!
			swal({
				title: "Delete layer",
				text: "Are you sure you want to delete " + datasource_id,
				type: "warning",	 
				showCancelButton: true,	 
				confirmButtonColor: "#DD6B55",	 
				confirmButtonText: "Yes, delete it!",	 
				cancelButtonText: "No, cancel pls!",	 
				closeOnConfirm: false,	 
				closeOnCancel: false,
				showLoaderOnConfirm: true
			},
			function(isConfirm){
				if (isConfirm) {
					$.ajax({
						url: '/api/v1/layer/'+ datasource_id +'?apikey=' + self.apikey,
						type: 'DELETE',
						success: function(result) {
							swal("Deleted!", result, "success");
							var model = self.vectorlayers.get(datasource_id);
							self.vectorlayers.remove(model);
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
				else {
					swal("Cancelled", "Your data is safe :)", "error");
				}
			});
		},

		viewLayer: function() {
			var self = this;
			var datasource_id = $(event.target).attr("ds_id");
			var url = '/api/v1/layer/'+ datasource_id +'?apikey=' + self.apikey
			console.log(url);
		},

		render: function() {
			var self = this;
			$("#layers_list").html("");
			this.get_vector_layers();
			self.vectorlayers.each(function(model) {
				var ds = model.get("id");
				console.log(model);
				var html = '';
				html += "<tr>"; // row start
				html += "<td>" + ds + "</td>"; // cell
				html += '<td><button class="btn btn-sm btn-info viewLayer" title="view" ds_id=' + ds + '><i class="fa fa-file-text" aria-hidden="true"></i></button></li></td>'; // cell
				html += "<td><button class='btn btn-sm btn-danger deleteLayer' title='delete' ds_id=" + ds + "><i class='fa fa-trash'></i></button></li></td>"; // cell
				html += "</tr>" // row end						
				var elem = $(html);
				elem.id = ds;
				$("#layers_list").append(elem);
			});
			// return this;
		},

		get_vector_layers: function() {
			var self = this;
			var customer = GoSpatial.getCustomer();
			var datasources = customer.datasources;
			if (datasources) {
				for (var i=0; i < datasources.length; i++) {
					var lyr = new VectorLayer({id: datasources[i]});
					self.vectorlayers.add(lyr);
				}
			}
		}

	});
