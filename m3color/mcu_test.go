package m3color

import (
	"testing"
)

func TestEval(t *testing.T) {
	if v, err := eval[int64](`c`, `return argbFromHex(c)`, "#11223344"); err != nil {
		t.Errorf("m3color: failed to test: %v", err)
	} else if v != 4280431428 {
		t.Errorf("m3color: failed to test: incorrect result")
	}
}

func BenchmarkEval(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := eval[int64](`c`, `return argbFromHex(c)`, "#11223344"); err != nil {
			panic(err)
		}
	}
}

func TestPaletteCSS(t *testing.T) {
	t.Run("0074a4", func(t *testing.T) {
		act, err := PaletteCSS("0074a4")
		if err != nil {
			t.Fatal(err)
		}

		exp := `:root{--md-source:#0074a4;--md-ref-palette-primary0:#000000;--md-ref-palette-primary4:#00101b;--md-ref-palette-primary5:#00131f;--md-ref-palette-primary6:#001522;--md-ref-palette-primary10:#001e2e;--md-ref-palette-primary12:#002234;--md-ref-palette-primary17:#002d43;--md-ref-palette-primary20:#00344c;--md-ref-palette-primary22:#003953;--md-ref-palette-primary24:#003d59;--md-ref-palette-primary25:#00405c;--md-ref-palette-primary30:#004c6d;--md-ref-palette-primary35:#00587e;--md-ref-palette-primary40:#00658f;--md-ref-palette-primary50:#007fb3;--md-ref-palette-primary60:#3199d1;--md-ref-palette-primary70:#54b4ed;--md-ref-palette-primary80:#86ceff;--md-ref-palette-primary87:#b5dfff;--md-ref-palette-primary90:#c8e6ff;--md-ref-palette-primary92:#d3ebff;--md-ref-palette-primary94:#dff0ff;--md-ref-palette-primary95:#e5f2ff;--md-ref-palette-primary96:#ebf5ff;--md-ref-palette-primary98:#f6faff;--md-ref-palette-primary99:#fbfcff;--md-ref-palette-primary100:#ffffff;--md-ref-palette-secondary0:#000000;--md-ref-palette-secondary4:#00101b;--md-ref-palette-secondary5:#01131e;--md-ref-palette-secondary6:#031520;--md-ref-palette-secondary10:#0b1d29;--md-ref-palette-secondary12:#10212d;--md-ref-palette-secondary17:#1b2c38;--md-ref-palette-secondary20:#21323e;--md-ref-palette-secondary22:#263743;--md-ref-palette-secondary24:#2a3b48;--md-ref-palette-secondary25:#2c3d4a;--md-ref-palette-secondary30:#384956;--md-ref-palette-secondary35:#435562;--md-ref-palette-secondary40:#4f616e;--md-ref-palette-secondary50:#687987;--md-ref-palette-secondary60:#8193a1;--md-ref-palette-secondary70:#9baebd;--md-ref-palette-secondary80:#b6c9d8;--md-ref-palette-secondary87:#cadcec;--md-ref-palette-secondary90:#d2e5f5;--md-ref-palette-secondary92:#d8ebfb;--md-ref-palette-secondary94:#dff0ff;--md-ref-palette-secondary95:#e5f2ff;--md-ref-palette-secondary96:#ebf5ff;--md-ref-palette-secondary98:#f6faff;--md-ref-palette-secondary99:#fbfcff;--md-ref-palette-secondary100:#ffffff;--md-ref-palette-tertiary0:#000000;--md-ref-palette-tertiary4:#110827;--md-ref-palette-tertiary5:#140b2a;--md-ref-palette-tertiary6:#160e2c;--md-ref-palette-tertiary10:#1f1635;--md-ref-palette-tertiary12:#231a39;--md-ref-palette-tertiary17:#2e2544;--md-ref-palette-tertiary20:#342b4b;--md-ref-palette-tertiary22:#393050;--md-ref-palette-tertiary24:#3d3454;--md-ref-palette-tertiary25:#3f3657;--md-ref-palette-tertiary30:#4b4263;--md-ref-palette-tertiary35:#574d6f;--md-ref-palette-tertiary40:#63597c;--md-ref-palette-tertiary50:#7c7196;--md-ref-palette-tertiary60:#968bb1;--md-ref-palette-tertiary70:#b1a5cc;--md-ref-palette-tertiary80:#cdc0e9;--md-ref-palette-tertiary87:#e1d4fd;--md-ref-palette-tertiary90:#e9ddff;--md-ref-palette-tertiary92:#eee4ff;--md-ref-palette-tertiary94:#f3eaff;--md-ref-palette-tertiary95:#f6edff;--md-ref-palette-tertiary96:#f8f1ff;--md-ref-palette-tertiary98:#fdf7ff;--md-ref-palette-tertiary99:#fffbff;--md-ref-palette-tertiary100:#ffffff;--md-ref-palette-neutral0:#000000;--md-ref-palette-neutral4:#0c0e11;--md-ref-palette-neutral5:#0f1113;--md-ref-palette-neutral6:#111416;--md-ref-palette-neutral10:#191c1e;--md-ref-palette-neutral12:#1d2022;--md-ref-palette-neutral17:#282a2c;--md-ref-palette-neutral20:#2e3133;--md-ref-palette-neutral22:#333537;--md-ref-palette-neutral24:#37393c;--md-ref-palette-neutral25:#393c3e;--md-ref-palette-neutral30:#454749;--md-ref-palette-neutral35:#505355;--md-ref-palette-neutral40:#5d5e61;--md-ref-palette-neutral50:#75777a;--md-ref-palette-neutral60:#8f9193;--md-ref-palette-neutral70:#aaabae;--md-ref-palette-neutral80:#c5c6c9;--md-ref-palette-neutral87:#d9dadd;--md-ref-palette-neutral90:#e2e2e5;--md-ref-palette-neutral92:#e7e8eb;--md-ref-palette-neutral94:#edeef0;--md-ref-palette-neutral95:#f0f0f3;--md-ref-palette-neutral96:#f3f3f6;--md-ref-palette-neutral98:#f9f9fc;--md-ref-palette-neutral99:#fcfcff;--md-ref-palette-neutral100:#ffffff;--md-ref-palette-neutral-variant0:#000000;--md-ref-palette-neutral-variant4:#090f14;--md-ref-palette-neutral-variant5:#0b1216;--md-ref-palette-neutral-variant6:#0e1419;--md-ref-palette-neutral-variant10:#161c21;--md-ref-palette-neutral-variant12:#1a2025;--md-ref-palette-neutral-variant17:#242b30;--md-ref-palette-neutral-variant20:#2b3136;--md-ref-palette-neutral-variant22:#2f363b;--md-ref-palette-neutral-variant24:#333a3f;--md-ref-palette-neutral-variant25:#363c42;--md-ref-palette-neutral-variant30:#41484d;--md-ref-palette-neutral-variant35:#4d5359;--md-ref-palette-neutral-variant40:#595f65;--md-ref-palette-neutral-variant50:#71787e;--md-ref-palette-neutral-variant60:#8b9198;--md-ref-palette-neutral-variant70:#a6acb2;--md-ref-palette-neutral-variant80:#c1c7ce;--md-ref-palette-neutral-variant87:#d4dbe1;--md-ref-palette-neutral-variant90:#dde3ea;--md-ref-palette-neutral-variant92:#e3e9f0;--md-ref-palette-neutral-variant94:#e8eef5;--md-ref-palette-neutral-variant95:#ebf1f8;--md-ref-palette-neutral-variant96:#eef4fb;--md-ref-palette-neutral-variant98:#f6faff;--md-ref-palette-neutral-variant99:#fbfcff;--md-ref-palette-neutral-variant100:#ffffff;--md-ref-palette-error0:#000000;--md-ref-palette-error4:#280001;--md-ref-palette-error5:#2d0001;--md-ref-palette-error6:#310001;--md-ref-palette-error10:#410002;--md-ref-palette-error12:#490002;--md-ref-palette-error17:#5c0004;--md-ref-palette-error20:#690005;--md-ref-palette-error22:#710005;--md-ref-palette-error24:#790006;--md-ref-palette-error25:#7e0007;--md-ref-palette-error30:#93000a;--md-ref-palette-error35:#a80710;--md-ref-palette-error40:#ba1a1a;--md-ref-palette-error50:#de3730;--md-ref-palette-error60:#ff5449;--md-ref-palette-error70:#ff897d;--md-ref-palette-error80:#ffb4ab;--md-ref-palette-error87:#ffcfc9;--md-ref-palette-error90:#ffdad6;--md-ref-palette-error92:#ffe2de;--md-ref-palette-error94:#ffe9e6;--md-ref-palette-error95:#ffedea;--md-ref-palette-error96:#fff0ee;--md-ref-palette-error98:#fff8f7;--md-ref-palette-error99:#fffbff;--md-ref-palette-error100:#ffffff}`
		if act != exp {
			t.Fatalf("incorrect result:\n\texp:%s\n\tact:%s", exp, act)
		}
	})
}

func BenchmarkPaletteCSS(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := PaletteCSS("0074a4"); err != nil {
			panic(err)
		}
	}
}
