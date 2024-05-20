package util

import "time"

type ToxOptions struct {
	Password           string
	Tunnel             string
	LocalAddress       string
	CertFile           string
	KeyFile            string
	ConnectTimeout     time.Duration
	Timeout            time.Duration
	InsecureSkipVerify bool
}

func ToToxOptions(opts []ToxOption) ToxOptions {
	toxOptions := ToxOptions{}

	for _, opt := range opts {
		opt.apply(&toxOptions)
	}

	return toxOptions
}

type ToxOption interface {
	apply(*ToxOptions)
}

type funcToxOption struct {
	f func(*ToxOptions)
}

func (fdo funcToxOption) apply(do *ToxOptions) {
	fdo.f(do)
}

func newToxOptionFunc(f func(*ToxOptions)) *funcToxOption {
	return &funcToxOption{
		f: f,
	}
}

func WithPassword(password string) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.Password = HashedPassword(password)
	})
}

func WithTunnel(tunnel string) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.Tunnel = tunnel
	})
}

func WithLocalAddress(localAddress string) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.LocalAddress = localAddress
	})
}

func WithCertFile(certFile string) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.CertFile = certFile
	})
}

func WithKeyFile(keyFile string) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.KeyFile = keyFile
	})
}

func WithConnectTimeout(timeout time.Duration) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.ConnectTimeout = timeout
	})
}

func WithTimeout(timeout time.Duration) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.Timeout = timeout
	})
}

func WithInsecureSkipVerify(insecureSkipVerify bool) ToxOption {
	return newToxOptionFunc(func(o *ToxOptions) {
		o.InsecureSkipVerify = insecureSkipVerify
	})
}
