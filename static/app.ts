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
