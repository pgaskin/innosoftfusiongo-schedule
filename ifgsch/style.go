package ifgsch

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/pgaskin/go-lightningcss"
	"github.com/pgaskin/innosoftfusiongo-schedule/m3color"
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

// colorCSS caches generated MD3 palettes by lowercased hex color.
var colorCSS sync.Map

// md3CSS returns the MD3 palette CSS for a hex color, caching the result. It
// panics if the color is invalid.
func md3CSS(c string) string {
	c = strings.ToLower(c)
	if v, ok := colorCSS.Load(c); ok {
		return v.(string)
	}
	v, err := m3color.PaletteCSS(c)
	if err != nil {
		panic(fmt.Errorf("generate md3 palette css for color %s: %w", c, err))
	}
	colorCSS.Store(c, v)
	return v
}

func fontCSS() string {
	return strings.Join([]string{
		"@font-face {",
		"\tfont-family: 'Asap SemiCondensed';",
		"\tfont-style: normal;",
		"\tfont-weight: 100 900;",
		"\tfont-stretch: 87.5%;",
		"\tfont-display: swap;",
		"\tsrc: url('data:font/woff2;base64," + base64.StdEncoding.EncodeToString(asapWOFF2()) + "') format('woff2-variations');",
		"}",
		"@font-face {",
		"\tfont-family: 'Material Symbols Subset';",
		"\tfont-style: normal;",
		"\tfont-weight: 300;",
		"\tsrc: url('data:font/woff2;base64," + base64.StdEncoding.EncodeToString(symbolsWOFF2()) + "') format('woff2');",
		"}",
	}, "\n")
}

func init() {
	go func() {
		slog.Info("compiling css")
		defer slog.Info("css ready")

		styleCSS()
	}()
}
