package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/isi-nc/sharepoint2sync/internal"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Azure/go-ntlmssp"

	s2kClient "github.com/isi-nc/sync2kafka/client"
)

var (
	// sharepoint
	url,
	user,
	password string
	ntlmAuth,
	unescapeSp,
	unescapeHtml bool

	// sync2kafka
	useTls      = flag.Bool("use-tls", false, "use TLS connection")
	skipVerify  = flag.Bool("skip-tls-verify", false, "skip tls verification")
	tlsCertPath = flag.String("tls-cert", "", "TLS certificate path (required if key is set)")
	token       = flag.String("token", "", "sync2kafka server token")
	server      = flag.String("server", ":9084", "sync2kafka server url")
	topic       = flag.String("topic", "sync2kafka", "destination topic")
	doDelete    = flag.Bool("do-delete", false, "Instruct sync2kafka server to perform deletions too")

	s2klient *s2kClient.BinarySync2KafkaClient
)

func main() {
	fmt.Println("hey")
	flag.StringVar(&url, "url", "", "sharepoint API url to request")
	flag.BoolVar(&ntlmAuth, "ntlm-auth", false, "use NTLM for sharepoint server auth")
	flag.StringVar(&user, "user", "", "username for sharepoint")
	flag.StringVar(&password, "password", "", "password for sharepoint")
	flag.BoolVar(&unescapeSp, "unescapeSp", false, "unescape _xXXXX_ characters produced by sharepoint")
	flag.BoolVar(&unescapeHtml, "unescapeHtml", true, "unescape html escaping (%20 etc...)")

	flag.Set("logtostderr", "true")
	flag.Parse()

	log.Printf("Start sharepoint2kafka with user %s for url %s on topic %s, ntml-auth %t",
		user,
		url,
		*topic,
		ntlmAuth)

	SetupCloseHandler()

	// connect to sync2kafka
	log.Printf("connect to sync2kafka %s", *server)
	crt := ""
	if *useTls {
		crt = readCert()
	}
	s2klient = s2kClient.NewBinary(&s2kClient.SyncInitInfo{
		Format:   "",
		DoDelete: false,
		Token:    *token,
		Topic:    *topic,
	}, *server, *skipVerify, *useTls, crt)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err := s2klient.Connect(ctx)
	defer func() {
		err := s2klient.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	// connect to sharepoint
	log.Printf("connect to sharepoint %s", url)
	shpClient := &http.Client{}

	// if we use the NTLM auth we create the ntlm transport
	if ntlmAuth {
		log.Println("sharepoint - use ntlm transport")
		shpClient = &http.Client {
			Transport: ntlmssp.Negotiator {
				RoundTripper: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			},
		}
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("Accept", "application/json;odata=nometadata")
	req.Header.Add("Charset", "utf-8")

	if user != "" && password != "" {
		log.Println("sharepoint - set basic auth")
		req.SetBasicAuth(user, password)
	}

	resp, err := shpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		contents, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("sharepoint - response statusCode: %d, payload %s", resp.StatusCode, string(contents))
	}

	log.Printf("sharepoint - parsing entries %s", url)
	entries, err := internal.ParseJsonSharepointValues(resp.Body)
	log.Printf("sharepoint - read %d entries from sharepoint", len(entries))
	log.Printf("sync2kafka - start transfer to sync2kafka %s", url)
	if err = s2klient.StartTransfer(); err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		// process value
		value := entry.Value
		if unescapeSp {
			value = internal.ReplaceEscapedXml(value)
		}

		if unescapeHtml {
			value = []byte(html.UnescapeString(string(value)))
		}

		// send to sync2kafa
		kv := s2kClient.BinaryKV{
			Key:   []byte(entry.Id),
			Value: value,
		}

		if err = s2klient.SendValue(kv); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("sync2kafka - %d messages sent to sync2kafka", len(entries))

	log.Printf("sync2kafka - end transfer to sync2kafka %s", url)
	err = s2klient.EndTransfer()
	if err != nil {
		log.Fatal(err)
	}
}

func readCert() string {
	var crt string
	if tlsCertPath != nil {
		crtBytes, err := ioutil.ReadFile(*tlsCertPath)
		if err != nil {
			log.Fatal(err)
		}
		crt = string(crtBytes)
	}
	return crt
}


func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)
	go func() {
		<-c
		if s2klient != nil {
			if err := s2klient.Close(); err != nil {
				os.Exit(1)
			}
		}
		os.Exit(0)
	}()
}