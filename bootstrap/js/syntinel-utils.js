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
	// Hard coded testing
	if(url.indexOf("/project/all") != -1) {
		callback("[{\"id\": 1,\"tests\": [1],\"name\": \"UltimateCode\"}]");
	} else if(url.indexOf("test/") != -1) {
		callback("{\"id\":1,\"name\":\"The greatest song in the world\",\"script\":\"#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover\",\"dockerfile\":\"FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh\",\"environmentVariables\":[\"a=b\"],\"health\":100,\"suite\":null}");
	} else {
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
}

String.prototype.replaceAll = function(search, replacement) {
  var target = this;
  return target.split(search).join(replacement);
};