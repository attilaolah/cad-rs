package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// The eKatastar Public Access URL.
var eKatURL *url.URL

func init() {
	var err error
	if eKatURL, err = url.Parse("https://katastar.rgz.gov.rs/eKatastarPublic/PublicAccess.aspx"); err != nil {
		panic(err)
	}
}

// Municipality represents a municipality as described here:
// https://en.wikipedia.org/wiki/Municipalities_and_cities_of_Serbia
type Municipality struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	CadastralMunicipalities []*CadastralMunicipality `json:"cadastral_municipalities,omitempty"`
}

// CadastralMunicipality represents a cadastral municipality.
type CadastralMunicipality struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func FetchMunicipalities() ([]*Municipality, error) {
	res, err := http.Get(eKatURL.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching %q: %w", eKatURL, err)
	}
	defer func() {
		err = res.Body.Close()
	}()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing page at %q: %w", eKatURL, err)
	}

	const (
		tag = "select"
		id  = "ContentPlaceHolder1_getOpstinaKO_dropOpstina"
	)
	sel := byTagID(doc, tag, id)
	if sel == nil {
		return nil, fmt.Errorf("error parsing page at %q: <%s id=%q> not found", eKatURL, tag, id)
	}

	ms := []*Municipality{}
	for opt := sel.FirstChild; opt != nil; opt = opt.NextSibling {
		if opt.Data != "option" {
			continue
		}
		tn := opt.FirstChild
		if tn == nil {
			return nil, fmt.Errorf("error parsing %s: no child node", opt)
		}
		m := Municipality{
			Name: strings.TrimSpace(strings.ToUpper(RemoveDigraphs.Replace(ToLatin.Replace(tn.Data)))),
		}

		// Find the ID
		for _, a := range opt.Attr {
			if a.Key == "value" {
				if m.ID, err = strconv.ParseInt(a.Val, 10, 64); err != nil {
					return nil, fmt.Errorf("error parsing %s: error parsing ID from %q: %w", opt, a.Val, err)
				}
				break
			}
		}
		if m.ID == 0 {
			return nil, fmt.Errorf("error parsing %s: could not parse ID", opt)
		}
		ms = append(ms, &m)
	}

	sort.Slice(ms, func(i, j int) bool { return ms[i].ID < ms[j].ID })

	return ms, err
}

// Populate fetches cadastral municipalities.
func (m *Municipality) Populate() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("error creating cookie jar: %w", err)
	}
	jar.SetCookies(eKatURL, []*http.Cookie{
		&http.Cookie{Name: "KnWebPublicGetOpstinaKO", Value: fmt.Sprintf("SelectedValueOpstina=%d", m.ID)},
	})

	res, err := (&http.Client{Jar: jar}).Get(eKatURL.String())
	if err != nil {
		return err
	}
	defer func() {
		err = res.Body.Close()
	}()

	//io.Copy(os.Stdout, res.Body)

	doc, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("error parsing page at %q: %w", eKatURL, err)
	}

	const (
		tag = "table"
		id  = "ContentPlaceHolder1_getOpstinaKO_GridView"
	)
	tbl := byTagID(doc, tag, id)
	if tbl == nil {
		return fmt.Errorf("error parsing page at %q: <%s id=%q> not found", eKatURL, tag, id)
	}

	return err
}

// Deptch-first search by HTML tag & ID attribute.
func byTagID(n *html.Node, tag, id string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		//fmt.Println("###", n)
		for _, a := range n.Attr {
			//fmt.Println("### > ", a)
			//if a.Key == "id" {
			//	fmt.Println("???", a.Val, "==", id, "->", a.Val == id)
			//}
			if a.Key == "id" && a.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c := byTagID(c, tag, id); c != nil {
			return c
		}
	}
	return nil
}

func main() {
	ms, err := FetchMunicipalities()
	if err != nil {
		log.Fatalf("error fetching municipalities: %v", err)
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(&ms); err != nil {
		log.Fatalf("error encoding municipalities: %v", err)
	}
}
