
			function GoSpatialApi(apikey, server) {

				this.apikey = apikey;
				this.server = server || "";

				this.getCustomer = function() {
					var self = this;
					var data;
					this.GET("/api/v1/customer" + "?apikey=" + self.apikey, function(error, result){
						if (error) {
							throw error;
						} else {
							data = result;
						}
					});
					return data;
				}

				this.getLayer = function(datasource, callback) {
					var self = this;
					this.GET("/api/v1/layer/" + datasource + "?apikey=" + self.apikey, function(error, result){
						callback(error, result);
					});
				}

				this.submitFeature = function(datasource, feature, callback) {
					var self = this;
					this.POST(
						'/api/v1/layer/' + datasource + '/feature?apikey=' + self.apikey,
						feature,
						function(error, result) {
							callback(error, result);
						}
					)
				}

				this.GET = function(route, callback) {
					var self = this;
					$.ajax({
						crossDomain: true,
						type: "GET",
						async: false,
						url: route,
						dataType: 'JSON',
						success: function (data) {
							return callback(null, data);
						},
						error: function(xhr,errmsg,err) {
							console.log(xhr.status,xhr.responseText,errmsg,err);
							result = null;
							var message = "status: " + xhr.status + "<br>";
							message += "responseText: " + xhr.responseText + "<br>";
							message += "errmsg: " + errmsg + "<br>";
							message += "Error:" + err;
							return callback(new Error(message));
						}
					});
				}

				this.POST = function(route, data, callback) {
					var self = this;
					$.ajax({
						crossDomain: true,
						type: "POST",
						async: false,
						data: data,
						url: route,
						dataType: 'JSON',
						success: function (data) {
							callback(null, data);
						},
						error: function(xhr,errmsg,err) {
							console.log(xhr.status,xhr.responseText,errmsg,err);
							result = null;
							var message = "status: " + xhr.status + "<br>";
							message += "responseText: " + xhr.responseText + "<br>";
							message += "errmsg: " + errmsg + "<br>";
							message += "Error:" + err;
							callback(new Error(message));
						}
					});
				}
			}

