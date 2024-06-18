package page

import (
	"fmt"
	"testing"
)

func TestPage(t *testing.T) {

	var content []string
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")
	content = append(content, "1")

	page := NewPage(content, 1000, nil)

	fmt.Println(page)

}
