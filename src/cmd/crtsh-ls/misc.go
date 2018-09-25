package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CertificateRecord contains the fields used in the JSON response from crt.sh queries
type CertificateRecord struct {
	IssuerCaID        int    `json:"issuer_ca_id"`
	IssuerName        string `json:"issuer_name"`
	NameValue         string `json:"name_value"`
	MinCertID         int    `json:"min_cert_id"`
	MinEntryTimestamp string `json:"min_entry_timestamp"`
	NotBefore         string `json:"not_before"`
	NotAfter          string `json:"not_after"`
}

func getCertStream(domain string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: viper.GetDuration("timeout"),
	}

	baseurl, err := url.Parse(viper.GetString("crtsh.base_uri"))
	if err != nil {
		logrus.Panic(err)
	}

	query := baseurl.Query()
	query.Add("output", "json")
	query.Add("q", domain)
	baseurl.RawQuery = query.Encode()

	logrus.Debugf("Requesting: %s", baseurl.String())

	resp, err := client.Get(baseurl.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned %d http status code (%s)", resp.StatusCode, resp.Status)
	}

	return resp.Body, err
}
