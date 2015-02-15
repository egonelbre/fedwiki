package plugin

type Plugin interface {
}

// simple client-side plugin
type Client struct {
	Name string
	*folderstore
}
