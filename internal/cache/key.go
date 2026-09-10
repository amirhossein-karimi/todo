package cache

import (
	"fmt"
	"strings"
	"time"
)

const (
	TodosKey = "todos"
	TTL      = 5 * time.Minute
)

func CreateKeyName(params ...interface{}) string {
	var name string

	for _, param := range params {
		name += fmt.Sprint(param) + ":"
	}

	return strings.TrimSuffix(name, ":")
}
