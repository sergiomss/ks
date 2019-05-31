package main

import "github.com/sergiomss/ks/cmd"

var (
	Version   string
	BuildDate string
)

func main() {
	cmd.Execute(Version, BuildDate)
}
