package scrapers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"github.com/attilaolah/cad-rs/proto"
	"github.com/attilaolah/cad-rs/text"
)

const eKatSearchStreets = eKatURL + "/FindAdresa.aspx/PretragaUlica"

type StreetSearchResults struct {
	Query   string          `json:"query"`
	Results []*proto.Street `json:"results"`
}

// ScrapeStreets fetches streets for a single municipality.
func ScrapeStreets(dir string, mID int64) (chan *StreetSearchResults, chan error) {
	ss := make(chan *StreetSearchResults)
	errs := make(chan error)

	done := func() {
		close(errs)
		close(ss)
	}
	fail := func(err error) {
		errs <- err
		done()
	}

	coll := colly.NewCollector(
		// Allow revisits, since we need to fetch many captchas.
		colly.AllowURLRevisit(),
	)
	// Disable cookies; they're not needed for the captchas.
	coll.DisableCookies()
	// Set a longer timeout, since the server can be pretty slow.
	coll.SetRequestTimeout(time.Minute * 2)

	coll.OnRequest(func(r *colly.Request) {
		r.Headers.Set("content-type", "application/json; charset=utf-8")
	})

	// Buffer size 1 to make sure only a single request is handled at a time.
	buf := make(chan *StreetSearchResults, 1)

	coll.OnResponse(func(res *colly.Response) {
		if ct := strings.ToLower(res.Headers.Get("content-type")); ct != "application/json; charset=utf-8" {
			errs <- fmt.Errorf("got unexpected response with content-type %q", ct)
			return
		}

		data := struct {
			D []string `json:"d"`
		}{}
		if err := json.Unmarshal(res.Body, &data); err != nil {
			errs <- fmt.Errorf("failed to decode response: %w", err)
		}

		sr := <-buf
		for _, s := range data.D {
			row := struct {
				First, Second string
			}{}
			if err := json.Unmarshal([]byte(s), &row); err != nil {
				errs <- fmt.Errorf("failed to unmarshal row %q: %w", s, err)
				continue
			}
			id, err := strconv.ParseInt(row.Second, 10, 64)
			if err != nil {
				errs <- fmt.Errorf("failed to parse %q as integer: %w", row.Second, err)
			}

			st := proto.Street{
				Id:       id,
				FullName: cleanup(row.First),
			}

			if st.Id != -1 && st.FullName != "NEMA REZULTATA PRETRAGE" {
				sr.Results = append(sr.Results, &st)
			}
		}

		ss <- sr
	})

	subdir := filepath.Join(dir, strconv.FormatInt(mID, 10))

	go func() {
		if err := coll.Limit(&colly.LimitRule{
			DomainGlob:  "katastar.rgz.gov.rs",
			Parallelism: 1,
		}); err != nil {
			fail(fmt.Errorf("failed to set limit rule: %w", err))
			return
		}

		qs := make(chan string)
		go genStreetSearchQueries(qs, errs, subdir)
		for q := range qs {
			buf <- &StreetSearchResults{
				Query:   cleanup(q),
				Results: []*proto.Street{},
			}

			data, err := json.Marshal(struct {
				PrefixText string `json:"prefixText"`
				ContextKey string `json:"contextKey"`
				Count      int    `json:"count"`
			}{
				PrefixText: q,
				ContextKey: strconv.FormatInt(mID, 10),
				Count:      1000,
			})
			if err != nil {
				errs <- fmt.Errorf("failed to encode query data: %w", err)
				continue
			}

			if err := coll.PostRaw(eKatSearchStreets, data); err != nil {
				errs <- fmt.Errorf("failed to fetch page at %q w/ data = %s: %w", eKatSearchStreets, data, err)
				<-buf
			}
		}

		coll.Wait()
		close(buf)
		done()
	}()

	return ss, errs
}

func genStreetSearchQueries(qs chan<- string, errs chan<- error, subdir string) {
	defer close(qs)

	tmp := make([]string, len(text.Azbuka)*len(text.Azbuka))
	for i, a := range text.Azbuka {
		for j, b := range text.Azbuka {
			tmp[len(text.Azbuka)*i+j] = string([]rune{a, b})
		}
	}

	rand.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})

	for _, q := range tmp {
		fn := filepath.Join(subdir, fmt.Sprintf("%s.json", asciil(q)))
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			qs <- q
		} else if err != nil {
			errs <- fmt.Errorf("error checking for file %q: %w", fn, err)
		}
	}
}

func asciil(s string) string {
	return strings.ToLower(ascii(s))
}

func ascii(s string) string {
	s = text.ToLatin.Replace(s)
	s = text.RemoveDigraphs.Replace(s)
	s = text.ToASCII.Replace(s)
	return s
}
