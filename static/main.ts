import { SessionResponse } from "./types/SessionResponse"

export function setupSession(xhr: XMLHttpRequest) {
	if (xhr.status === 200) {
		try {
			const tokensResponse: SessionResponse = JSON.parse(xhr.responseText)
			const expiresAt = tokensResponse.expires_at

			sessionStorage.setItem('expiresAt', expiresAt)
			const redirectEvent = new CustomEvent('session-setup')
			document.body.dispatchEvent(redirectEvent)
		} catch(error) {
			// parse failed
		}
	}
}
