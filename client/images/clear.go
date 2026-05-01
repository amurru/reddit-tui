package images

import termimg "github.com/blacktop/go-termimg"

func ClearTerminalImages() error {
	return termimg.ClearAll()
}
