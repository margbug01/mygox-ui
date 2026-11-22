//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
)

func initialHash() string {
	loc := js.Global().Get("location")
	if !loc.Truthy() {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(loc.Get("hash").String(), "#"))
}

func initialDark() bool {
	doc := js.Global().Get("document")
	if doc.Truthy() {
		el := doc.Get("documentElement")
		if el.Truthy() && el.Call("getAttribute", "data-mygo-theme").String() == "dark" {
			return true
		}
	}
	return false
}

func syncHash(hash string) {
	hist := js.Global().Get("history")
	loc := js.Global().Get("location")
	if !hist.Truthy() || !loc.Truthy() {
		return
	}
	base := loc.Get("pathname").String() + loc.Get("search").String()
	if hash != "" {
		hist.Call("replaceState", nil, "", base+"#"+hash)
	} else {
		hist.Call("replaceState", nil, "", base)
	}
}

func notifyWebReady() {
	fn := js.Global().Get("__mygoxReady")
	if fn.Type() == js.TypeFunction {
		fn.Invoke()
	}
}

func syncWebTheme(dark bool) {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}
	el := doc.Get("documentElement")
	if el.Truthy() {
		if dark {
			el.Call("setAttribute", "data-mygo-theme", "dark")
		} else {
			el.Call("setAttribute", "data-mygo-theme", "light")
		}
	}
}

func watchHashChange(onChange func(string)) {
	win := js.Global()
	if !win.Truthy() {
		return
	}
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		onChange(initialHash())
		evCtor := win.Get("Event")
		if evCtor.Truthy() {
			win.Call("dispatchEvent", evCtor.New("resize"))
		}
		return nil
	})
	win.Call("addEventListener", "hashchange", cb)
}

func webViewportWidth() float32 {
	w := js.Global().Get("innerWidth")
	if w.Type() == js.TypeNumber {
		return float32(w.Float())
	}
	return 0
}
