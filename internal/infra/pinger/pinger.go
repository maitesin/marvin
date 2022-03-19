package pinger

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Pinger struct {
	address   string
	frequency int
}

func NewPinger(address string, frequency int) Pinger {
	return Pinger{
		address:   address,
		frequency: frequency,
	}
}

func (p Pinger) Start(ctx context.Context) {
	t := time.NewTicker(time.Duration(p.frequency) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			p.ping()
		}
	}
}

func (p Pinger) ping() {
	_, err := http.Get(p.address)
	if err != nil {
		fmt.Println(err)
	}
}
