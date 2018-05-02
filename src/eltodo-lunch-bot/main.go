package main

import (
	"eltodo-lunch-bot/menu"
	"eltodo-lunch-bot/cfg"
	"fmt"
	"os"
	"eltodo-lunch-bot/webhook"
	"errors"
	"time"
	"gopkg.in/robfig/cron.v2"
	"sync"
	"strconv"
)

var menicka []menu.IMenu
var loc *time.Location

func init() {
	menicka = []menu.IMenu{
		&menu.NaKamyku{},
		&menu.BozskaKomedie{},
		&menu.DiCarlo{},
		&menu.Pepino{},
	}
}

func main() {
	err := cfg.Load()
	if err != nil {
		panic(err)
	}
	loc, err = time.LoadLocation(cfg.Get().TimeZone)
	if err != nil {
		panic(err)
	}

	c := cron.New()
	cronPattern := fmt.Sprintf(`TZ=%s %s`, cfg.Get().TimeZone, cfg.Get().Cron)

	entryId, err := c.AddFunc(cronPattern, getAndSend)
	if err != nil {
		panic(err)
	}


	fmt.Fprintln(os.Stdout, `Add job to cron, id: ` + strconv.FormatInt(int64(entryId), 10))

	wg := sync.WaitGroup{}
	wg.Add(1)

	c.Start()
	fmt.Fprintln(os.Stdout, `Next run in : ` + c.Entry(entryId).Next.String())
	wg.Wait()
}

func getAndSend(){
	var errArr []error
	for _, m := range menicka {
		dailyMenu, err := m.Load(menu.Weekday(time.Now().In(loc).Weekday()))
		time.Now().Location()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n",err.Error())
			errArr = append(errArr, errors.New(`getting dm from `+ m.GetPlaceName() + `error: ` + err.Error()))
			continue
		}

		a := webhook.Attachment{
			Fallback: 	fmt.Sprintf(`%s - denni  menu - %s`, m.GetPlaceName(), m.GetUrl()),
			Title: 		m.GetPlaceName(),
			TitleLink: 	m.GetUrl(),
			Text: dailyMenu,
		}
		webhook.NewAttachment(a)
	}
	j, _ := webhook.BuildJSON()
	fmt.Fprint(os.Stdout, string(j))
	err := webhook.Send()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n",err.Error())
		errArr = append(errArr,  errors.New(`calling obedbot webhook error: ` + err.Error()))
	}

	if len(errArr) > 0 {
		err := webhook.NotifyAboutErrors(errArr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}


