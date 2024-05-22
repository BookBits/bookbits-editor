import { FileAutoSave } from "./types/FileAutoSave"
import { SessionResponse } from "./types/SessionResponse"

interface CsrfTokenResponse {
	csrfToken: string
}

export async function getCSRFToken(): Promise<string> {
	let csrfToken: string|undefined = await fetch("/csrf").then((res) => {
		if (res.status === 200) {
			return res.json().then((body: CsrfTokenResponse) => {
				return body.csrfToken
			})
		}
	})
	return csrfToken!
}

async function refreshTokens() {
	await fetch("/refresh", {
		method: 'POST',
		headers: {
			"X-CSRF-Token": await getCSRFToken()
		}
	}).then((res) => {
		if (res.status === 200) {
			res.json().then((contents: SessionResponse) => {
				sessionStorage.setItem('expiresAt', contents.expires_at)
				setupSessionRefresh()
			})
		}
	})
}

function setupSessionRefresh() {
	const expiresAtVal: string = sessionStorage.getItem("expiresAt")!
	const expiresAt = Date.parse(expiresAtVal)
	const bufferTime = 60 * 1000

	const timeOut = expiresAt - Date.now() - bufferTime

	setTimeout(refreshTokens, timeOut)
}

setupSessionRefresh()

export function logout() {
	sessionStorage.clear()
}

export function saveFile(evt: any, content: string, fileID: string) {
	console.log('saving to server')
	localStorage.removeItem(`autosave:${fileID}`)
	evt.detail.parameters['content'] = content
}

function getCookie(cname: string) : string {
  const name = cname + "=";
  const decodedCookie = decodeURIComponent(document.cookie);
  const ca = decodedCookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) === 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

async function refreshFileLock(fileID: string) {
	var url = "/app/projects/files/" + fileID + "/lock"
	await fetch(url, {
		method: 'POST',
		headers: {
			"X-CSRF-Token": await getCSRFToken()
		}
	}).then((res) => {
		if (res.status === 200) {
			setupFileLockRefresh(fileID)
		}
	})
}

export function setupFileLockRefresh(fileID: string) {
	const expiresAt: string = getCookie("File-Lock-Expire")
	const expires = Date.parse(expiresAt)

	const buffer = 60 * 1000
	const timeOut = expires - Date.now() - buffer

	setTimeout(async () => {await refreshFileLock(fileID)}, timeOut)
}

export function autoSave(fileID: string, fileVersion: number, fileContent: string) {
	const autoSaveKey = `autosave:${fileID}`
	var autoSave: FileAutoSave = {
		fileID,
		fileVersion,
		fileContent
	}
	localStorage.setItem(autoSaveKey, JSON.stringify(autoSave))
}

export function unsavedChanges(fileID: string) : boolean {
	const changes = localStorage.getItem(`autosave:${fileID}`)
	if (changes) {
		return true
	}
	return false
}

export function getUnsavedChanges(fileID: string) : string {
	const changes = localStorage.getItem(`autosave:${fileID}`) ?? ""
	if (changes === "") {
		return ""
	}
	const changesContent: FileAutoSave = JSON.parse(changes)
	return changesContent.fileContent
}

export function discardUnsavedChanges(fileID: string) {
	localStorage.removeItem(`autosave:${fileID}`)
}

export async function unlockFile(fileID: string, event: Event) {
	var url = "/app/projects/files/" + fileID + "/unlock"
	await fetch(url, {
		method: 'POST',
		headers: {
			"X-CSRF-Token": await getCSRFToken()
		}
	}).then((res) => {
		if (res.status !== 200) {
			window.dispatchEvent(new Event("editor:unlock-error"))
			event.preventDefault()
		}
	})

}
