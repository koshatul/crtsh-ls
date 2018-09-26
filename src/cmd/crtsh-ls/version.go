package main

import "fmt"

var version = "dev"
var buildDate = "notset"
var gitHash = ""

func init() {
	rootCmd.Version = fmt.Sprintf("%s [%s] (%s)", version, gitHash, buildDate)
}
