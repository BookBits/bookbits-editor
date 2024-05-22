import Alpine from "alpinejs"
import focus from '@alpinejs/focus'
import htmx from "htmx.org"
import "./main.css"
import { createIcons, Search, UserRound, LogOut, Wrench, ChevronDown, X, Plus, Ellipsis, Pencil, FolderPlus, FolderDot, FolderOpen, Trash2, ChevronRight, FileText, FilePlus2, UserRoundPlus, Eye, Heading1, Heading2, Bold, Italic, Underline, List, ListOrdered, Pilcrow, Quote, ChevronLeft } from "lucide"
import { Editor } from '@tiptap/core'
import StarterKit from '@tiptap/starter-kit'
import { Underline as tiptapUnderline} from '@tiptap/extension-underline';

document.addEventListener('alpine:init', () => {
	Alpine.data('editor', (content) => {
		let editor;

		return {
			updatedAt: Date.now(),
			lastAutoSave: Date.now(),
			lastAutoSync: Date.now(),
			init() {
				const _this = this

				editor = new Editor({
					element: this.$refs.element,
					extensions: [StarterKit, tiptapUnderline],
					content: content,
					onCreate({}) {
						_this.updatedAt = Date.now()
					},
					onUpdate({}) {
						const autoSaveBuffer = 2 * 60 * 1000
						const autoSyncBuffer = 10 * 60 * 1000
						const now = Date.now()
						_this.updatedAt = now
						if ((now - _this.lastAutoSave) > autoSaveBuffer) {
							window.dispatchEvent(new Event('editor-auto-save'))
							_this.lastAutoSave = now
						}
						if ((now - _this.lastAutoSync) > autoSyncBuffer) {
							window.dispatchEvent(new Event('editor-auto-sync'))
							_this.lastAutoSync = now
						}
					},
					onSelectionUpdate({}) {
						_this.updatedAt = Date.now()
					}
				})
			},
			isLoaded() {
				return editor
			},
			isActive(type, opts = {}) {
      		  return editor.isActive(type, opts)
      		},
      		toggleHeading(opts) {
      		  editor.chain().toggleHeading(opts).focus().run()
      		},
      		toggleBold() {
      		  editor.chain().toggleBold().focus().run()
      		},
      		toggleItalic() {
      		  editor.chain().toggleItalic().focus().run()
      		},
		toggleUnderline() {
			editor.commands.toggleUnderline()
		},
		toggleBulletList() {
			editor.chain().toggleBulletList().focus().run()
		},
		toggleOrderedList() {
			editor.chain().toggleOrderedList().focus().run()
		},
		setParagraph() {
			editor.chain().setParagraph().focus().run()
		},
		toggleQuote() {
			editor.chain().toggleBlockquote().focus().run()
		},
		getContent() {
			return editor.getHTML()
		}
		}
	})
})

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
		ChevronLeft,
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
		Eye,
		Heading1,
		Heading2,
		Bold,
		Italic,
		Underline,
		List,
		ListOrdered,
		Pilcrow,
		Quote
	}})
}

loadIcons()
