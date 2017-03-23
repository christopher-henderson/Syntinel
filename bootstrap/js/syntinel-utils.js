var SYNTINEL_URL = "http://syntinel.chenderson.org/api/v1";

function getQueryVariable(variable) {
	var query = window.location.search.substring(1);
	var vars = query.split("&");
	for (var i=0;i<vars.length;i++) {
		var pair = vars[i].split("=");
		if(pair[0] == variable){return pair[1];}
	}
}

function escapeNewLineChars(valueToEscape) {
	if (valueToEscape != null && valueToEscape != "") {
		return valueToEscape.replace(/\n/g, "\\n");
	} else {
		return valueToEscape;
	}
}

function apiGet(url, params, callback) {
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("GET", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send();
	request.onreadystatechange = function() {
		if(request.readyState >= 4)
			handler(request);
	};
}