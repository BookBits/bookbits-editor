import Alpine from "alpinejs"
import htmx from "htmx.org"
import "./main.css"

window.Alpine = Alpine
window.htmx = htmx

Alpine.start()

export function successfulRefresh(xhr) {
	if (xhr.status === 200) {
		try {
			const tokensResponse = JSON.parse(xhr.responseText)
			const expiresAt = tokensResponse.expires_at

			sessionStorage.setItem('expiresAt', expiresAt)
			const redirectEvent = new CustomEvent('tokens-setup')
			document.body.dispatchEvent(redirectEvent)
		} catch(error) {
			// parse failed
		}
	}
}
