import Alpine from "alpinejs"
import htmx from "htmx.org"
import "./main.css"
import { createIcons, Search, UserRound, LogOut, Wrench } from "lucide"

window.Alpine = Alpine
window.htmx = htmx

Alpine.start()
createIcons({icons: {
	Search,
	UserRound,
	Wrench,
	LogOut
}})
