
export function handleLoginError(xhr: XMLHttpRequest): string {
	return xhr.responseText
}

interface LoginResponse {
	accessToken: string
	expires_at: string
}

export function setupTokens(xhr: XMLHttpRequest) {
	if (xhr.status === 200) {
		try {
			const loginResponse: LoginResponse = JSON.parse(xhr.responseText)
			const accessToken = loginResponse.accessToken
			const expiresAt = loginResponse.expires_at

			sessionStorage.setItem('accessToken', accessToken)
			sessionStorage.setItem('expiresAt', expiresAt)
			const redirectEvent = new CustomEvent('login-successful')
			document.body.dispatchEvent(redirectEvent)
			console.log("event dispatched")
		} catch(error) {
			// login failed
		}
	} else {
		return
	}
}

export function configureRedirect(event) {
	const accessToken = sessionStorage.getItem("accessToken")
	event.detail.headers["Authorization"] = `Bearer ${accessToken}`
}

