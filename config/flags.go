package config

import (
	"strings"
)

type flagList []string

func (l *flagList) String() string {
	return strings.Join(*l, ", ")
}

func (l *flagList) Set(value string) error {
	*l = append(*l, value)
	return nil
}
