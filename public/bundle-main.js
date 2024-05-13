var mainBundle = (function (exports) {
    'use strict';

    function setupSession(xhr) {
        if (xhr.status === 200) {
            try {
                var tokensResponse = JSON.parse(xhr.responseText);
                var expiresAt = tokensResponse.expires_at;
                sessionStorage.setItem('expiresAt', expiresAt);
                var redirectEvent = new CustomEvent('session-setup');
                document.body.dispatchEvent(redirectEvent);
            }
            catch (error) {
                // parse failed
            }
        }
    }

    exports.setupSession = setupSession;

    return exports;

})({});
