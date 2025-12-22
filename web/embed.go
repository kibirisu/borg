package web

import (
	"embed"
	"io/fs"
)

//go:embed dist
var assets embed.FS

func GetAssets() fs.FS {
	res, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}
	return res
}
