// Package m3color wraps some of the functionality from the Material Color
// Utilities JavasScript library.
package m3color

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/dop251/goja"
)

//go:generate go run ./m3fetch
//go:embed mcu.js
var mcuJS []byte

var mcu sync.Pool

func init() {
	if prog, err := goja.Compile("mcu.js", string(mcuJS), true); err != nil {
		panic(fmt.Errorf("m3color: failed to compile: %w", err))
	} else {
		mcu.New = func() any {
			vm := goja.New()
			if _, err := vm.RunProgram(prog); err != nil {
				panic(fmt.Errorf("m3color: failed to init: %w", err))
			}
			return vm
		}
		mcu.Put(mcu.Get())
	}
}

func eval[T string | int64 | float64 | bool](args, fn string, arg ...any) (T, error) {
	vm := mcu.Get().(*goja.Runtime)
	defer mcu.Put(vm)

	var z T

	f, err := vm.RunString(`(` + args + `)=>{` + fn + `}`)
	if err != nil {
		return z, err
	}
	c, _ := goja.AssertFunction(f)

	a := make([]goja.Value, len(arg))
	for i, x := range arg {
		a[i] = vm.ToValue(x)
	}

	v, err := c(nil, a...)
	if err != nil {
		return z, err
	}

	z, ok := v.Export().(T)
	if !ok {
		return z, fmt.Errorf("value %q is not %T", v, v)
	}
	return z, nil
}

func PaletteCSS(c string) (string, error) {
	if c == "" {
		c = "6750A4" // M3 baseline color
	}
	return eval[string](`c`, `
		const a = argbFromHex(c)
		const t = themeFromSourceColor(a)
		const v = [["primary","primary"],["secondary","secondary"],["tertiary","tertiary"],["neutral","neutral"],["neutralVariant","neutral-variant"],["error","error"]]
		const n = [0,4,5,6,10,12,17,20,22,24,25,30,35,40,50,60,70,80,87,90,92,94,95,96,98,99,100]
		return ":root{--md-source:" + hexFromArgb(a) + ";" + v.flatMap(([x,y]) => n.map(n=>"--md-ref-palette-"+y+n+":"+hexFromArgb(t.palettes[x].tone(n)))).join(";") + "}"
	`, c)
}
