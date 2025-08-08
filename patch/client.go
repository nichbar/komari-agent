package patch

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

var (
	Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * 2000,
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}

	Dialer = &net.Dialer{Resolver: Resolver}

	Transport = &http.Transport{
		DialContext:       Dialer.DialContext,
		DisableKeepAlives: false,
	}

	Client = &http.Client{
		Transport: Transport,
		Timeout:   30 * time.Second,
	}
)
