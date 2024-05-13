import Alpine from "alpinejs"
import focus from '@alpinejs/focus'
import htmx from "htmx.org"
import "./main.css"
import { createIcons, Search, UserRound, LogOut, Wrench, ChevronDown, X, Plus, Ellipsis, Pencil } from "lucide"

window.Alpine = Alpine
window.htmx = htmx

Alpine.plugin(focus)
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
		Ellipsis,
		Pencil
	}})
}

loadIcons()
