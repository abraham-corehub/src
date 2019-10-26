package main

import (
	"fmt"

	"github.com/qor/assetfs"
)

func main() {
	// Default implemention based on filesystem, you could overwrite with other implemention, for example bindatafs will do this.
	aFS := assetfs.AssetFS()

	// Register path to AssetFS
	aFS.RegisterPath("web/app/views")

	// Prepend path to AssetFS
	aFS.PrependPath("web/app/views")

	// Get file's content with name from path `/web/app/views`
	b, err := aFS.Asset("filename.tmpl")

	fmt.Println(string(b), err)

	// List matched files from assetfs
	lS, err := aFS.Glob("*.tmpl")

	for _, s := range lS {
		fmt.Println(s)
	}

	// NameSpace return namespaced filesystem
	namespacedFS := aFS.NameSpace("asset")
	err = namespacedFS.RegisterPath("web/app/myviews")
	if err != nil {
		fmt.Println(err)
	}
	err = namespacedFS.PrependPath("web/app/myviews")
	if err != nil {
		fmt.Println(err)
	}
	// Will lookup file with name "filename.tmpl" from path `/web/app/myspecialviews` but not `/web/app/views`
	b, err = namespacedFS.Asset("filename.tmpl")
	fmt.Println(string(b), err)
	namespacedFS.Glob("*.tmpl")
}
