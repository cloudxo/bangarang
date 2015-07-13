function NewPolicyController($scope, $http) {
	this.np = {};
	this.compOps = ["greater", "less", "exactly"];
	this.escalations = [];

	this.fetchEscalations = function() {
		$http.get("api/escalation/config/*").success(function(data, status) {
			if (data != null) {
				this.escalations = data;
			}
		}).error(function(data, status) {
			this.escalations = ["No Escalations Available"]
		});
	}

	this.createPolicyStruct = function() {
		if (!this.np.name) {
			return null;
		}

		var p = {
			name: this.np.name
		};

		// set up match
		if (this.matchChips.length > 0) {
			p.match = {};
			for (var i = 0; i < this.matchChips.length; i++) {
				p.match[this.matchChips[i].key] = this.matchChips[i].val;
			}
		}
		if (this.notMatchChips.length > 0) {
			p.not_match = {
				occurences: this.wOcc
			};
			for (var i = 0; i < this.notMatchChips.length; i++) {
				p.not_match[this.notMatchChips[i].key] = this.notMatchChips[i].val;
			}
		}
		if (this.critOpChips.length > 0) {
			p.crit = {
				occurences: this.cOcc
			};
			for (var i = 0; i < this.critOpChips.length; i++) {
				p.crit[this.critOpChips[i].key] = this.critOpChips[i].val;
			}
		}
		if (this.warnOpChips.length > 0) {
			p.warn = {};
			for (var i = 0; i < this.warnOpChips.length; i++) {
				p.warn[this.warnOpChips[i].key] = this.warnOpChips[i].val;
			}
		}
		return p;
	}

	this.addPolicy = function() {
		var pol = this.createPolicyStruct();
		if (pol) {
			$http.post("api/policy/config/" + pol.name, this.createPolicyStruct());
			this.reset();
		}
	}

	this.cancelPolicy = function() {
		this.reset();
	}

	this.addNewMatch = function() {
		if (this.matchChips == null ) {
			this.matchChips = [];
		}
		this.matchChips.push({"key": this.newMatchKey, "val": this.newMatchVal});
		this.newMatchKey = "";
		this.newMatchVal = "";
	}

	this.addNewNotMatch = function() {

		if (this.not_matchChips == null) {
			this.not_matchChips = [];
		}
		this.notMatchChips.push({"key": this.newNotMatchKey, "val": this.newNotMatchVal});
		this.newNotMatchKey = "";
		this.newNotMatchVal = "";
	}

	this.addNewCritOp = function() {
		if (this.cOpKey && this.cOpVal ) {
			this.critOpChips.push({"key": this.cOpKey, "val": this.cOpVal});
			this.cOpKey = "";
			this.cOpVal = "";
		}
	}
	this.addNewWarnOp = function() {
		if (np.wOpKey && np.wOpVal ) {
			this.warnOpChips.push({"key": np.wOpKey, "val": np.wOpVal});
			this.wOpVal = "";
			this.wOpKey = "";
		}
	}



	this.init = function() {
		this.cOpVal = ""
		this.cOpKey = ""
		this.wOpVal = ""
		this.wOpKey = ""
		this.wOcc = 1
		this.cOcc = 1
		this.fetchEscalations();
	}

	this.reset = function() {
		this.init();
		this.matchChips = [];
		this.notMatchChips = [];
		this.critOpChips = [];
		this.warnOpChips = []
	}

	this.reset();
}
angular.module('bangarang').controller("NewPolicyController", NewPolicyController);

function PolicyController($scope, $http) {
	$scope.policies = null;

	this.fetchPolicies = function() {
		$http.get("api/policy/config/*").success(function(data, status) {
			$scope.policies = data;
		});
	}
	this.init = function() {
		this.fetchPolicies();
	}

	this.init();
}

angular.module('bangarang').controller("PolicyController", PolicyController);
