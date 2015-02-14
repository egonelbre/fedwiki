package page

import "time"

type Header struct {
	Slug     Slug      `json:"slug"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Synopsis string    `json:"synopsis"`
}

type Page struct {
	Header  Header
	Story   Story   `json:"story"`
	Journal Journal `json:"journal"`
}
