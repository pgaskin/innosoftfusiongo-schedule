//go:build ignore

package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
)

func main() {
	var (
		u = index(os.Args, 1)
		n = index(os.Args, 2)
	)
	if err := fetch(u, n); err != nil {
		slog.Error("failed to fetch", "url", u, "error", err)
		os.Exit(1)
	}
}

func fetch(u, n string) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	if n == "" {
		n = path.Base(req.URL.Path)
	}
	if n == "" {
		return errors.New("no filename")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(n, buf, 0666); err != nil {
		return err
	}

	slog.Info("fetched", "url", u, "name", n, "size", len(buf))
	return nil
}

func index[T any](a []T, i int) T {
	if i < len(a) {
		return a[i]
	}
	var z T
	return z
}
