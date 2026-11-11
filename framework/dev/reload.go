// Package dev supplies development-only browser reload endpoints.
package dev

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"
)

type Reload struct {
	roots []string
	boot  int64
}

func NewReload(roots ...string) *Reload {
	return &Reload{roots: roots, boot: time.Now().UnixNano()}
}

func (reload *Reload) Version(w http.ResponseWriter, r *http.Request) {
	latest := reload.boot
	for _, root := range reload.roots {
		fs.WalkDir(os.DirFS("."), root, func(path string, entry fs.DirEntry, err error) error {
			if err == nil && !entry.IsDir() {
				if info, err := entry.Info(); err == nil && info.ModTime().UnixNano() > latest {
					latest = info.ModTime().UnixNano()
				}
			}
			return nil
		})
	}
	fmt.Fprint(w, latest)
}

func (reload *Reload) Script(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte(`let version;setInterval(async()=>{try{const r=await fetch('/dev/version',{cache:'no-store'});if(!r.ok)return;const next=await r.text();if(version&&version!==next)location.reload();version=next;}catch{}},1500);`))
}
