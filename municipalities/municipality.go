package municipalities

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/attilaolah/ekat/latin"
	"github.com/gocolly/colly"
)

// The eKatastar Public Access URL.
const eKatURL = "https://katastar.rgz.gov.rs/eKatastarPublic/PublicAccess.aspx"

// Municipality represents a municipality as described here:
// https://en.wikipedia.org/wiki/Municipalities_and_cities_of_Serbia
type Municipality struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	CadastralMunicipalities []*CadastralMunicipality `json:"cadastral_municipalities,omitempty"`
}

// FetchAll fetches all municipality data.
func FetchAll() (chan *Municipality, chan error) {
	ms := make(chan *Municipality)
	errs := make(chan error)
	seen := map[int64]struct{}{}

	c := colly.NewCollector(
		// Allow revisits, since only the cookie differs.
		colly.AllowURLRevisit(),
	)
	// But disable cookie handling; we'll set the cookie manually.
	c.DisableCookies()

	c.OnHTML("select#ContentPlaceHolder1_getOpstinaKO_dropOpstina>option", func(opt *colly.HTMLElement) {
		id, err := strconv.ParseInt(opt.Attr("value"), 10, 64)
		if err != nil {
			errs <- fmt.Errorf("error parsing ID: %w", err)
			return
		}
		if _, ok := seen[id]; ok {
			return // already visited
		}
		seen[id] = struct{}{}

		ms <- &Municipality{
			ID:   id,
			Name: cleanup(opt.Text),
		}

		if id != 80438 { // Subotica
			return
		}

		if err := c.Request(http.MethodGet, eKatURL, nil, nil, http.Header{
			"cookie": []string{fmt.Sprintf("KnWebPublicGetOpstinaKO=SelectedValueOpstina=%d", id)},
		}); err != nil {
			errs <- err
		}
	})

	c.OnHTML("table#ContentPlaceHolder1_getOpstinaKO_GridView", func(tbl *colly.HTMLElement) {
		tbl.ForEach("tr", func(row int, tr *colly.HTMLElement) {
			cm := CadastralMunicipality{}
			tr.ForEach("td", func(col int, td *colly.HTMLElement) {
				val := cleanup(td.Text)
				if col == 1 {
					cm.Name = val
					return
				}
				if col == 2 {
					id, err := strconv.ParseInt(val, 10, 64)
					if err != nil {
						errs <- fmt.Errorf("error parsing ID: %w", err)
						return
					}
					cm.ID = id
					return
				}
				if col == 3 {
					// TODO: Parse Municipality ID!
				}
			})
			if cm.ID == 0 {
				return
			}
			//fmt.Println("KM:", cm)
		})
	})

	c.OnResponse(func(res *colly.Response) {
		//fmt.Println("Got response for:", res.Request.URL, res.Request.Headers)
	})

	c.OnError(func(res *colly.Response, err error) {
		errs <- err
	})

	go func() {
		if err := c.Limit(&colly.LimitRule{
			DomainGlob:  "katastar.rgz.gov.rs",
			Parallelism: 2,
		}); err != nil {
			errs <- err
		}
		if err := c.Visit(eKatURL); err != nil {
			errs <- err
		}
		c.Wait()
		close(errs)
		close(ms)
	}()

	return ms, errs
}

func cleanup(text string) string {
	return strings.TrimSpace(strings.ToUpper(latin.RemoveDigraphs.Replace(latin.ToLatin.Replace(text))))
}
