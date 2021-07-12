package scrapers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"

	"github.com/attilaolah/cad-rs/proto"
	"github.com/attilaolah/cad-rs/text"
)

// The eKatastar Public Access URL.
const (
	eKatURL       = "https://katastar.rgz.gov.rs/eKatastarPublic"
	eKatPubAccess = eKatURL + "/PublicAccess.aspx"
)

// ScrapeMunicipalities fetches all municipality data.
func ScrapeMunicipalities() ([]*proto.Municipality, error) {
	mmap := map[int64]*proto.Municipality{}
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
		mmap[id] = &proto.Municipality{
			Id:   id,
			Name: cleanup(opt.Text),
		}

		if err := c.Request(http.MethodGet, eKatPubAccess, nil, nil, http.Header{
			"cookie": []string{fmt.Sprintf("KnWebPublicGetOpstinaKO=SelectedValueOpstina=%d", id)},
		}); err != nil {
			errs <- err
		}
	})

	c.OnHTML("table#ContentPlaceHolder1_getOpstinaKO_GridView>tbody>tr:not(.header)", func(tr *colly.HTMLElement) {
		cm := proto.CadastralMunicipality{}
		tr.ForEach("td", func(col int, td *colly.HTMLElement) {
			if col == 0 {
				s := td.ChildAttr("img", "src")
				s = strings.TrimSuffix(strings.TrimPrefix(s, "images/kn_status_"), ".gif")
				typ, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					errs <- fmt.Errorf("error parsing cadastre type: %w", err)
					return
				}
				cm.CadastreType = proto.CadastralMunicipality_CadastreType(typ)
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
				cm.Id = id
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
			errs <- fmt.Errorf("failed to set limit rule: %w", err)
		}
		if err := c.Visit(eKatPubAccess); err != nil {
			errs <- fmt.Errorf("failed to fetch page at %q: %w", eKatPubAccess, err)
		}
		c.Wait()
		close(errs)
	}()

	for err := range errs {
		return nil, err
	}

	ms := []*proto.Municipality{}
	for _, m := range mmap {
		sort.Slice(m.CadastralMunicipalities, func(i, j int) bool {
			return m.CadastralMunicipalities[i].Id < m.CadastralMunicipalities[j].Id
		})
		ms = append(ms, m)
	}
	sort.Slice(ms, func(i, j int) bool { return ms[i].Id < ms[j].Id })

	return ms, nil
}

func cleanup(s string) string {
	s = strings.TrimSpace(s)
	s = text.ToLatin.Replace(s)
	s = text.RemoveDigraphs.Replace(s)
	s = strings.ToUpper(s)
	return s
}
