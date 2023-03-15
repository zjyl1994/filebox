package assets

import _ "embed"

//go:embed favicon.ico
var FaviconBytes []byte

//go:embed index.html
var IndexTemplateString string
