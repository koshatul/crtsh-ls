package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

// basicFunctions are the set of initial
// functions provided to every template.
// nolint: gochecknoglobals // it's copypasta
var basicFunctions = template.FuncMap{
	"json": func(v interface{}) string {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(v)
		if err != nil {
			logrus.Fatal(err)
		}
		// Remove the trailing new line added by the encoder
		return strings.TrimSpace(buf.String())
	},
	"split":    strings.Split,
	"join":     strings.Join,
	"title":    strings.Title,
	"lower":    strings.ToLower,
	"upper":    strings.ToUpper,
	"pad":      padWithSpace,
	"padlen":   padToLength,
	"padmax":   padToMaxLength,
	"truncate": truncateWithLength,
	"tf":       stringTrueFalse,
	"yn":       stringYesNo,
	"t":        stringTab,
}

// nolint: gochecknoglobals // it's copypasta
var maxLenPrefix = 0

// padToLength adds whitespace to pad to the supplied length
func padToMaxLength(source string) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", maxLenPrefix), source)
}

// padToLength adds whitespace to pad to the supplied length
func padToLength(source string, prefix int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", prefix), source)
}

// padWithSpace adds whitespace to the input if the input is non-empty
func padWithSpace(source string, prefix, suffix int) string {
	if source == "" {
		return source
	}

	return strings.Repeat(" ", prefix) + source + strings.Repeat(" ", suffix)
}

// stringTrueFalse returns "true" or "false" for boolean input
func stringTrueFalse(source bool) string {
	if source {
		return "true"
	}

	return "false"
}

// stringYesNo returns "yes" or "no" for boolean input
func stringYesNo(source bool) string {
	if source {
		return "yes"
	}

	return "no"
}

// stringTab returns a tab character
func stringTab() string {
	return "\t"
}

// truncateWithLength truncates the source string up to the length provided by the input
func truncateWithLength(source string, length int) string {
	if len(source) < length {
		return source
	}

	return source[:length]
}
