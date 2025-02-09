package main

import (
	"github.com/alserok/url_shortener/internal/app"
	"github.com/alserok/url_shortener/internal/config"
)

func main() {
	app.MustStart(config.MustLoad())
}
