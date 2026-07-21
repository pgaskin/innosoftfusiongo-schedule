package ifgsch

import (
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/pgaskin/go-lightningcss"
)

//go:embed style.css
var rawStyleCSS []byte

var styleCSS = mustOnce("compile css", func() (string, error) {
	res, err := lightningcss.Transform(rawStyleCSS, &lightningcss.Options{
		Filename: "style.css",
		Minify:   true,
		Nesting:  true,
		Targets: lightningcss.Targets{
			Chrome:  lightningcss.Version(110, 0, 0),
			Safari:  lightningcss.Version(15, 0, 0),
			Firefox: lightningcss.Version(110, 0, 0),
		},
	})
	if err != nil {
		return "", fmt.Errorf("lightningcss: %w", err)
	}
	return string(res.Code), nil
})

func init() {
	go func() {
		slog.Info("compiling css")
		defer slog.Info("css ready")

		styleCSS()
	}()
}
