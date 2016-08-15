
	function GoSpatialApi(apikey, server) {

		this.apikey = apikey;
		this.server = server || "";
		// this.url = new URL("http://localhost:8888/map?apikey=iNx1xvBPDrZb#");

		this.getCustomer = function(callback) {
			var self = this;
			this.GET(this.server + "/api/v1/customer" + "?apikey=" + self.apikey, function(error, result){
				return callback(error, result);
			});
		}

		this.getLayer = function(datasource, callback) {
			var self = this;
			this.GET(this.server + "/api/v1/layer/" + datasource + "?apikey=" + self.apikey, function(error, result){
				return callback(error, result);
			});
		}

		this.submitFeature = function(datasource, feature, callback) {
			var self = this;
			this.POST(
				this.server + '/api/v1/layer/' + datasource + '/feature?apikey=' + self.apikey,
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
				// async: false,
				url: route,
				dataType: 'JSON',
				success: function (data) {
					return callback(null, data);
				},
				error: function(xhr,errmsg,err) {
					console.log(xhr.status,xhr.responseText,errmsg,err);
					console.log(xhr);
					// var message = "status: " + xhr.status + "<br>";
					// message += "responseText: " + xhr.responseText + "<br>";
					// message += "errmsg: " + errmsg + "<br>";
					// message += "Error:" + err;
					var message = xhr.status + " " + xhr.responseText;
					return callback(new Error(message));
				}
			});
		}

		this.POST = function(route, data, callback) {
			var self = this;
			$.ajax({
				crossDomain: true,
				type: "POST",
				// async: false,
				data: data,
				url: route,
				dataType: 'JSON',
				success: function (data) {
					callback(null, data);
				},
				error: function(xhr,errmsg,err) {
					console.log(xhr.status,xhr.responseText,errmsg,err);
					console.log(xhr);
					// var message = "status: " + xhr.status + "<br>";
					// message += "responseText: " + xhr.responseText + "<br>";
					// message += "errmsg: " + errmsg + "<br>";
					// message += "Error:" + err;
					var message = xhr.status + " " + xhr.responseText;
					callback(new Error(message));
				}
			});
		}
	}

