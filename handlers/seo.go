package handlers

import (
	"fmt"
	"net/http"
)

var domain = "https://customize.fly.dev"

func RobotsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, `User-agent: *
Allow: /
Allow: /dashboard
Allow: /basic-api-examples

Disallow: /api/v1/
Disallow: /api/v1/examples/
Disallow: /api/v1/custom/
Disallow: /save
Disallow: /health

Sitemap: `+domain+`/sitemap.xml`)
}

func SitemapHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>`+domain+`</loc>
    </url>
    <url>
        <loc>`+domain+`/dashboard</loc>
    </url>
    <url>
        <loc>`+domain+`/basic-api-examples</loc>
    </url>
</urlset>`)
}
