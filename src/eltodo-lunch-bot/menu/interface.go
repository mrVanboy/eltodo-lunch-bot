package menu

type IMenu interface {
	GetUrl() string
	GetPlaceName() string
	Load(weekday Weekday) (string, error)
}

type place struct {
	url string
	placeName string
}

func (p *place) setUrl(url string)  {
	p.url = url
}

func (p *place) setName(placeName string)  {
	p.placeName = placeName
}

func (p place) GetUrl() string {
	return p.url
}

func (p place) GetPlaceName() string{
	return p.placeName
}