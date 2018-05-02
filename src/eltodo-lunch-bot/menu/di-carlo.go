package menu

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strconv"
	"errors"
	"fmt"
	"strings"
	"regexp"
	"eltodo-lunch-bot/cfg"
)

type DiCarlo struct {
	place
	weekday Weekday
}


func (p *DiCarlo) Load(weekday Weekday) (string, error) {
	p.placeName = `Di Carlo Lhotka`
	p.url = cfg.Get().UrlDC
	p.weekday = weekday

	var menu string

	n, err := htmlquery.LoadURL(p.url)
	if err != nil {
		return ``, err
	}
	dm := p.getDailyMenuNode(n)
	heading, id, err := p.getHeading(dm)
	if err != nil {
		return ``, err
	}

	menu += fmt.Sprintf("%s\n", heading)

	items, err := p.getItems(dm, id)
	if err != nil {
		return "", err
	}
	menu += items

	menu = regexp.MustCompile(`\.+.*?\.+\s*\n`).ReplaceAllString(menu, "\n")
	menu = regexp.MustCompile("\n(.+)").ReplaceAllString(menu, "\n• $1")

	return menu, nil
}

func (p DiCarlo) getDailyMenuNode(root *html.Node) *html.Node{
	return htmlquery.FindOne(root, `//*[@id='daily-menu']`)
}

func (p DiCarlo) getHeading(root *html.Node) (string, int, error) {
	var heading, sId string
	founded := false
	rx := p.weekday.buildRegexp(`(?i)%s`)
	htmlquery.FindEach(root, `//a[@data-daily-menu]`, func(_ int, node *html.Node) {
		if founded { return }

		text := htmlquery.InnerText(node)
		if rx.MatchString(text) {
			sId = htmlquery.SelectAttr(node, `data-daily-menu`)
			heading = htmlquery.InnerText(node)
			founded = true
		}
	})

	if !founded {
		return "", 0, errors.New(`can't find current day of week in //a[@data-daily-menu]`)
	}

	id, err := strconv.ParseInt(sId, 10, 32)
	if err != nil {
		return "", 0, err
	}

	return heading, int(id), nil
}

func (p DiCarlo) getItems(root *html.Node, id int) (string, error) {
	var items string

	n := htmlquery.Find(root, `//div[contains(@class,"daily-menu-section__table")]`)
	if id > len(n){
		return ``, errors.New(`id from href is greater than length of array of nodes .daily-menu-section__table`)
	}
	currentDayNode := n[id-1]

	htmlquery.FindEach(currentDayNode, `//div[contains(@class, "food-menu__row")]`, func(_ int, node *html.Node) {
		priceNode := htmlquery.FindOne(node, `//div[contains(@class, "price")]`)
		price := htmlquery.InnerText(priceNode)
		price = strings.TrimSpace(price)

		foodNodes := htmlquery.Find(node, `//div[contains(@class, "food-menu__desc")]/*[not(self::div)]`)
		var food string
		for _, foodNode := range foodNodes {
			content := strings.TrimSpace(htmlquery.InnerText(foodNode))
			if len(content) > 0 {
				if len(food) > 0 { food += " "}
				food += content
			}
		}

		if (len(food) + len(price)) > 0 {
			item := food

			if len(price) > 0 {
				item += fmt.Sprintf(` *%s*`, price)
			}

			items += fmt.Sprintf("\n%s", item)
		}
	})

	return items, nil
}
