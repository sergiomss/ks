package main

import "github.com/sergiomss/ks/cmd"

var (
	version string
	date    string
	commit  string
)

func main() {
	cmd.Execute(version, date, commit)
}
