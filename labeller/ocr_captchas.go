package labeller

import (
	"github.com/otiai10/gosseract"
)

type OCR struct {
	*gosseract.Client
}

func OCR1(tessdir string) (*OCR, error) {
	c := gosseract.NewClient()

	if err := c.SetTessdataPrefix(tessdir); err != nil {
		c.Close()
		return nil, err
	}

	if err := c.SetConfigFile("foo/bar"); err != nil {
		c.Close()
		return nil, err
	}

	return &OCR{c}, nil
}
