import Alpine from "alpinejs"
import focus from '@alpinejs/focus'
import htmx from "htmx.org"
import "./main.css"
import { createIcons, Search, UserRound, LogOut, Wrench, ChevronDown, X, Plus, Ellipsis, Pencil, FolderPlus, FolderDot, FolderOpen, Trash2, ChevronRight, FileText, FilePlus2, UserRoundPlus, Eye } from "lucide"

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
		ChevronRight,
		X,
		Plus,
		Ellipsis,
		Pencil,
		FolderPlus,
		FolderDot,
		FolderOpen,
		Trash2,	
		FileText,
		FilePlus2,
		UserRoundPlus,
		Eye
	}})
}

loadIcons()
