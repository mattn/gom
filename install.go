package main

import (
	"fmt"
)

func install(goms []Gom) error {
	for _, gom := range goms {
		fmt.Printf("installing %s(tag: %s, options: %s)\n",
			gom.name,
			gom.tag,
			gom.options)
	}
	return nil
}
