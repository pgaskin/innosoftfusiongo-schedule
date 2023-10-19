package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		slog.Error("failed to create temp dir", "error", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmp)

	if err := npm(tmp, "@material/material-color-utilities", "0.2.7"); err != nil {
		slog.Error("failed to download package", "error", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(tmp, "bundle.js"), []byte(`
		import * as mcu from '@material/material-color-utilities'
		Object.assign(globalThis, mcu)
	`), 0644); err != nil {
		slog.Error("failed to write entry source", "error", err)
		os.Exit(1)
	}

	slog.Info("building")
	result := api.Build(api.BuildOptions{
		LogLevel:          api.LogLevelError,
		AbsWorkingDir:     tmp,
		EntryPoints:       []string{"bundle.js"},
		Outfile:           "bundle.dist.js",
		Bundle:            true,
		Write:             true,
		Drop:              api.DropConsole | api.DropDebugger,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		LegalComments:     api.LegalCommentsEndOfFile,
		Sourcemap:         api.SourceMapNone,
		Target:            api.ES2015,
		Platform:          api.PlatformNeutral,
		Format:            api.FormatIIFE,
	})
	if len(result.Errors) != 0 {
		slog.Error("failed to build package")
		os.Exit(1)
	}

	if buf, err := os.ReadFile(filepath.Join(tmp, "bundle.dist.js")); err != nil {
		slog.Error("failed to read output", "error", err)
		os.Exit(1)
	} else if err := os.WriteFile("mcu.js", buf, 0644); err != nil {
		slog.Error("failed to write output", "error", err)
		os.Exit(1)
	}
	slog.Info("done")
}

func npm(dest, name, version string) error {
	dest = filepath.Join(dest, "node_modules", filepath.FromSlash(name))
	slog.Info("downloading", "dest", dest, slog.Group("package", "name", name, "version", version))

	tarball, err := npmResolve(name, version)
	if err != nil {
		return fmt.Errorf("resolve package: %w", err)
	}

	slog.Info("using tarball", "url", tarball)

	if err := npmFetch(dest, tarball); err != nil {
		return fmt.Errorf("fetch package: %w", err)
	}
	return nil
}

func npmResolve(name, version string) (string, error) {
	resp, err := http.Get("https://registry.npmjs.org/" + name)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status %d (%s)", resp.StatusCode, resp.Status)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		if mt, _, _ := mime.ParseMediaType(ct); mt != "application/json" {
			return "", fmt.Errorf("expected application/json, got %q", mt)
		}
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var obj struct {
		DistTags struct {
			Latest string `json:"latest"`
		} `json:"dist-tags"`
		Versions map[string]struct {
			Dist struct {
				Tarball string `json:"tarball"`
			} `json:"dist"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(buf, &obj); err != nil {
		return "", err
	}
	if version == "" {
		if obj.DistTags.Latest == "" {
			return "", fmt.Errorf("failed to get latest version")
		}
		version = obj.DistTags.Latest
	}
	if obj.Versions[version].Dist.Tarball == "" {
		return "", fmt.Errorf("latest version has no tarball link")
	}
	return obj.Versions[version].Dist.Tarball, nil
}

func npmFetch(dest, tarball string) error {
	resp, err := http.Get(tarball)
	if err != nil {
		return fmt.Errorf("download package: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download package: response status %d (%s)", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("download package: %w", err)
	}
	defer zr.Close()

	for tr := tar.NewReader(zr); ; {
		h, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("download package: %w", err)
		}

		if h.Typeflag != tar.TypeReg {
			continue
		}

		n, ok := strings.CutPrefix(h.Name, "package/")
		if !ok {
			return fmt.Errorf("download package: filename %q does not begin with package/", h.Name)
		}
		slog.Info("download", "name", n)

		if err := os.MkdirAll(filepath.Join(dest, filepath.FromSlash(path.Dir(n))), 0755); err != nil {
			return fmt.Errorf("download package: %w", err)
		}

		if f, err := os.OpenFile(filepath.Join(dest, filepath.FromSlash(n)), os.O_CREATE|os.O_TRUNC|os.O_EXCL|os.O_WRONLY, 0644); err != nil {
			return fmt.Errorf("download package: %w", err)
		} else if _, err := io.CopyN(f, tr, h.Size); err != nil {
			f.Close()
			return fmt.Errorf("download package: %w", err)
		} else if err := f.Close(); err != nil {
			return fmt.Errorf("download package: %w", err)
		}
	}
	return nil
}
