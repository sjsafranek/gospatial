
var Utils = {

    postRequest: function(route, data) {
        var results;
        $.ajax({
            crossDomain: true,
            // dataType: 'jsonp',
            type: "POST",
            async: false,
            data: data,
            url: route,
            dataType: 'JSON',
            success: function (data) {
                try {
                    results = data;
                }
                catch(err){  console.log('Error:', err);  }
            },
            error: function(xhr,errmsg,err) {
                console.log(xhr.status,xhr.responseText,errmsg,err);
                result = null;
            }
        });
        return results;
    },

    getRequest: function(route, data) {
        var results;
        $.ajax({
            crossDomain: true,
            // dataType: 'jsonp',
            type: "GET",
            async: false,
            data: data,
            url: route,
            dataType: 'JSON',
            success: function (data) {
                try {
                    results = data;
                }
                catch(err){  console.log('Error:', err);  }
            },
            error: function(xhr,errmsg,err) {
                console.log(xhr.status,xhr.responseText,errmsg,err);
                result = null;
            }
        });
        return results;
    }

}