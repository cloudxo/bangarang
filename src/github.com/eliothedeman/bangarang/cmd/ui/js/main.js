function Router($scope, $cookies, $http) {
	this.selected = 0;
    this.auth_token
	this.getSelected = function() {
		var s = $cookies.get("router:tab");
		if (s) {
			this.selected = s;
		}

		return this.selected;
	}

	this.updateSelected = function(index) {
		$cookies.put("router:tab", index); 
		this.selected = index;
	}

    get_auth_token = function() {
        return $cookies.get("session:token");
    }

    set_auth_token = function(token) {
        $cookies.put("session:token", token);
        $http.defaults.headers.common["BANG_SESSION"] = token;
    }

    delete_auth_token = function() {
        $cookies.remove("session:token");
    }

    this.login = function(username, password) {
        $http.get("api/auth/user?user="+username+"&pass="+password).then(function(response){
            set_auth_token(response.data.token);
            $cookies.put("session:logged_in", true);

            // refresh
            document.location.reload();

        }, function(response){
            alert("Invalid username/password")
        });
    }

    this.logout = function() {
        delete_auth_token();

        // refresh
        document.location.reload();
    }

    init = function() {
        if (get_auth_token()) {
            $scope.logged_in = true;
            session = get_auth_token();
            if (session) {
                set_auth_token(session);
            }

            $scope.logged_in = true;
        } else {
            $scope.logged_in = false;
        }
    }

    init()
}

angular.module("bangarang").controller("Router", Router);

function Config($scope, $cookies, $http, $mdDialog) {
	$scope.snapshots = [];
	$scope.snapshotsByHash = {};
	$scope.fetchSnapshots = function() {
		$http.get("api/config/version/*").success(function(data) {
			data.sort(function(x,y){
				return new Date(x.time_stamp) - new Date(y.time_stamp)

			});

			data.reverse()
			$scope.snapshotsByHash = {};
			for (var i = data.length - 1; i >= 0; i--) {
				data[i].date = new Date(data[i].time_stamp).format("h:M:s mmmm-dd- yyyy")
				$scope.snapshotsByHash[data[i].hash] = data[i]
			};
			$scope.snapshots = data;
		});
	}
	$scope.showRevert = function(hash) {
		var confirm = $mdDialog.confirm()
		.title("Revert Config")
		.content("Are you sure you want to revert to config version: " + hash)
		.ariaLabel('Alert Dialog Demo')
		.ok("Yes")
		.cancel("Cancel")

		$mdDialog.show(confirm).then(function(){
			$scope.updateCurrent(hash);
		}, function(){

		});

	}

	$scope.updateCurrent = function(hash) {
		$http.post("api/config/version/" + hash).success(function(data) {
			$scope.fetchSnapshots();
		})
	}
	this.selected = 0;
	this.getSelected = function() {
		var s = $cookies.get("conf:tab");
		if (s) {
			this.selected = s;
		}
		return this.selected;
	}
	this.updateSelected = function(name) {
		$cookies.put("conf:tab", name);
		this.selected = name;
	}
}

angular.module("bangarang").controller("Config", Config)

/*
 * Date Format 1.2.3
 * (c) 2007-2009 Steven Levithan <stevenlevithan.com>
 * MIT license
 *
 * Includes enhancements by Scott Trenda <scott.trenda.net>
 * and Kris Kowal <cixar.com/~kris.kowal/>
 *
 * Accepts a date, a mask, or a date and a mask.
 * Returns a formatted version of the given date.
 * The date defaults to the current date/time.
 * The mask defaults to dateFormat.masks.default.
 */

