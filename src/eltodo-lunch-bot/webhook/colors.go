package webhook

import "math/rand"

var colors = []string{
	`#001f3f`,
	`#0074D9`,
	`#7FDBFF`,
	`#39CCCC`,
	`#3D9970`,
	`#2ECC40`,
	`#01FF70`,
	`#FFDC00`,
	`#FF851B`,
	`#FF4136`,
	`#85144b`,
	`#F012BE`,
	`#B10DC9`,
}

func getColors(length int) []string{
	c := colors
	rand.Shuffle(len(c), func(i, j int) {
		c[i], c[j] = c[j], c[i]
	})
	for len(c) <= length{
		c = append(c, colors[ rand.Intn( len(colors) - 1 ) ])
	}
	return c[0:length]
}