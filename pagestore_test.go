package fedwiki

import "testing"

var slugcases = []struct {
	In  string
	Exp Slug
}{
	{In: "", Exp: "-"},
	{In: "Hello  World 90", Exp: "hello-world-90"},
	{In: "Hello, 世界", Exp: "hello-世界"},
	{In: "90Things", Exp: "90things"},
	{In: "90 Things", Exp: "90-things"},
	{In: "KÜSIMUSED", Exp: "küsimused"},
	{In: "Küsimused Öösel", Exp: "küsimused-öösel"},
	{In: "nested / _paths", Exp: "nested/paths"},
	{In: "nested-/-paths", Exp: "nested/paths"},
	{In: "example_test.go", Exp: "example-test-go"},
}

func TestSlugify(t *testing.T) {
	for _, test := range slugcases {
		got := Slugify(test.In)
		if got != test.Exp {
			t.Errorf("Slugify(%q): got %q expected %q", test.In, got, test.Exp)
		}
	}
}

func TestValidateSlug(t *testing.T) {
	for _, test := range slugcases {
		err := ValidateSlug(test.Exp)
		if err != nil {
			t.Errorf("Invalid %q: %v", test.Exp, err)
		}
	}
}
