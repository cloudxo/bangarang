function DashboardController($scope, $http) {
	$scope.incidents = [];
	$scope.fetching = false;
	var stopper = null;

	$scope.startFetching = function() {
		$scope.fetchIncidents();

		if (!$scope.fetching) {
			$scope.fetching = true;
			stopper = setInterval(function(){
				$scope.fetchIncidents()
			}, 5000)
		}
	}

	$scope.stopFetching = function() {
		clearInterval(stopper);
		$scope.fetching = false;
	}


	$scope.resolveIncident = function(id) {
		$http.delete("api/incident/" + id, null).success(function(){
			$scope.fetchIncidents();
		});
	}

	$scope.formatDescription = function(incident) {
		return incident.service + (incident.sub_service ? "." + incident.sub_service : " ") + " on " + incident.host + " is " + incident.metric.toFixed(2) + " at " + new Date(incident.time * 1000).format("h:M:sTT mmmm-dd-yyyy"); 
	}

	var codes = {
		"0": "OK",
		"1": "WARNING",
		"2": "CRITICAL"
	}

	$scope.getStatusCode = function(status) {
		return codes[status];
	}

	var colors = {
		"0": "green",
		"1": "#FFFD82",
		"2": "#FB5C5C"
	}

	$scope.getStatusColor = function(status) {
		return colors[status]
	}

	$scope.fetchIncidents = function() {
		$http.get("api/incident/*").success(function(data) {
			var ins = [];
			for (k in data) {
				ins.push({key:k, val:data[k]})
			}

			ins.sort(function(x,y) {
				if (x.val.status != y.val.status) {
					return y.val.status - x.val.status
				}
				return y.val.time - x.val.time;
			})
			$scope.incidents = ins;
		});
	}
}

angular.module("bangarang").controller("DashboardController", DashboardController);
