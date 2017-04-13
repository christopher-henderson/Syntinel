var PI = 0;
var TI = 0;

var p = [];

function pageLoad() {
	document.getElementById("button-project-create").addEventListener('click', function() {
		// Open the modal
		$("#modal-create").modal();
	});

	var createProject = document.getElementById("modal-create-project").addEventListener('click', function() {
		apiPost(SYNTINEL_URL + "/project/", {"name" : document.getElementById("modal-create-project-name").value}, function(res) {
			if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
				var qs = {};
				if(res.responseText && res.responseText.length > 0) {
					qs.reason = res.responseText;
				}
				if(res.status) {
					qs.status = res.status;
				}
				window.location = buildUrl("error.html", qs);
				return;
			}

			var project = res;
			project = escapeNewLineChars(project);
			project = JSON.parse(project);

			window.location = "project.html?project="+project.id;
		});
	});

	// Make all the calls
	apiGet(SYNTINEL_URL + "/project/all", "", function(res) {
		if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
			var qs = {};
			if(res.responseText && res.responseText.length > 0) {
				qs.reason = res.responseText;
			}
			if(res.status) {
				qs.status = res.status;
			}
			window.location = buildUrl("error.html", qs);
			return;
		}

		var projects = res;
		projects = escapeNewLineChars(projects);
		projects = JSON.parse(projects).results;

		for(var i = 0; i < projects.length; i++) {
			var project = projects[i];
			p[i] = {"project" : project, "tests" : []};
		}

		projectLoaded();
	});
}

function projectLoaded() {
	if(PI < p.length) {
		var project = p[PI].project;
		if(TI < project.tests.length) {
			apiGet(SYNTINEL_URL + "/test/" + project.tests[TI], null, function(res) {
				if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
					var qs = {};
					if(res.responseText && res.responseText.length > 0) {
						qs.reason = res.responseText;
					}
					if(res.status) {
						qs.status = res.status;
					}
					window.location = buildUrl("error.html", qs);
					return;
				}

				p[PI].tests.push(JSON.parse(escapeNewLineChars(res)));
				TI++;
				projectLoaded();
			});
		} else {
			TI = 0;
			PI++;
			projectLoaded();
		}
	} else {
		populatePage();
	}
}

function populatePage() {
	var tabAll = document.getElementById("projectsTab-all");
	var tabPassing = document.getElementById("projectsTab-passing");
	var tabFailing = document.getElementById("projectsTab-failing");

	for(var i = 0; i < p.length; i++) {
		var project = p[i].project;
		var projectHealth = 0;

		var tab = "";
		tab += "<div aria-multiselectable=\"true\" class=\"panel-group\" id=\"project_" + project.id + "_#PROJECT_TAB#_accordion\" role=\"tablist\">";
		tab += "	<div class=\"panel panel-#PROJECT-HEALTH#\">";
		tab += "		<div class=\"panel-heading prev-emp activestate\" id=\"project_" + project.id + "_#PROJECT_TAB#_tab\" role=\"tab\">";
		tab += "			<h4 class=\"panel-title\"><a aria-controls=\"project_" + project.id + "_#PROJECT_TAB#_collapse\" aria-expanded=\"true\" data-parent=\"#project_" + project.id + "_#PROJECT_TAB#_accordion\" data-toggle=\"collapse\" href=\"#project_" + project.id + "_#PROJECT_TAB#_collapse\">" + p[i].project.name + " (Details)</a>";
		tab += "			<small> - <a href=\"project.html?project=" + project.id + "\">Open Project</a></small>";
		tab += "			</h4>";
		tab += "		</div>";
		tab += "		<div aria-labelledby=\"project_" + project.id + "_#PROJECT_TAB#_tab\" class=\"panel-collapse collapse\" id=\"project_" + project.id + "_#PROJECT_TAB#_collapse\" role=\"tabpanel\">";
		tab += "			<div class=\"container panel-body\">";
		// Project details
		tab += "				<h3>Project details";
		tab += "				<button type=\"button\" onclick=\"window.location=\'project.html?project=" + project.id + "\'\" class=\"btn btn-sm btn-info\">Open</button>";
		tab += "				</h3>";
		tab += "				<p>ID: " + project.id + "</p>";
		tab += "				<p>Name: " + project.name + "</p>";
		tab += "				<p>Test count: " + p[i].tests.length + "</p>";
		tab += "				<h3>Tests</h3>";
		// Test table
		tab += "				<div class=\"table-responsive\">";
		tab += "					<table id=\"table-project-" + project.id + "-#PROJECT_TAB#-tests\" class=\"table table-bordered table-hover table-striped syntinel-test-table\">";
		tab += "						<thead>";
		tab += "							<tr>";
		tab += "								<th>ID</th>";
		tab += "								<th>Name</th>";
		tab += "								<th>Health</th>";
		tab += "							</tr>";
		tab += "						</thead>";
		tab += "						<tbody>";
		// Tests
		// Sort tests
		p[i].tests.sort(function(a,b) {return (a.id > b.id) ? 1 : ((b.id > a.id) ? -1 : 0);} );
		for(var j = 0; j < p[i].tests.length; j++) {
			var test = p[i].tests[j];
			projectHealth += test.health;
			var testHealth = (test.health >= SYNTINEL_HEALTH.SUCCESS_MIN ? "success" : (test.health > SYNTINEL_HEALTH.WARN_MIN ? "warning" : "danger"));
			tab += "							<tr project=" + project.id + " class=\"" + testHealth + " syntinel-test-row\">";
			tab += "								<td>" + test.id + "</td>";
			tab += "								<td>" + test.name + "</td>";
			tab += "								<td>" + test.health + "%</td>";
			tab += "							</tr>";
		}
		tab += "						</tbody>";
		tab += "					</table>";
		tab += "				</div>";
		tab += "			</div>";
		tab += "		</div>";
		tab += "	</div>";
		tab += "</div>";

		// Change project health color
		projectHealth = projectHealth / p[i].tests.length;
		projectHealth = (projectHealth >= SYNTINEL_HEALTH.SUCCESS_MIN ? "success" : (projectHealth > SYNTINEL_HEALTH.WARN_MIN ? "warning" : "danger"));
		tab = tab.replace("#PROJECT-HEALTH#", projectHealth);

		// Add to main "All" tab
		tabAll.innerHTML += tab.replaceAll("#PROJECT_TAB#", "all");

		// Add to other tabs
		if(projectHealth == "success") {
			tabPassing.innerHTML += tab.replaceAll("#PROJECT_TAB#", "passing");
		} else if(projectHealth == "warning" || projectHealth == "danger") {
			tabFailing.innerHTML += tab.replaceAll("#PROJECT_TAB#", "failing");
		}

	}

	var testTables = document.getElementsByClassName("syntinel-test-table");
	for(var i = 0; i < testTables.length; i++) {
		var table = testTables[i];
		var rows = table.getElementsByTagName("tbody")[0].getElementsByTagName("tr");

		for(var j = 0; j < rows.length; j++) {
			var row = rows[j];
			var handler = function(table, row) {
				return function() {
					var testID = row.childNodes[1].innerHTML;
					var projectID = table.id.split("-")[2];
					window.location = ("test.html?project="+ projectID +"&test=" + testID);
				}
			};

			row.onclick = handler(table, row);
		}
	}
}
