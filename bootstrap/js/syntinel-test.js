function pageLoad() {
	var id = getQueryVariable("idTest");
	var test;
	if(typeof id !== "undefined") {
		// TODO API call to get current test
		test = "{\"id\":1,\"name\":\"The greatest song in the world\",\"script\":\"#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover\",\"dockerfile\":\"FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh\",\"environmentVariables\":[\"a=b\"],\"health\":100,\"suite\":null}";
		test = escapeNewLineChars(test);
		try {
			test = JSON.parse(test);
		} catch(err) {
			console.log(err);
		}
	}

	var projects = [
		{name : "CompoZed"}
	]

	// Setting - Project dropdown
	var projectNames = "";
	for(var i = 0; i < projects.length; i++) {
		projectNames += "<option>" + projects[i].name + "</option>";
	}
	if(projectNames.length == 0) {
		projectNames += "<option>No project selected</option>";
	}
	document.getElementById("setting-project").innerHTML = projectNames;

	if(typeof test === "undefined") {
		// Page header
		document.getElementById("header-test-name").innerHTML = "Create new test <small>Syntinel Test</small>";
		document.getElementById("breadcrumb-test-name").innerHTML = "<i class=\"fa fa-file\"></i> Create new test";
	} else {
		// Page header
		document.getElementById("header-test-name").innerHTML = test.name + " <small>Syntinel Test</small>";
		document.getElementById("breadcrumb-test-name").innerHTML = "<i class=\"fa fa-file\"></i> " + test.name;

		// Setting - Name
		document.getElementById("setting-testName").value = test.name;

		// Setting - ID
		document.getElementById("setting-testID").value = test.id;

		// Setting - Script
		document.getElementById("setting-testScript").value = test.script;

		// Setting - Docker
		document.getElementById("setting-testDocker").value = test.dockerfile;

		// Setting - Environment Variables
		var envs = test.environmentVariables;
		var envStr = "";
		for(var i = 0; i < envs.length; i++) {
			var env = envs[i].split("=");
			envStr += "<tr><td>" + env[0] + "</td><td>" + env[1] + "</td><td>";
			envStr += "<button type=\"button\" class=\"btn btn-xs btn-info\">Edit</button>";
			envStr += "<button type=\"button\" class=\"btn btn-xs btn-danger\">Remove</button>";
			envStr += "</td></tr>";
		}
		document.getElementById("setting-environmentVariables").innerHTML = envStr;
	}
}