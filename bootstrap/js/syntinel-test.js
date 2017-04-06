var project = {};
var test = {};
var r = [];

function pageLoad() {
	var projectID = getQueryVariable("project");
	var testID = getQueryVariable("test");

	var populatePage = function() {
		// Page header
		document.getElementById("header-test-name").innerHTML = test.name + " <small>Syntinel Test</small>";
		document.getElementById("breadcrumb-project-name").innerHTML = "<i class=\"fa fa-sitemap\"></i> <a href=\"project.html?project=" + projectID + "\">" + project.name + "</a>";
		document.getElementById("breadcrumb-test-name").innerHTML = "<i class=\"fa fa-file\"></i> " + test.name;

		// Save
		document.getElementById("setting-button-save").addEventListener('click', function() {
			var postBody = {};
			var script = document.getElementById("setting-testScript");
			if(script.value != test.script)
				postBody.script = script.value;

			var docker = document.getElementById("setting-testDocker");
			if(docker.value != test.dockerfile)
				postBody.dockerfile = docker.value;

			var envs = document.getElementById("setting-environmentVariables").childNodes;
			var envsChanged = false;
			if (envs !== undefined) {
				for(var i = 0; i < envs.length; i++) {
					var env = envs[i];
					env = env.childNodes[0].innerHTML + "=" + env.childNodes[1].innerHTML;

					if(!postBody.environmentVariables)
						postBody.environmentVariables = [];

					postBody.environmentVariables.push(env[0] + "=" + env[1]);

					if(env != test.environmentVariables) {
						envsChanged = true;
					}
				}
			}

			if(envsChanged == false) {
				delete postBody.environmentVariables;
			}

			var run = document.getElementById("setting-run");
			if(run.value == "off" && test.interval != null) {
				postBody.interval = null;
			} else if(run.value == "single") {
				apiPost(SYNTINEL_URL + "/testrun/", {"test" : Number(test.id)}, function(res) {
					if(res.error && SYNTINEL_ERRORREDIRECT) {
						var qs = {};
						if(res.responseText && res.responseText.length > 0) {
							qs.reason = res.responseText;
						}
						if(res.status) {
							qs.status = res.status;
						}
						
						qs.project = projectID;
						qs.test = testID;

						window.location = buildUrl("error.html", qs);
						return;
					}
					window.location = "test.html?project="+projectID+"&test="+testID;
				});
			} else if(run.value == "schedule" && test.interval == null) {
				postBody.interval = Number(document.getElementById("setting-run-interval").getElementsByTagName("input")[0].value);
			}

			if(postBody.script || postBody.dockerfile || postBody.environmentVariables || postBody.hasOwnProperty("interval")) {
				apiPatch(SYNTINEL_URL + "/test/" + testID, postBody, function(res) {
					if(res.error && SYNTINEL_ERRORREDIRECT) {
						var qs = {};
						if(res.responseText && res.responseText.length > 0) {
							qs.reason = res.responseText;
						}
						if(res.status) {
							qs.status = res.status;
						}
						
						qs.project = projectID;
						qs.test = testID;

						window.location = buildUrl("error.html", qs);
						return;
					}
					window.location = "test.html?project="+projectID+"&test="+testID;
				});
			}
		});

		document.getElementById("setting-button-delete").addEventListener('click', function() {
			apiDelete(SYNTINEL_URL + "/test/" + testID, {}, function(res) {
					if(res.error && SYNTINEL_ERRORREDIRECT) {
						var qs = {};
						if(res.responseText && res.responseText.length > 0) {
							qs.reason = res.responseText;
						}
						if(res.status) {
							qs.status = res.status;
						}
						
						qs.project = projectID;
						qs.test = testID;

						window.location = buildUrl("error.html", qs);
						return;
					}
				window.location = "project.html?project="+projectID;
			});
		});

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
		else {
			document.getElementById("setting-run").value = "schedule";
			document.getElementById("setting-run-interval").value = test.interval;
		}

		settingsRunChanged();

		// Setting - Environment Variables
		var envs = test.environmentVariables;
		var envStr = "";
		if (envs !== null) {
			for(var i = 0; i < envs.length; i++) {
				var env = envs[i].split("=");
				envStr += "<tr id=\"" + env[0] + "\"><td>" + env[0] + "</td><td>" + env[1] + "</td></tr>";
			}
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

		$('#table-test-runs').find('tr').click(function() {
			var index = ($(this).index());
			var row = document.getElementById("table-test-runs-body").childNodes[index];
			var id = row.childNodes[1].innerHTML;
			window.location = "run.html?project="+projectID+"&test="+testID+"&run="+id;
		});
	}

	// Get the project
	apiGet(SYNTINEL_URL + "/project/" + projectID, "", function(res) {
		if(res.error && SYNTINEL_ERRORREDIRECT) {
			var qs = {};
			if(res.responseText && res.responseText.length > 0) {
				qs.reason = res.responseText;
			}
			if(res.status) {
				qs.status = res.status;
			}
			
			qs.project = projectID;
			qs.test = testID;

			window.location = buildUrl("error.html", qs);
			return;
		}

		project = res;
		project = escapeNewLineChars(project);
		project = JSON.parse(project);

		// Now get the test
		apiGet(SYNTINEL_URL + "/test/" + testID, "", function(res) {
			if(res.error && SYNTINEL_ERRORREDIRECT) {
				var qs = {};
				if(res.responseText && res.responseText.length > 0) {
					qs.reason = res.responseText;
				}
				if(res.status) {
					qs.status = res.status;
				}
				
				qs.project = projectID;
				qs.test = testID;

				window.location = buildUrl("error.html", qs);
				return;
			}

			test = res;
			test = escapeNewLineChars(test);
			test = JSON.parse(test);

			// Get test runs
			apiGet(SYNTINEL_URL + "/testrun/all?test=" + testID, "", function(res) {
				if(res.error && SYNTINEL_ERRORREDIRECT) {
					var qs = {};
					if(res.responseText && res.responseText.length > 0) {
						qs.reason = res.responseText;
					}
					if(res.status) {
						qs.status = res.status;
					}
					
					qs.project = projectID;
					qs.test = testID;

					window.location = buildUrl("error.html", qs);
					return;
				}

				var runs = res;
				runs = escapeNewLineChars(runs);
				runs = JSON.parse(runs).results;

				for(var i = 0; i < runs.length; i++) {
					var run = runs[i];
					r.push(run);
				}

				populatePage();
			});
		});
	});
}

function settingsRunChanged() {
	var runSelect = document.getElementById("setting-run");
	if(runSelect.value == "schedule")
		document.getElementById("setting-run-interval").hidden = false;
	else
		document.getElementById("setting-run-interval").hidden = true;
}

var modalEnvs = [];

function updateModalsEnv() {
	var modalBody = document.getElementById("modal-env-body");
	modalBody.innerHTML = "";
	modalEnvs = [];

	var envs = test.environmentVariables;
	if(envs.hasOwnProperty("length")) {
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
