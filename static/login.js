export function handleLoginError(xhr) {
	return xhr.responseText
}

export function setupTokens(xhr) {
	if (xhr.status === 200) {
		try {
			const loginResponse = JSON.parse(xhr.responseText)
			const expiresAt = loginResponse.expires_at

			sessionStorage.setItem('expiresAt', expiresAt)
			const redirectEvent = new CustomEvent('login-successful')
			document.body.dispatchEvent(redirectEvent)
		} catch(error) {
			// parse failed
		}
	} else {
		return
	}
}
