import Alpine from "alpinejs"
import htmx from "htmx.org"
import "./main.css"
import { createIcons, Search, UserRound, LogOut, Wrench, ChevronDown, X, Plus, Ellipsis } from "lucide"

window.Alpine = Alpine
window.htmx = htmx

Alpine.start()

export function loadIcons() {
	createIcons({icons: {
		Search,
		UserRound,
		Wrench,
		LogOut,
		ChevronDown,
		X,
		Plus,
		Ellipsis
	}})
}

loadIcons()
