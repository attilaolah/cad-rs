package scrapers

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/attilaolah/cad-rs/proto"
)

const (
	eKatFindObj  = eKatURL + "/FindObjekat.aspx?OpstinaID=%d"
	eKatFindAddr = eKatURL + "/FindAdresa.aspx?OpstinaID=%d"
	eKatFindParc = eKatURL + "/FindParcela.aspx?KoID=%d"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Scrape5Captchas scrapes 4-digit captchas.
func Scrape4Captchas(ctx context.Context, municipalities string, samples int) (chan *proto.Captcha, chan error) {
	return ScrapeCaptchas(ctx, proto.Captcha_ALPHANUM_4, municipalities, samples)
}

// Scrape5Captchas scrapes 5-digit captchas.
func Scrape5Captchas(ctx context.Context, municipalities string, samples int) (chan *proto.Captcha, chan error) {
	return ScrapeCaptchas(ctx, proto.Captcha_ALPHANUM_5, municipalities, samples)
}

// ScrapeCaptchas fetches captchas of any type.
// It will keep generating captchas until the context is cancelled.
func ScrapeCaptchas(ctx context.Context, typ proto.Captcha_Type, municipalities string, samples int) (cs chan *proto.Captcha, errs chan error) {
	cs = make(chan *proto.Captcha)
	errs = make(chan error)

	done := func() {
		close(errs)
		close(cs)
	}
	fail := func(err error) {
		errs <- err
		done()
	}

	f, err := os.Open(municipalities)
	if err != nil {
		go fail(fmt.Errorf("failed to open %q: %w", municipalities, err))
		return
	}

	ms := []*proto.Municipality{}
	if err := json.NewDecoder(f).Decode(&ms); err != nil {
		f.Close()
		go fail(fmt.Errorf("failed to decode %q: %w", f.Name(), err))
		return
	}
	if err = f.Close(); err != nil {
		go fail(fmt.Errorf("failed to close %q: %w", f.Name(), err))
		return
	}

	coll := colly.NewCollector(
		// Allow revisits, since we need to fetch many captchas.
		colly.AllowURLRevisit(),
	)
	// Disable cookies; they're not needed for the captchas.
	coll.DisableCookies()

	var cm sync.Map

	coll.OnHTML(`img[src^="CaptchaImage.aspx?guid="]`, func(img *colly.HTMLElement) {
		u, err := url.Parse(img.Request.AbsoluteURL(img.Attr("src")))
		if err != nil {
			errs <- fmt.Errorf("failed to parse img src: %w", err)
			return
		}
		id, err := uuid.Parse(u.Query().Get("guid"))
		if err != nil {
			errs <- fmt.Errorf("failed to parse ID as UUID: %w", err)
			return
		}

		cm.Store(id, &proto.Captcha{
			Id:   id.String(),
			Type: typ,
		})

		for n := samples; n > 0; n-- {
			if err := coll.Visit(u.String()); err != nil {
				errs <- fmt.Errorf("failed to fetch image at %q: %w", u, err)
			}
		}
	})

	coll.OnResponse(func(res *colly.Response) {
		if res.Headers.Get("content-type") != "image/jpeg" {
			return
		}
		ts, err := time.Parse(time.RFC1123, res.Headers.Get("date"))
		if err != nil {
			errs <- fmt.Errorf("failed to parse date header: %w", err)
			return
		}

		id, err := uuid.Parse(res.Request.URL.Query().Get("guid"))
		if err != nil {
			errs <- fmt.Errorf("failed to parse ID as UUID: %w", err)
			return
		}
		val, ok := cm.Load(id)
		if !ok {
			errs <- fmt.Errorf("captcha with UUID %q not found in map", id)
		}

		c := val.(*proto.Captcha)
		c.Samples = append(c.Samples, &proto.Captcha_Sample{
			Data:        res.Body,
			ContentType: "image/jpeg",
			GeneratedAt: timestamppb.New(ts),
			Sha1:        fmt.Sprintf("%x", sha1.Sum(res.Body)),
		})

		if len(c.Samples) == samples {
			cm.Delete(id)
			cs <- c
		}
	})

	go func() {
		if err := coll.Limit(&colly.LimitRule{
			DomainGlob: "katastar.rgz.gov.rs",
			Delay:      time.Second,
		}); err != nil {
			fail(fmt.Errorf("failed to set limit rule: %w", err))
			return
		}

		defer done()

		urls := make(chan string)
		go genCaptchaPageURLs(ctx, urls, ms, typ)

		for {
			select {
			case url := <-urls:
				if err := coll.Visit(url); err != nil {
					errs <- fmt.Errorf("failed to fetch page at %q: %w", url, err)
				}
			case <-ctx.Done():
				errs <- ctx.Err()
				coll.Wait()
				return
			}
		}
	}()

	return
}

// Generate URLs that should contain a Captcha image of the given type.
func genCaptchaPageURLs(ctx context.Context, urls chan<- string, ms []*proto.Municipality, typ proto.Captcha_Type) {
	for {
		var url string
		if typ == proto.Captcha_ALPHANUM_4 {
			url = gen4CaptchaPageURL(ms)
		} else if typ == proto.Captcha_ALPHANUM_5 {
			url = gen5CaptchaPageURL(ms)
		} else {
			close(urls)
			return
		}

		select {
		case urls <- url:
			// URL sent to output.
		case <-ctx.Done():
			close(urls)
			return
		}
	}
}

func gen4CaptchaPageURL(ms []*proto.Municipality) string {
	var url string
	if rand.Intn(2) == 0 {
		url = eKatFindObj
	} else {
		url = eKatFindAddr
	}

	if len(ms) > 0 {
		m := ms[rand.Intn(len(ms))]
		return fmt.Sprintf(url, m.Id)
	}

	// This should not happen.
	return strings.TrimSuffix(url, "?OpstinaID=%d")
}

func gen5CaptchaPageURL(ms []*proto.Municipality) string {
	url := eKatFindParc
	if len(ms) > 0 {
		m := ms[rand.Intn(len(ms))]
		if len(m.CadastralMunicipalities) > 0 {
			cm := m.CadastralMunicipalities[rand.Intn(len(m.CadastralMunicipalities))]
			return fmt.Sprintf(url, cm.Id)
		}
	}

	// This should not happen.
	return strings.TrimSuffix(url, "?KoID=%d")
}
