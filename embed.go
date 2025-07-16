package embed

import "embed"

//go:embed client/dist/**
var ViteFS embed.FS

//go:embed assets/**
var AssetsFS embed.FS
