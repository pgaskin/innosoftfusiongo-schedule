//go:build ignore

package main

import (
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	asap := url.Values{
		"family": {"Asap:wdth,wght@87.5,100..900"},
		"text":   {""},
	}
	for _, r := range [][2]rune{
		{33, 126},
		{'\u2002', '\u201E'},
		{'\u2022', '\u2022'},
		{'\u2026', '\u2026'},
	} {
		for c := r[0]; c <= r[1]; c++ {
			asap["text"][0] += string(c)
		}
	}
	font("asap.woff2", asap)

	symbols := url.Values{
		"family": {"Material Symbols Outlined:opsz,wght,FILL,GRAD@20,300,0,0"},
		"text": {string([]rune{
			'\uE192', // schedule
			'\uE55F', // location_on
		})},
	}
	font("symbols.woff2", symbols)

	// note: use https://wakamaifondue.com/ to see font details
}

func font(name string, params url.Values) {
	css, err := css2(params)
	if err != nil {
		slog.Error("failed to fetch font css", "error", err)
		os.Exit(1)
	}

	if n := strings.Count(css, "@font-face"); n != 1 {
		slog.Error("expected a single variable font-face declaration, got "+strconv.Itoa(n), "css", css)
		os.Exit(1)
	}

	var u string
	if m := regexp.MustCompile(`url\((.+?)\)`).FindStringSubmatch(css); m == nil {
		slog.Error("failed to extract font url", "css", css)
		os.Exit(1)
	} else {
		u = m[1]
	}

	buf, err := woff2(u)
	if err != nil {
		slog.Error("failed to fetch font file", "error", err)
		os.Exit(1)
	}

	if err := os.WriteFile(name, buf, 0644); err != nil {
		slog.Error("failed to write font file", "error", err)
		os.Exit(1)
	}

	slog.Info("done", "name", name, "size", len(buf))
}

func css2(p url.Values) (string, error) {
	slog.Info("fetching font css", "params", p.Encode())

	req, err := http.NewRequest(http.MethodGet, "https://fonts.googleapis.com/css2?"+p.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0") // supports variable woff2

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status %d (%s)", resp.StatusCode, resp.Status)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func woff2(u string) ([]byte, error) {
	slog.Info("fetching font woff2", "url", u)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d (%s)", resp.StatusCode, resp.Status)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		if mt, _, _ := mime.ParseMediaType(ct); mt != "font/woff2" {
			return nil, fmt.Errorf("unexpected mimetype %q", ct)
		}
	}

	return io.ReadAll(resp.Body)
}