var dateFormat = function () {
    var token = /d{1,4}|m{1,4}|yy(?:yy)?|([HhMsTt])\1?|[LloSZ]|"[^"]*"|'[^']*'/g,
        timezone = /\b(?:[PMCEA][SDP]T|(?:Pacific|Mountain|Central|Eastern|Atlantic) (?:Standard|Daylight|Prevailing) Time|(?:GMT|UTC)(?:[-+]\d{4})?)\b/g,
        timezoneClip = /[^-+\dA-Z]/g,
        pad = function (val, len) {
            val = String(val);
            len = len || 2;
            while (val.length < len) val = "0" + val;
            return val;
        };

    // Regexes and supporting functions are cached through closure
    return function (date, mask, utc) {
        var dF = dateFormat;

        // You can't provide utc if you skip other args (use the "UTC:" mask prefix)
        if (arguments.length == 1 && Object.prototype.toString.call(date) == "[object String]" && !/\d/.test(date)) {
            mask = date;
            date = undefined;
        }

        // Passing date through Date applies Date.parse, if necessary
        date = date ? new Date(date) : new Date;
        if (isNaN(date)) throw SyntaxError("invalid date");

        mask = String(dF.masks[mask] || mask || dF.masks["default"]);

        // Allow setting the utc argument via the mask
        if (mask.slice(0, 4) == "UTC:") {
            mask = mask.slice(4);
            utc = true;
        }

        var _ = utc ? "getUTC" : "get",
            d = date[_ + "Date"](),
            D = date[_ + "Day"](),
            m = date[_ + "Month"](),
            y = date[_ + "FullYear"](),
            H = date[_ + "Hours"](),
            M = date[_ + "Minutes"](),
            s = date[_ + "Seconds"](),
            L = date[_ + "Milliseconds"](),
            o = utc ? 0 : date.getTimezoneOffset(),
            flags = {
                d:    d,
                dd:   pad(d),
                ddd:  dF.i18n.dayNames[D],
                dddd: dF.i18n.dayNames[D + 7],
                m:    m + 1,
                mm:   pad(m + 1),
                mmm:  dF.i18n.monthNames[m],
                mmmm: dF.i18n.monthNames[m + 12],
                yy:   String(y).slice(2),
                yyyy: y,
                h:    H % 12 || 12,
                hh:   pad(H % 12 || 12),
                H:    H,
                HH:   pad(H),
                M:    M,
                MM:   pad(M),
                s:    s,
                ss:   pad(s),
                l:    pad(L, 3),
                L:    pad(L > 99 ? Math.round(L / 10) : L),
                t:    H < 12 ? "a"  : "p",
                tt:   H < 12 ? "am" : "pm",
                T:    H < 12 ? "A"  : "P",
                TT:   H < 12 ? "AM" : "PM",
                Z:    utc ? "UTC" : (String(date).match(timezone) || [""]).pop().replace(timezoneClip, ""),
                o:    (o > 0 ? "-" : "+") + pad(Math.floor(Math.abs(o) / 60) * 100 + Math.abs(o) % 60, 4),
                S:    ["th", "st", "nd", "rd"][d % 10 > 3 ? 0 : (d % 100 - d % 10 != 10) * d % 10]
            };

        return mask.replace(token, function ($0) {
            return $0 in flags ? flags[$0] : $0.slice(1, $0.length - 1);
        });
    };
}();

// Some common format strings
dateFormat.masks = {
    "default":      "ddd mmm dd yyyy HH:MM:ss",
    shortDate:      "m/d/yy",
    mediumDate:     "mmm d, yyyy",
    longDate:       "mmmm d, yyyy",
    fullDate:       "dddd, mmmm d, yyyy",
    shortTime:      "h:MM TT",
    mediumTime:     "h:MM:ss TT",
    longTime:       "h:MM:ss TT Z",
    isoDate:        "yyyy-mm-dd",
    isoTime:        "HH:MM:ss",
    isoDateTime:    "yyyy-mm-dd'T'HH:MM:ss",
    isoUtcDateTime: "UTC:yyyy-mm-dd'T'HH:MM:ss'Z'"
};

// Internationalization strings
dateFormat.i18n = {
    dayNames: [
        "Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat",
        "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"
    ],
    monthNames: [
        "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
        "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"
    ]
};

// For convenience...
Date.prototype.format = function (mask, utc) {
    return dateFormat(this, mask, utc);
};