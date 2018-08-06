package main

import "github.com/majestic-fox/ks/cmd"

var (
	Version   string
	BuildDate string
)

func main() {
	cmd.Execute(Version, BuildDate)
}
