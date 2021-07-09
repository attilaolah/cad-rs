package labeller

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/attilaolah/ekat/proto"
)

// Pre-crop 4-letter captchas, before slicing them up:
const (
	crop4t = 5  // crop top
	crop4b = 5  // crop bottom
	crop4l = 5  // crop left
	crop4r = 15 // crop right
)

// Split4Captchas generates single-letter segments from 4-letter captchas.
func Split4Captchas(datadir string) (ch chan image.Image, errs chan error) {
	ch = make(chan image.Image)
	errs = make(chan error)

	done := func() {
		close(errs)
		close(ch)
	}
	fail := func(err error) {
		errs <- err
		done()
	}

	pat := filepath.Join(datadir, "????????-????-????-????-????????????.json")
	fs, err := filepath.Glob(pat)
	if err != nil {
		go fail(fmt.Errorf("failed to match pattern %q: %w", pat, err))
		return
	}

	go func() {
		defer done()

		for _, fn := range fs {
			f, err := os.Open(fn)
			if err != nil {
				errs <- fmt.Errorf("failed to open %q: %w", fn, err)
				continue
			}
			closef := func() {
				if err = f.Close(); err != nil {
					errs <- fmt.Errorf("failed to close %q: %w", f.Name(), err)
				}
			}

			c := proto.Captcha{}
			if err = json.NewDecoder(f).Decode(&c); err != nil {
				errs <- fmt.Errorf("failed to decode %q: %w", f.Name(), err)
				closef()
				continue
			}
			closef()

			if c.Type != proto.Captcha_ALPHANUM_4 {
				continue // ignore other types
			}

			for _, s := range c.Samples {
				fn = filepath.Join(datadir, "samples", fmt.Sprintf("%s.jpg", s.Sha1))
				f, err = os.Open(fn)
				if err != nil {
					errs <- fmt.Errorf("failed to open %q: %w", fn, err)
					continue
				}

				img, err := jpeg.Decode(f)
				if err != nil {
					errs <- fmt.Errorf("failed to decode %q: %w", f.Name(), err)
					closef()
					continue
				}
				closef()

				for part := range cut4(img) {
					ch <- part
				}
			}
		}
	}()

	return ch, errs
}

// Cut a 4-letter captcha in four.
func cut4(img image.Image) <-chan image.Image {
	ch := make(chan image.Image)

	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	go func() {
		defer close(ch)
		for x := crop4l; x < dx-crop4l-crop4r; x += (dx - crop4l - crop4r) / 4 {
			cut := image.NewRGBA(image.Rectangle{
				Min: image.Point{},
				Max: image.Point{(dx - crop4l - crop4r) / 4, dy - crop4t - crop4b},
			})
			draw.Over.Draw(cut, cut.Bounds(), img, image.Point{X: x, Y: crop4t})
			ch <- cut
		}
	}()

	return ch
}
