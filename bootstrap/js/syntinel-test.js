var project = {};
var test = {};
var r = [];

function pageLoad() {
	var projectID = getQueryVariable("project");
	var testID = getQueryVariable("test");

	var populatePage = function() {
		// Page header
		document.getElementById("header-test-name").innerHTML = test.name + " <small>Syntinel Test</small>";
		document.getElementById("breadcrumb-project-name").innerHTML = "<i class=\"fa fa-file\"></i> <a href=\"project.html?project=" + projectID + "\">" + project.name + "</a>";
		document.getElementById("breadcrumb-test-name").innerHTML = "<i class=\"fa fa-file\"></i> " + test.name;

		// Setting - Name
		document.getElementById("setting-project").value = projectID;

		// Setting - Name
		document.getElementById("setting-testName").value = test.name;

		// Setting - ID
		document.getElementById("setting-testID").value = test.id;

		// Setting - Script
		document.getElementById("setting-testScript").value = test.script;

		// Setting - Docker
		document.getElementById("setting-testDocker").value = test.dockerfile;

		// Setting - Run
		if(!test.interval || test.interval == null)
			document.getElementById("setting-run").value = "off";
		else
			document.getElementById("setting-run").value = "schedule";

		settingsRunChanged();

		// Setting - Environment Variables
		var envs = test.environmentVariables;
		var envStr = "";
		for(var i = 0; i < envs.length; i++) {
			var env = envs[i].split("=");
			envStr += "<tr id=\"" + env[0] + "\"><td>" + env[0] + "</td><td>" + env[1] + "</td></tr>";
		}
		document.getElementById("setting-environmentVariables").innerHTML = envStr;

		var buttonEnvAdd = document.getElementById("setting-env-button-edit").addEventListener('click', function() {
			updateModalsEnv();

			// Open the modal
			$("#modal-env").modal();
		});

		// Run histories
		var testRuns = document.getElementById("table-test-runs-body");
		testRuns.innerHTML = "";
		for(var i = 0; i < r.length; i++) {
			var run = r[i];
			var runRow = "";

			var runStatus;
			if(!run.successful || run.successful == null) {
				runStatus = "Running";
			} else if(run.successful == true) {
				runStatus = "Successful";
			} else {
				runStatus = "Failed";
			}

			runRow += "<tr class=\"" + (runStatus == "Successful" ? "success" : (runStatus == "Running" ? "warning" : "danger")) + "\">";
			runRow += "	<td>" + run.id + "</td>";
			runRow += "	<td>" + runStatus + "</td>";
			runRow += "</tr>";

			testRuns.innerHTML += runRow;
		}
	}

	// Get the project
	apiGet(SYNTINEL_URL + "/project/" + projectID, "", function(res) {
		project = res;
		project = escapeNewLineChars(project);
		project = JSON.parse(project);

		// Now get the test
		apiGet(SYNTINEL_URL + "/test/" + testID, "", function(res) {
			test = res;
			test = escapeNewLineChars(test);
			test = JSON.parse(test);

			// Get test runs
			apiGet(SYNTINEL_URL + "/testrun/all?test=" + testID, "", function(res) {
				var runs = res;
				runs = escapeNewLineChars(runs);
				runs = JSON.parse(runs);

				for(var i = 0; i < runs.length; i++) {
					var run = runs[i];
					r.push(run);
				}

				populatePage();
			});
		});
	});

	$('#table-test-runs').find('tr').click(function() {
		var index = ($(this).index());
		var row = document.getElementById("table-test-runs-body").childNodes[index];
		var id = row.childNodes[1].innerHTML;
		window.location = "run.html?project="+projectID+"&test="+testID+"&run="+id;
	});
}

function settingsRunChanged() {
	var runSelect = document.getElementById("setting-run");
	if(runSelect.value == "schedule")
		document.getElementById("setting-run-interval").hidden = false;
	else
		document.getElementById("setting-run-interval").hidden = true;
}

var modalTest = {};
var modalEnvs = [];

function updateModalsEnv() {
	var modalBody = document.getElementById("modal-env-body");
	modalBody.innerHTML = "";
	modalEnvs = [];

	var envs = test.environmentVariables;
	for(var i = 0; i < envs.length + 1; i++) {
		var env;
		if(i < envs.length) {
			env = envs[i].split("=");
		} else {
			env = ["",""];
		}

		modalEnvs[i] = env;

		var str = "";
		str += "<div class=\"row\">";
		str += "	<div class=\"col-lg-5\">";
		str += "		<input onchange=\"modalEnvInputChanged(" + i + ")\" id=\"modal-env-variable-index-" + i + "\" class=\"form-control\" placeholder=\"Variable\" value=\"" + env[0] + "\">";
		str += "	</div>";
		str += "	<div class=\"col-lg-5\">";
		str += "		<input onchange=\"modalEnvInputChanged(" + i + ")\" id=\"modal-env-value-index-" + i + "\" class=\"form-control\" placeholder=\"Value\" value=\"" + env[1] + "\">";
		str += "	</div>";
		str += "</div>";

		modalBody.innerHTML += str;
	}
}

function modalEnvInputChanged(index) {
	var envVariable = document.getElementById("modal-env-variable-index-" + index);
	var envValue = document.getElementById("modal-env-value-index-" + index);

	// Update the model
	if(index == modalEnvs.length - 1) {
		// Add
		test.environmentVariables.push(envVariable.value + "=" + envValue.value);
	} else {
		if(envVariable.value.length <= 0 && envValue.value.length <= 0) {
			// Remove
			test.environmentVariables.splice(index, 1);
		} else {
			// Edit
			test.environmentVariables[index] = envVariable.value + "=" + envValue.value;
		}
	}

	// Update main page
	var envs = test.environmentVariables;
	var envStr = "";
	for(var i = 0; i < envs.length; i++) {
		var env = envs[i].split("=");
		envStr += "<tr id=\"" + env[0] + "\"><td>" + env[0] + "</td><td>" + env[1] + "</td></tr>";
	}
	document.getElementById("setting-environmentVariables").innerHTML = envStr;

	updateModalsEnv();
}