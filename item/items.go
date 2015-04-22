// This package implements common federated wiki types
package item

import "github.com/egonelbre/fedwiki"

func Paragraph(text string) fedwiki.Item {
	return fedwiki.Item{
		"type": "paragraph",
		"id":   fedwiki.NewID(),
		"text": text,
	}
}

func HTML(text string) fedwiki.Item {
	return fedwiki.Item{
		"type": "html",
		"id":   fedwiki.NewID(),
		"text": text,
	}
}

func Reference(title, site, text string) fedwiki.Item {
	return fedwiki.Item{
		"type":  "reference",
		"id":    fedwiki.NewID(),
		"title": title,
		"site":  site,
		"text":  text,
	}
}

func Image(caption, url, text string) fedwiki.Item {
	return fedwiki.Item{
		"type":    "image",
		"id":      fedwiki.NewID(),
		"url":     url,
		"text":    text,
		"caption": caption,
	}
}
