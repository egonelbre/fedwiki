package page

type Story []Item

type Item map[string]interface{}

func (item Item) Val(key string) string {
	if s, ok := item[key].(string); ok {
		return s
	}
	return ""
}

func (item Item) Id() string { return item.Val("id") }
