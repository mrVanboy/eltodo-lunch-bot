package menu

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"errors"
	"strings"
	"regexp"
	"eltodo-lunch-bot/cfg"
)

type Pepino struct {
	place
}

func (p *Pepino) Load(weekday Weekday) (string, error) {
	p.placeName = `Pizzeria Pepino`
	p.url = cfg.Get().UrlPP
	// load web
	n, err := htmlquery.LoadURL(p.url)
	if err != nil {
		return ``, err
	}
	if n == nil {
		return "", errors.New(`can't load html'`)
	}

	// find day
	dm, err := p.getDailyMenuNode(n, weekday)
	if err != nil {
		return ``, err
	}
	// get menu
	var menu string

	menu += p.getHeading(dm) + "\n"

	menu += p.getItems(dm)
	menu = regexp.MustCompile(`(?m) (\d{2,3}[[:punct:]]*\d*\s\p{L}+)(\n|\z)`).ReplaceAllString(menu,` *$1*$2`)
	return menu, nil
}

func (p Pepino) getDailyMenuNode(node *html.Node, weekday Weekday) (*html.Node, error) {
	rx := weekday.buildRegexp(`(?i)%s`)
	dms := htmlquery.Find(node, `//div[contains(@class, "content")]`)
	for _, dm := range dms {
		text := htmlquery.InnerText(dm)
		if rx.MatchString(text) {
			return dm, nil
		}
	}
	return nil, errors.New(`can't find daily menu node for current week`)
}
func (p Pepino) getHeading(node *html.Node) string {
	var heading string
	htmlquery.FindEach(node, `//*[substring(name(), 0, 1) = "h"]/text()`, func(_ int, h *html.Node) {
		heading += strings.TrimSpace(h.Data) + " "
	})
	return heading
}
func (p Pepino) getItems(node *html.Node) string {
	var items string
	htmlquery.FindEach(node, `//tr`, func(_ int, rowNode *html.Node) {
		items += "\n"
		htmlquery.FindEach(rowNode, `/td/text()`, func(_ int, columnNode *html.Node) {
			items += ` ` + strings.TrimSpace(columnNode.Data)
		})
	})
	return items
}

