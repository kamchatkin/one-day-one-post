package static

import _ "embed"

//go:embed index.html
var IndexPage []byte

//go:embed style.css
var StylePage string

//go:embed style.css.map
var StyleMapPage string

//go:embed favicon.png
var FaviconImage []byte
