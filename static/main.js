export function setupSession(xhr) {
	if (xhr.status === 200) {
		try {
			const tokensResponse = JSON.parse(xhr.responseText)
			const expiresAt = tokensResponse.expires_at

			sessionStorage.setItem('expiresAt', expiresAt)
			const redirectEvent = new CustomEvent('session-setup')
			document.body.dispatchEvent(redirectEvent)
		} catch(error) {
			// parse failed
		}
	}
}

export function handleLoginError(xhr) {
	return xhr.responseText
}
