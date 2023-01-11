package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/urfave/cli/v2"
)

var (
	keyPairID string
	rawPriKey string
)

func main() {
	pKey, err := parseSignPrivateKey(rawPriKey)
	if err != nil {
		panic(err)
	}
	signer := sign.NewURLSigner(keyPairID, pKey)
	app := buildCli(signer)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func buildCli(signer *sign.URLSigner) *cli.App {
	var (
		ttl     uint
		ip      string
		startAt string
		endAt   string
	)

	app := &cli.App{
		Name:                 "CloudFront URL Signer",
		Usage:                "sign your provide raw url",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "policy",
				Usage: "Sign URL by Custom Policy",
				Action: func(ctx *cli.Context) error {
					urlArg := ctx.Args().Get(0)
					if len(urlArg) == 0 {
						return fmt.Errorf("you have to provide the url for signing")
					}

					u, err := parseURL(urlArg)
					if err != nil {
						return err
					}

					lt, err := time.Parse(time.RFC3339, startAt)
					if err != nil {
						return err
					}

					gt, err := time.Parse(time.RFC3339, endAt)
					if err != nil {
						return err
					}

					p := getPolicy(u.String(), lt, gt, ip)
					s, err := signer.SignWithPolicy(u.String(), p)
					if err != nil {
						return err
					}
					fmt.Printf("Signed URL: \n")
					fmt.Println(s)
					return nil

				},
			},
			{
				Name:  "time",
				Usage: "Sign URL by TTL",
				Action: func(ctx *cli.Context) error {
					urlArg := ctx.Args().Get(0)
					if len(urlArg) == 0 {
						return fmt.Errorf("you have to provide the url for signing")
					}
					u, err := parseURL(urlArg)
					if err != nil {
						return err
					}

					policy := sign.NewCannedPolicy(u.String(), time.Now().Add(time.Second*time.Duration(ttl)))
					s, err := signer.SignWithPolicy(u.String(), policy)
					if err != nil {
						return err
					}
					fmt.Printf("Signed URL: \n")
					fmt.Println(s)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:        "ttl",
				Usage:       "Provide the URL time-to-live in second",
				Destination: &ttl,
			},
			&cli.StringFlag{
				Name:        "start",
				Usage:       "Provide the URL be valid from. Valid in Policy.",
				Destination: &startAt,
			},
			&cli.StringFlag{
				Name:        "end",
				Usage:       "Provide the URL be valid to. Valid in Policy.",
				Destination: &endAt,
			},
			&cli.StringFlag{
				Name:        "ip",
				Usage:       "Allowed IP. Valid in Policy.",
				Destination: &ip,
			},
		},
	}
	return app
}

func parseURL(urlString string) (*url.URL, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	fmt.Printf("URL Scheme: %s\n", u.Scheme)
	fmt.Printf("URL Host: %s\n", u.Host)
	fmt.Printf("URL Path: %s\n", u.Path)
	fmt.Printf("URL Query: %s\n", u.RawQuery)
	fmt.Println("")
	return u, nil
}

func getPolicy(url string, lsThan time.Time, gtThan time.Time, ip string) *sign.Policy {
	return &sign.Policy{
		Statements: []sign.Statement{
			{
				Resource: url,
				Condition: sign.Condition{
					DateLessThan:    sign.NewAWSEpochTime(lsThan),
					DateGreaterThan: sign.NewAWSEpochTime(gtThan),
					IPAddress:       &sign.IPAddress{SourceIP: ip},
				},
			},
		},
	}
}

func parseSignPrivateKey(key string) (*rsa.PrivateKey, error) {
	kb, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(kb)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("private key parsing error")
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}
