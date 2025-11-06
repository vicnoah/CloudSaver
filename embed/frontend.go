package embed

import "embed"

//go:embed all:dist
var FrontendAssets embed.FS
