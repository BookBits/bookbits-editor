var appBundle = (function (exports) {
    'use strict';

    /******************************************************************************
    Copyright (c) Microsoft Corporation.

    Permission to use, copy, modify, and/or distribute this software for any
    purpose with or without fee is hereby granted.

    THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
    REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
    AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
    INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
    LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
    OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
    PERFORMANCE OF THIS SOFTWARE.
    ***************************************************************************** */
    /* global Reflect, Promise, SuppressedError, Symbol */


    function __awaiter(thisArg, _arguments, P, generator) {
        function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
        return new (P || (P = Promise))(function (resolve, reject) {
            function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
            function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
            function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
            step((generator = generator.apply(thisArg, _arguments || [])).next());
        });
    }

    function __generator(thisArg, body) {
        var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
        return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
        function verb(n) { return function (v) { return step([n, v]); }; }
        function step(op) {
            if (f) throw new TypeError("Generator is already executing.");
            while (g && (g = 0, op[0] && (_ = 0)), _) try {
                if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
                if (y = 0, t) op = [op[0] & 2, t.value];
                switch (op[0]) {
                    case 0: case 1: t = op; break;
                    case 4: _.label++; return { value: op[1], done: false };
                    case 5: _.label++; y = op[1]; op = [0]; continue;
                    case 7: op = _.ops.pop(); _.trys.pop(); continue;
                    default:
                        if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                        if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                        if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                        if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                        if (t[2]) _.ops.pop();
                        _.trys.pop(); continue;
                }
                op = body.call(thisArg, _);
            } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
            if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
        }
    }

    typeof SuppressedError === "function" ? SuppressedError : function (error, suppressed, message) {
        var e = new Error(message);
        return e.name = "SuppressedError", e.error = error, e.suppressed = suppressed, e;
    };

    function getCSRFToken() {
        return __awaiter(this, void 0, void 0, function () {
            var csrfToken;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, fetch("/csrf").then(function (res) {
                            if (res.status === 200) {
                                return res.json().then(function (body) {
                                    return body.csrfToken;
                                });
                            }
                        })];
                    case 1:
                        csrfToken = _a.sent();
                        return [2 /*return*/, csrfToken];
                }
            });
        });
    }
    function refreshTokens() {
        return __awaiter(this, void 0, void 0, function () {
            var _a, _b, _c;
            var _d, _e;
            return __generator(this, function (_f) {
                switch (_f.label) {
                    case 0:
                        _a = fetch;
                        _b = ["/refresh"];
                        _d = {
                            method: 'POST'
                        };
                        _e = {};
                        _c = "X-CSRF-Token";
                        return [4 /*yield*/, getCSRFToken()];
                    case 1: return [4 /*yield*/, _a.apply(void 0, _b.concat([(_d.headers = (_e[_c] = _f.sent(),
                                _e),
                                _d)])).then(function (res) {
                            if (res.status === 200) {
                                res.json().then(function (contents) {
                                    sessionStorage.setItem('expiresAt', contents.expires_at);
                                    setupSessionRefresh();
                                });
                            }
                        })];
                    case 2:
                        _f.sent();
                        return [2 /*return*/];
                }
            });
        });
    }
    function setupSessionRefresh() {
        var expiresAtVal = sessionStorage.getItem("expiresAt");
        var expiresAt = Date.parse(expiresAtVal);
        var bufferTime = 60 * 1000;
        var timeOut = expiresAt - Date.now() - bufferTime;
        setTimeout(refreshTokens, timeOut);
    }
    setupSessionRefresh();
    function logout() {
        sessionStorage.clear();
    }
    function saveFile(evt, content, fileID) {
        console.log('saving to server');
        localStorage.removeItem("autosave:".concat(fileID));
        evt.detail.parameters['content'] = content;
    }
    function getCookie(cname) {
        var name = cname + "=";
        var decodedCookie = decodeURIComponent(document.cookie);
        var ca = decodedCookie.split(';');
        for (var i = 0; i < ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0) === ' ') {
                c = c.substring(1);
            }
            if (c.indexOf(name) === 0) {
                return c.substring(name.length, c.length);
            }
        }
        return "";
    }
    function refreshFileLock(fileID) {
        return __awaiter(this, void 0, void 0, function () {
            var url, _a, _b, _c;
            var _d, _e;
            return __generator(this, function (_f) {
                switch (_f.label) {
                    case 0:
                        url = "/app/projects/files/" + fileID + "/lock";
                        _a = fetch;
                        _b = [url];
                        _d = {
                            method: 'POST'
                        };
                        _e = {};
                        _c = "X-CSRF-Token";
                        return [4 /*yield*/, getCSRFToken()];
                    case 1: return [4 /*yield*/, _a.apply(void 0, _b.concat([(_d.headers = (_e[_c] = _f.sent(),
                                _e),
                                _d)])).then(function (res) {
                            if (res.status === 200) {
                                setupFileLockRefresh(fileID);
                            }
                        })];
                    case 2:
                        _f.sent();
                        return [2 /*return*/];
                }
            });
        });
    }
    function setupFileLockRefresh(fileID) {
        var _this = this;
        var expiresAt = getCookie("File-Lock-Expire");
        var expires = Date.parse(expiresAt);
        var buffer = 60 * 1000;
        var timeOut = expires - Date.now() - buffer;
        setTimeout(function () { return __awaiter(_this, void 0, void 0, function () { return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, refreshFileLock(fileID)];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        }); }); }, timeOut);
    }
    function autoSave(fileID, fileVersion, fileContent) {
        var autoSaveKey = "autosave:".concat(fileID);
        var autoSave = {
            fileID: fileID,
            fileVersion: fileVersion,
            fileContent: fileContent
        };
        localStorage.setItem(autoSaveKey, JSON.stringify(autoSave));
    }
    function unsavedChanges(fileID) {
        var changes = localStorage.getItem("autosave:".concat(fileID));
        if (changes) {
            return true;
        }
        return false;
    }
    function getUnsavedChanges(fileID) {
        var _a;
        var changes = (_a = localStorage.getItem("autosave:".concat(fileID))) !== null && _a !== void 0 ? _a : "";
        if (changes === "") {
            return "";
        }
        var changesContent = JSON.parse(changes);
        return changesContent.fileContent;
    }
    function discardUnsavedChanges(fileID) {
        localStorage.removeItem("autosave:".concat(fileID));
    }
    function unlockFile(fileID, event) {
        return __awaiter(this, void 0, void 0, function () {
            var url, _a, _b, _c;
            var _d, _e;
            return __generator(this, function (_f) {
                switch (_f.label) {
                    case 0:
                        url = "/app/projects/files/" + fileID + "/unlock";
                        _a = fetch;
                        _b = [url];
                        _d = {
                            method: 'POST'
                        };
                        _e = {};
                        _c = "X-CSRF-Token";
                        return [4 /*yield*/, getCSRFToken()];
                    case 1: return [4 /*yield*/, _a.apply(void 0, _b.concat([(_d.headers = (_e[_c] = _f.sent(),
                                _e),
                                _d)])).then(function (res) {
                            if (res.status !== 200) {
                                window.dispatchEvent(new Event("editor:unlock-error"));
                                event.preventDefault();
                            }
                        })];
                    case 2:
                        _f.sent();
                        return [2 /*return*/];
                }
            });
        });
    }

    exports.autoSave = autoSave;
    exports.discardUnsavedChanges = discardUnsavedChanges;
    exports.getCSRFToken = getCSRFToken;
    exports.getUnsavedChanges = getUnsavedChanges;
    exports.logout = logout;
    exports.saveFile = saveFile;
    exports.setupFileLockRefresh = setupFileLockRefresh;
    exports.unlockFile = unlockFile;
    exports.unsavedChanges = unsavedChanges;

    return exports;

})({});
