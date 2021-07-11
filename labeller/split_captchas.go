package labeller

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/attilaolah/cad-rs/proto"
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

				for _, part := range cut4(img) {
					ch <- part
				}
			}
		}
	}()

	return ch, errs
}

// Extract single letters from a 4-letter captcha.
func cut4(img image.Image) [4]image.Image {
	imgs := [4]image.Image{}

	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()

	for i := 0; i < 4; i++ {
		cut := image.NewRGBA(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: (dx - crop4l - crop4r) / 4,
				Y: dy - crop4t - crop4b,
			},
		})
		ref := image.Point{
			X: ((dx-crop4l-crop4r)/4)*i + crop4l,
			Y: crop4t,
		}
		draw.Over.Draw(cut, cut.Bounds(), img, ref)
		imgs[i] = cut
	}

	return imgs
}
