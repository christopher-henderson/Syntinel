function pageLoad() {
	var p = [];
	var projectCount = 0;

	var populatePage = function() {
		var tabAll = document.getElementById("projectsTab-all");
		var tabPassing = document.getElementById("projectsTab-passing");
		var tabFailing = document.getElementById("projectsTab-failing");

		if(p.length < projectCount) {
			return;
		}

		for(var i = 0; i < p.length; i++) {
			var project = p[i].project;
			var projectHealth = 0;

			var tab = "";
			tab += "<div aria-multiselectable=\"true\" class=\"panel-group\" id=\"project_" + project.id + "_#PROJECT_TAB#_accordion\" role=\"tablist\">";
			tab += "	<div class=\"panel panel-#PROJECT-HEALTH#\">";
			tab += "		<div class=\"panel-heading prev-emp activestate\" id=\"rc_branch\" role=\"tab\">";
			tab += "			<h4 class=\"panel-title\"><a aria-controls=\"project_" + project.id + "_#PROJECT_TAB#_collapse\" aria-expanded=\"true\" data-parent=\"#project_" + project.id + "_#PROJECT_TAB#_accordion\" data-toggle=\"collapse\" href=\"#collapseOne\">" + p[i].project.name + "</a>";
			tab += "			</h4>";
			tab += "		</div>";
			tab += "		<div aria-labelledby=\"rc_branch\" class=\"panel-collapse collapse\" id=\"project_" + project.id + "_#PROJECT_TAB#_collapse\" role=\"tabpanel\">";
			tab += "			<div class=\"container panel-body\">";
			// Project details
			tab += "				<h3>Project details</h3>";
			tab += "				<p>[" + project.id + "] " + project.name + "</p>";
			tab += "				<h3>Tests</h3>";
			// Test table
			tab += "				<div class=\"table-responsive\">";
			tab += "					<table class=\"table table-bordered table-hover table-striped\">";
			tab += "						<thead>";
			tab += "							<tr>";
			tab += "								<th>ID</th>";
			tab += "								<th>Name</th>";
			tab += "								<th>Health</th>";
			tab += "							</tr>";
			tab += "						</thead>";
			tab += "						<tbody>";
			// Tests
			for(var j = 0; j < p[i].tests.length; j++) {
				var test = p[i].tests[j];
				projectHealth += test.health;
				var testHealth = (test.health >= 100 ? "success" : (test.health > 50 ? "warning" : "danger"));
				tab += "							<tr class=\"" + testHealth + "\">";
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
			projectHealth = (projectHealth >= 100 ? "success" : (projectHealth > 50 ? "warning" : "danger"));
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
	}

	// Make all the calls
	apiGet(SYNTINEL_URL + "/project/all", "", function(res) {
		var projects = res;
		projects = escapeNewLineChars(projects);
		projects = JSON.parse(projects);
		projectCount = projects.length;

		for(var i = 0; i < projects.length; i++) {
			var project = projects[i];

			var tests = [];
			var count = 0;
			for(var j = 0; j < project.tests.length; j++) {
				apiGet("/test/" + project.tests[i], "", function(res) {
					tests.push(JSON.parse(escapeNewLineChars(res)));
					count++;
					if(count == project.tests.length) {
						p.push({"project" : project, "tests" : tests});
						populatePage();
					}
				});
			}
		}
	});
}