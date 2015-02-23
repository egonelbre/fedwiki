# Server specification

## Content URLs

* `GET /` - returns the main page.
* `GET /favicon.png` - returns the favicon for the current site.
* `GET /systems/slugs.json` - returns the list of slugs on the site.
* `GET /systems/sitemap.json` - returns the sitemap for the site.

## Slug

Each page is uniquely identified by a `slug`.
Slugs need to match regular expression `[a-z0-9\-_\?\&\=\%\#\";\/]+`.

## Page URLs

All page urls can be suffixed with `.json` to return JSON content.

* `GET /<slug>` - returns page content

* `PUT /<slug>` - creates a new page with specified slug. The title must use the specified conversion to slug, otherwise the request will be rejected. Request data:

```
{
	"slug": "/<slug>"
	"title": "<title>",
	"synopsis": "<synopsis>",
	"version": 0,
	"story": [],
	"journal": []
}
```

* `DELETE /<slug>` - deletes the page with specified slug. If the request data `slug` or `version` mismatch the data in the page store the request will be rejected. Request data:

```
{
	"slug": "<slug>",
	"version": <version>
}
```

* `PATCH /<slug>` - updates the page with specified slug. If the request data `version` mismatches the server may either do a merge or reject the request. Request data:

```	
{
	"slug": "<slug>"
	"version": <version>
	"action": <action data>
}
```