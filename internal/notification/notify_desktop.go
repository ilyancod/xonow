package notification

import (
	"github.com/gen2brain/beeep"
)

type NotifyDesktop struct {
	IconPath string
}

func (nd *NotifyDesktop) notify(title, text string) error {
	return beeep.Notify(title, text, nd.IconPath)
}
