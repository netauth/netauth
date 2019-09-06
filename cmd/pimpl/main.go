package main

import (
	"github.com/NetAuth/NetAuth/pkg/plugin/tree"
)

func main() {
	tree.PluginMain(stfu{tree.NullPlugin{}})
}

type stfu struct {
	tree.NullPlugin
}
