package embed

import "embed"

//go:embed assets/*
var AssetFS embed.FS

