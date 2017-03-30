var SYNTINEL_URL = "http://syntinel.chenderson.org/api/v1";
var SYNTINEL_HEALTH = {
	SUCCESS_MIN : 87.5,
	WARN_MIN : 75
}

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
		valueToEscape = valueToEscape.replace(/\r/g, "\\r")
		return valueToEscape.replace(/\n/g, "\\n");
	} else {
		return valueToEscape;
	}
}

function apiGet(url, params, callback) {
	// Hard coded testing
	if(url.indexOf("/project/all") != -1) {
		callback("[{\"id\": 1,\"tests\": [1],\"name\": \"UltimateCode\"}]");
	} else if(url.indexOf("/project/1") != -1) {
		callback("{\"id\": 1,\"tests\": [1],\"name\": \"UltimateCode\"}");
	} else if(url.indexOf("/test/") != -1) {
		callback("{\"id\":1,\"name\":\"The greatest song in the world\",\"script\":\"#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover\",\"dockerfile\":\"FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh\",\"environmentVariables\":[\"a=b\"],\"health\":" + getRandomIntInclusive(0,100) + ",\"suite\":null}");
	} else if(url.indexOf("/testrun/") != -1) {
		callback("[{\"id\":1,\"log\":\"Sending build context to Docker daemon 3.072 kB\r\nStep 1 : FROM docker.io/centos\n ---> 67591570dd29\nStep 2 : MAINTAINER Christopher Henderson\n ---> Using cache\n ---> d33126dcfbfc\nStep 3 : RUN yum install -y go git wget\nCloning into 'TestTheTester'...\n ---> Using cache\n ---> b4da48a3d37c\nStep 4 : COPY script.sh $HOME/script.sh\n ---> Using cache\n ---> 40743fa214db\nStep 5 : CMD chmod +x script.sh && ./script.sh\n ---> Using cache\n ---> 1da95f92dcd2\nSuccessfully built 1da95f92dcd2\n=== RUN   TestPow\n--- PASS: TestPow (0.00s)\n=== RUN   TestSquare\n--- PASS: TestSquare (0.00s)\nPASS\ncoverage: 46.2% of statements\n\",\"error\":\"\",\"status\":null,\"successful\":true,\"test\":1}]");
	}  else {
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

function getRandomIntInclusive(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}