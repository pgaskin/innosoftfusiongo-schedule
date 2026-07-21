package ifgsch

import (
	_ "embed"
	"fmt"
	"log/slog"
	"sync"
	"unicode"

	"github.com/pgaskin/go-hbsubset"
	"github.com/pgaskin/go-woff2"
)

//go:generate go run fetch.go https://github.com/Omnibus-Type/Asap/raw/ca471c0ccf90a5c66155d4bcaa020859804ffd00/fonts/variable/Asap%5Bwdth%2Cwght%5D.ttf asap.ttf
//go:generate go run fetch.go https://github.com/google/material-design-icons/raw/fe742c4072d4e3b8b899170109d9f710e89f082e/variablefont/MaterialSymbolsOutlined%5BFILL%2CGRAD%2Copsz%2Cwght%5D.ttf symbols.ttf

var (
	//go:embed asap.ttf
	asapTTF []byte
	//go:embed symbols.ttf
	symbolsTTF []byte
)

var keepNames = []hbsubset.NameID{
	hbsubset.NameUniqueID,
	hbsubset.NameCopyright,
	hbsubset.NameFontFamily,
	hbsubset.NameFontSubfamily,
	hbsubset.NameTypographicFamily,
	hbsubset.NameTypographicSubfamily,
	hbsubset.NameFullName,
	hbsubset.NameMacFullName,
	hbsubset.NamePostscriptName,
	hbsubset.NameManufacturer,
	hbsubset.NameDescription,
	hbsubset.NameVariationsPSPrefix,
}

var asapWOFF2 = mustOnce("subset asap", func() ([]byte, error) {
	sub, err := hbsubset.Subset(asapTTF, 0, &hbsubset.Options{
		UnicodeRanges: &unicode.RangeTable{
			R16: []unicode.Range16{
				{Stride: 1, Lo: 32, Hi: 126},            // space + ascii printable
				{Stride: 1, Lo: '\u2002', Hi: '\u201e'}, // spaces, smart punctuation
				{Stride: 1, Lo: '\u2022', Hi: '\u2022'}, // bullet
				{Stride: 1, Lo: '\u2026', Hi: '\u2026'}, // ellipsis
			},
		},
		PinAllAxesToDefault: true,
		AxisRanges: map[hbsubset.Tag]hbsubset.AxisRange{
			hbsubset.MakeTag("wght"): {Min: 100, Max: 900, Default: 400},
		},
		PinAxes: map[hbsubset.Tag]float32{
			hbsubset.MakeTag("wdth"): 87.5, // SemiCondensed
		},
		LayoutFeatures: []hbsubset.Tag{
			hbsubset.MakeTag("kern"),
			hbsubset.MakeTag("liga"),
		},
		LayoutScripts: []hbsubset.Tag{
			hbsubset.MakeTag("latn"), // latin
		},
		NameIDs:       keepNames,
		NameLanguages: []uint32{1033}, // english
	})
	if err != nil {
		return nil, fmt.Errorf("subset: %w", err)
	}
	enc, err := woff2.Encode(sub, nil)
	if err != nil {
		return nil, fmt.Errorf("woff2: %w", err)
	}
	return enc, nil
})

var symbolsWOFF2 = mustOnce("subset symbols", func() ([]byte, error) {
	sub, err := hbsubset.Subset(symbolsTTF, 0, &hbsubset.Options{
		Unicodes: []rune{
			'', // schedule
			'', // location_on
		},
		PinAllAxesToDefault: true,
		PinAxes: map[hbsubset.Tag]float32{
			hbsubset.MakeTag("opsz"): 20,
			hbsubset.MakeTag("wght"): 300,
			hbsubset.MakeTag("FILL"): 0,
			hbsubset.MakeTag("GRAD"): 0,
		},
		DropTables: []hbsubset.Tag{
			hbsubset.MakeTag("GSUB"), // ligature name lookups, unused
		},
		NameIDs:       keepNames,
		NameLanguages: []uint32{1033}, // english
	})
	if err != nil {
		return nil, fmt.Errorf("subset: %w", err)
	}
	enc, err := woff2.Encode(sub, nil)
	if err != nil {
		return nil, fmt.Errorf("woff2: %w", err)
	}
	return enc, nil
})

func init() {
	go func() {
		slog.Info("subsetting fonts")
		defer slog.Info("fonts ready")

		var wg sync.WaitGroup
		defer wg.Wait()

		wg.Go(func() { asapWOFF2() })
		wg.Go(func() { symbolsWOFF2() })
	}()
}
