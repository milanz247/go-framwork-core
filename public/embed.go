package public

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets
var assets embed.FS

func Handler() http.Handler { f, _ := fs.Sub(assets, "assets"); return http.FileServer(http.FS(f)) }
