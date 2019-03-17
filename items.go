package main

import (
	"fmt"
)

type item interface {
	String() string
}

type proxy struct {
	Scheme string
	IP     string
	Port   int
}

func (p proxy) String() string {
	return fmt.Sprintf("%s://%s:%d", p.Scheme, p.IP, p.Port)
}
