package main

import (
	"github.com/gen2brain/beeep"
)

func notify(title, text string) error {
	return beeep.Notify(title, text, "assets/xonotic.png")
}
