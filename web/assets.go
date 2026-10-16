package web

import "embed"

// Assets embeds all web front-end files into the binary.
//
//go:embed *
var Assets embed.FS
