package middleware

import "io"

var methods = map[string]func(io.ReadWriter, string) io.ReadWriter{
	"aes-256-cfb":             NewAes256Cfb,
	"aes-128-gcm":             NewAes128Gcm,
	"aes-256-gcm":             NewAes256Gcm,
	"chacha20-ietf-poly1305":  NewChacha20Poly1305,
	"xchacha20-ietf-poly1305": NewXChacha20Poly1305,
	"noop":                    NewNoop,
}

func NotExist(method string) bool {
	_, ok := methods[method]
	return !ok
}

func Get(method string) func(io.ReadWriter, string) io.ReadWriter {
	fn, ok := methods[method]
	if ok {
		return fn
	}

	panic(method + " not support")
}
