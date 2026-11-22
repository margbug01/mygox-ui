//go:build !(js && wasm)

package main

func initialHash() string       { return "" }
func initialDark() bool         { return false }
func syncHash(string)           {}
func notifyWebReady()           {}
func syncWebTheme(bool)         {}
func watchHashChange(func(string)) {}
func webViewportWidth() float32    { return 0 }
