package municipalities

import (
	"fmt"
	"net/http"
	"sort"
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
func FetchAll() ([]*Municipality, error) {
	mmap := map[int64]*Municipality{}
	errs := make(chan error)

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
		if _, ok := mmap[id]; ok {
			return // already visited
		}
		mmap[id] = &Municipality{
			ID:   id,
			Name: cleanup(opt.Text),
		}

		if err := c.Request(http.MethodGet, eKatURL, nil, nil, http.Header{
			"cookie": []string{fmt.Sprintf("KnWebPublicGetOpstinaKO=SelectedValueOpstina=%d", id)},
		}); err != nil {
			errs <- err
		}
	})

	c.OnHTML("table#ContentPlaceHolder1_getOpstinaKO_GridView>tbody>tr:not(.header)", func(tr *colly.HTMLElement) {
		cm := CadastralMunicipality{}
		tr.ForEach("td", func(col int, td *colly.HTMLElement) {
			if col == 0 {
				s := td.ChildAttr("img", "src")
				s = strings.TrimSuffix(strings.TrimPrefix(s, "images/kn_status_"), ".gif")
				typ, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					errs <- fmt.Errorf("error parsing cadastre type: %w", err)
					return
				}
				cm.Type = CadastreType(typ)
				return
			}
			if col == 1 {
				cm.Name = cleanup(td.Text)
				return
			}
			if col == 2 {
				s := strings.TrimSpace(td.Text)
				id, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					errs <- fmt.Errorf("error parsing ID: %w", err)
					return
				}
				cm.ID = id
				return
			}
			if col == 3 {
				s := td.ChildAttr("a", "href")
				s = strings.TrimPrefix(s, "FindObjekat.aspx?OpstinaID=")
				id, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					errs <- fmt.Errorf("error parsing ID: %w", err)
					return
				}
				m := mmap[id]
				if m == nil {
					errs <- fmt.Errorf("municipality id=%d not found", id)
					return
				}
				m.CadastralMunicipalities = append(m.CadastralMunicipalities, &cm)
			}
		})
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
	}()

	for err := range errs {
		return nil, err
	}

	ms := []*Municipality{}
	for _, m := range mmap {
		sort.Slice(m.CadastralMunicipalities, func(i, j int) bool {
			return m.CadastralMunicipalities[i].ID < m.CadastralMunicipalities[j].ID
		})
		ms = append(ms, m)
	}
	sort.Slice(ms, func(i, j int) bool { return ms[i].ID < ms[j].ID })

	return ms, nil
}

func cleanup(text string) string {
	text = strings.TrimSpace(text)
	text = latin.ToLatin.Replace(text)
	text = latin.RemoveDigraphs.Replace(text)
	text = strings.ToUpper(text)
	return text
}
