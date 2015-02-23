package item

import "github.com/egonelbre/fedwiki"

// just simple utility functions for creating items

func Paragraph(text string) fedwiki.Item {
	return fedwiki.Item{
		"type": "paragraph",
		"text": text,
	}
}

func HTML(text string) fedwiki.Item {
	return fedwiki.Item{
		"type": "html",
		"text": text,
	}
}

func Reference(title, site, text string) fedwiki.Item {
	return fedwiki.Item{
		"type":  "reference",
		"title": title,
		"site":  site,
		"text":  text,
	}
}

func Image(caption, url, text string) fedwiki.Item {
	return fedwiki.Item{
		"type":    "image",
		"url":     url,
		"text":    text,
		"caption": caption,
	}
}
