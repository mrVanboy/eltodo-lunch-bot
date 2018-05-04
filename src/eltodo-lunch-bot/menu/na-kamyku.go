package menu

import (
	"io/ioutil"
	"errors"
	"os"
	"net/http"
	"io"
	"os/exec"
	"eltodo-lunch-bot/cfg"
	"bytes"
	"regexp"
)

type NaKamyku struct {
	place
}


func (p *NaKamyku) Load(weekday Weekday) (string, error) {
	p.placeName = `Na Kamýku`
	p.url = cfg.Get().UrlNK

	pdf, err := p.downloadPdf()
	if err != nil {
		return ``, err
	}
	defer removeAndClose(pdf)

	txt, err := p.convertPdfToTxt(pdf)
	if err != nil {
		return ``, err
	}
	defer removeAndClose(txt)

	const dayRegexpTemplate = `(?i)((?:%v\n*(?s)(?:.+?)))\n{2,}`
	bContent,err := ioutil.ReadAll(txt)
	if err != nil {
		return ``, err
	}

	rx := Weekday(weekday).buildRegexp(dayRegexpTemplate)
	matches := rx.FindSubmatch(bContent)
	if len(matches) < 1 {
		return ``, errors.New(`can't find menu for current day`)
	}
	menu := string(matches[1])

	//cleanup spaces
	menu = regexp.MustCompile(`(?: {2,})`).ReplaceAllString(menu, ` `)
	// cleanup alergens
	menu = regexp.MustCompile(`(?m)(?: +\d+)+( \d)`).ReplaceAllString(menu, `$1`)
	// cleanup whitespaces
	menu = regexp.MustCompile(`([[:lower:]][[:punct:]])\s*\n`).ReplaceAllString(menu, `$1`)
	menu = regexp.MustCompile(` (\d{2,3}[[:punct:]]*\d*)(\n|\z)`).ReplaceAllString(menu, ` *$1*$2`)
	menu = regexp.MustCompile("\n(.+)").ReplaceAllString(menu, "\n• $1")
	return menu, nil
}

func (p *NaKamyku) downloadPdf() (*os.File, error){
	url := p.url

	output, err := p.createTempFile()
	fileName := output.Name()

	// fmt.Println("Downloading", url, "to", fileName)
	if err != nil {
		return nil, errors.New("Error while creating " + fileName + "-" + err.Error())
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New("error while downloading" + url + "-" + err.Error())
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return nil, errors.New("error while downloading" + url + "-" + err.Error())
	}

	// fmt.Println(n, "bytes downloaded.")

	return output, nil

}

func (p NaKamyku) createTempFile() (*os.File, error){
	return ioutil.TempFile("", "na-kamyku")
}

func (p NaKamyku) convertPdfToTxt(pdfFile *os.File) (*os.File, error){
	txtFile, err := p.createTempFile()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(cfg.Get().PathToPdfToText, `-layout`, pdfFile.Name(), txtFile.Name())
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		defer removeAndClose(txtFile)
		println(`STDERR:`, errBuf.String())
		println(`STDOUT:`, outBuf.String())
		return nil, errors.New(`pdf converting exec error - ` + err.Error())
	}

	if outBuf.Len() > 0 {
		defer removeAndClose(txtFile)
		return nil, errors.New(`unexpected stdout from pdftotxt -` + outBuf.String())
	}
	if errBuf.Len() > 0 {
		defer removeAndClose(txtFile)
		return nil, errors.New(`unexpected stderr from pdftotxt -` + errBuf.String())
	}
	return txtFile, nil
}

func removeAndClose(file *os.File) {
	file.Close()
	err := os.Remove(file.Name())
	if err != nil {
		println(err.Error())
	}
}

