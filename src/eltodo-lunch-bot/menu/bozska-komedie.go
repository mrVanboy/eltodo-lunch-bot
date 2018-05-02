package menu

import (
	"github.com/antchfx/htmlquery"
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"errors"
	"eltodo-lunch-bot/cfg"
)

type BozskaKomedie struct {
	place
}


func (p *BozskaKomedie) Load(weekday Weekday) (string, error) {

	p.placeName = `Božská komedie`
	p.url = cfg.Get().UrlBK


	var menu string

	n, err := htmlquery.LoadURL(p.url)
	if err != nil {
		return ``, err
	}

	dm := p.getDailyMenuNode(n, weekday)
	if dm == nil {
		return "", errors.New(`can't find daily menu for current weekday`)
	}

	menu += fmt.Sprintf("%s\n", p.getHeading(dm))

	menu += p.getItems(dm)

	return menu, nil
}

func (p BozskaKomedie) getDailyMenuNode(root *html.Node, weekday Weekday) *html.Node{
	return  htmlquery.FindOne(root, `//*[@id='dailymenu']//table/../..`)

}

func (p BozskaKomedie) getHeading(root *html.Node) string {
	return htmlquery.FindOne(root, `//div/h3/text()`).Data
}

func (p BozskaKomedie) getItems(root *html.Node) string {
	var items string

	htmlquery.FindEach(root, `//table/*/tr`, func(_ int, node *html.Node) {
		items += "\n•"
		tds := htmlquery.Find(node, `/td`)
		for i := 0; i < len(tds); i++{
			item := htmlquery.InnerText(tds[i])
			item = strings.TrimSpace(item)
			item = strings.Replace(item, "\n", "", -1)
			item = strings.Replace(item, "   ", "", -1)
			if i+1 == len(tds){
				items += fmt.Sprintf(" *%s*", item)
			} else {
				items +=  ` ` + item
			}
		}


	})
	return items
}
