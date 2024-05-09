var appBundle = (function (exports) {
	'use strict';

	function refreshTokens() {
		fetch("/refresh", {
			method: 'POST'
		}).then((res) => {
			if (res.status === 200) {
				res.json().then((contents) => {
					sessionStorage.setItem('expiresAt', contents.expires_at);
					setupSessionRefresh();
				});
			}
		});
	}

	function setupSessionRefresh() {
		const expiresAtVal = sessionStorage.getItem("expiresAt");
		const expiresAt = Date.parse(expiresAtVal);
		const bufferTime = 60 * 1000;

		const timeOut = expiresAt - Date.now() - bufferTime;

		setTimeout(refreshTokens, timeOut);
	}

	setupSessionRefresh();

	function logout() {
		sessionStorage.clear();
	}

	exports.logout = logout;

	return exports;

})({});
