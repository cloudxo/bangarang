function EscalationController($scope, $http, $cookies, $mdDialog) {
	$scope.escalations = [];
	$scope.fetchEscalations = function() {
		$http.get("api/escalation/config/*").success(function(data, status) {
			if (data) {
				if (data != "null") {
					$scope.escalations = data;
				}
			}
		});
	}
	this.selected = 0;
	this.getSelected = function() {
		var s = $cookies.get("nec:tab");
		if (s) {
			this.selected = s;
		}
		return this.selected;
	}
	this.updateSelected = function(name) {
		$cookies.put("nec:tab", name);
		this.selected = name;
	}

	$scope.removeSure = {}
	$scope.showRemoveDialog = function(name) {
		$scope.removeSure[name] = true;
	}

	$scope.hideRemoveDialog = function(name) {
		$scope.removeSure[name] = false;
	}

	$scope.shouldHideRemoveDialog = function(name) {
		var show = $scope.removeSure[name];
		return show != true;
	}

	$scope.removeEscalation = function(name)  {
		$http.delete("api/escalation/config/"+name).success(function(data) {
			$scope.fetchEscalations();
		});
	}

	$scope.fetchEscalations();

}
angular.module("bangarang").controller("EscalationController", EscalationController);


function NewEscalationController($scope, $http, $interval) {
	this.name = "";
	this.type = null;
	this.ots = {};

	this.type_list = [
		{
			title: "Pagerduty",
			name: "pager_duty"
		},
		{
			title: "Email",
			name: "email",
		},
		{
			title: "Console",
			name: "console"
		},
		{
			title: "Grafana Graphite Annotation",
			name: "grafana_graphite_annotation"
		}

	]

	this.pdOpts = [
		{
			title:"Api Key",
			name: "key",
			value: ""
		},
		{
			title: "Subdomain",
			name: "subdomain",
			value: ""
		}
	];

	this.ggaOpts = [
		{
			title:"Host",
			name: "host",
			value: ""
		},
		{
			title:"Port",
			name: "port",
			value: 2003
		}
	]

	this.emailOpts = [
		{
			title: "To",
			name: "recipients",
			value: "",
			format: function() {
				if (typeof this.value == "string") {
					this.value = this.value.split(",");
				}
			}
		},
		{
			title: "From",
			name:"sender",
			value:""
		},
		{
			title: "User",
			name:"user",
			value:""
		},
		{
			title: "Password",
			name:"password",
			value:""
		},
		{
			title:"Host",
			name:"host",
			value: "smtp.gmail.com"
		},
		{
			title:"Port",
			name:"port",
			value: 465
		}
	];
	this.consoleOpts = [];
	this.chips = [];

	this.getOpts = function(type) {
		switch (type) {
			case "pager_duty":
				return this.pdOpts;

			case "email":
				return this.emailOpts;

			case "console":
				return this.consoleOpts;

			case "grafana_graphite_annotation":
				return this.ggaOpts;

			default:
				return [];
		}
	}



	this.newEscalation = function() {

		if (!this.type) {
			return;
		}

		var e = {
			type: this.type
		};

		var opts = this.getOpts(this.type);
		for (var i = 0; i < opts.length; i++) {

			// if the opts value has a format function, all it
			if (opts[i].format) {
				opts[i].format()
			}
			e[opts[i].name] = opts[i].value;
		}

		this.chips.push(e);
	}

	this.submitNew = function() {
		if (!this.name) {
			return;
		}
		$scope.newEscalationProgress = 50;
		a = this;
		$http.post("api/escalation/config/" + this.name, this.chips).success(function(data) {
			a.reset();
		});

	}

	this.reset = function() {
		this.type = null;
		this.name = "";
		this.chips = [];
		this.opts = {};
		$scope.newEscalationProgress = 0;
	}
}


angular.module("bangarang").controller("NewEscalationController", NewEscalationController);
