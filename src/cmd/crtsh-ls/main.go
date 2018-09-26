package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

var rootCmd = &cobra.Command{
	Use:   "crtsh-ls <domain>",
	Short: "crtsh-ls lists domains from crt.sh database",
	Long:  "<domain> is \"%.github.com\" to show all subdomains of github.com or \"github.com\" to show single domain certificate",
	Args:  cobra.MinimumNArgs(1),
	Run:   mainCommand,
}

func init() {
	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().StringP("format", "f", "{{padlen .NameValue 20}}\t{{.NotBefore}}\t{{.NotAfter}}", "Output formatting (go template).\n Possible items are IssuerCaID, IssuerName, NameValue, MinCertID, MinEntryTimestamp, NotBefore, NotAfter.\n")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output. ")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 50*time.Second, "Request timeout.")
	rootCmd.PersistentFlags().Bool("only-valid", false, "Only display still (date) valid certificates.")

	err := multierr.Combine(
		viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format")),
		viper.BindEnv("format", "FORMAT"),
		viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout")),
		viper.BindEnv("timeout", "TIMEOUT"),
		viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")),
		viper.BindEnv("debug", "DEBUG"),
		viper.BindPFlag("only-valid", rootCmd.PersistentFlags().Lookup("only-valid")),
		viper.BindEnv("only-valid", "ONLY_VALID"),
	)
	if err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func mainCommand(cmd *cobra.Command, args []string) {
	if !strings.HasSuffix(viper.GetString("format"), "\n") {
		viper.Set("format", fmt.Sprintf("%s\n", viper.GetString("format")))
	}
	tmpl, err := template.New("").Funcs(basicFunctions).Parse(viper.GetString("format"))
	if err != nil {
		logrus.Fatal(err)
	}

	data, err := getCertStream(args[0])
	if err != nil {
		logrus.Fatal(err)
	}

	dec := json.NewDecoder(data)

	for {
		cert := CertificateRecord{}
		if err := dec.Decode(&cert); err == io.EOF {
			break
		} else if err != nil {
			logrus.Fatal(err)
		}
		if viper.GetBool("only-valid") {
			t, err := time.Parse("2006-01-02T15:04:05", cert.NotAfter)
			if err != nil {
				logrus.Debugf("Failed to parse time: %s (%s)", cert.NotAfter, err.Error())
				continue
			}
			if t.Before(time.Now()) {
				continue
			}
		}
		if err := tmpl.Execute(os.Stdout, cert); err != nil {
			logrus.Warnf("Unable to format line: %s", err.Error())
		}
	}

}
