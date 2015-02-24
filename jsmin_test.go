package jsmin_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/omeid/jsmin"
)

func TestMinify(t *testing.T) {
	files, err := ioutil.ReadDir(`testdata`)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		name := filepath.Join("testdata", file.Name())
		if !strings.HasSuffix(name, ".before") {
			continue
		}
		before, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		after, err := ioutil.ReadFile(name[:len(name)-7] + ".after")
		if err != nil {
			t.Fatal(err)
		}
		reader, err := jsmin.Minify(bytes.NewReader(before))
		if err != nil {
			t.Fatal(err)
		}

		out, _ := ioutil.ReadAll(reader)
		if string(out) != string(after) {
			println("---------------")
			println(string(out))
			println("---------------")
			println(string(after))
			println("---------------")
			t.Fatal("Whops. We failed.")
		}
	}
}
