package asmr_one

import "github.com/go-kid/ioc"

func init() {
	ioc.Register(NewClient())
}
