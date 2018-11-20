package main

import (
	"github.com/namsral/flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.WarnLevel)

	p := &PublishAPK{}

	flag.StringVar(&p.apkPath, "apk", "", "path to the apk which should be uploaded.")
	flag.StringVar(&p.pkgName, "pkg", "", "package name of the application.")
	flag.StringVar(&p.email, "email", "", "email address of the service account.")
	flag.StringVar(&p.keyPath, "key", "", "path to the key.pem file for the service account OR base64 encoded private key.")
	flag.StringVar(&p.track, "track", "", "track to which the apk should be pushed (e.g.: alpha, beta, internal, production)")
	flag.BoolVar(&p.debug, "debug", false, "flag indicating whether debug output should be printed.")
	flag.Parse()

	if p.apkPath == "" {
		log.Fatal("Missing 'apk'-parameter, please specify path to the apk you want to upload.")
	}

	if p.pkgName == "" {
		log.Fatal("Missing 'pkg'-parameter, please specify the name of the package.")
	}

	if p.email == "" {
		log.Fatal("Missing 'email'-parameter, please specify the email address of your service account.")
	}

	if p.keyPath == "" {
		log.Fatal("Missing 'key'-parameter, please specify the path to the key.pem file for the service account.")
	}

	if p.track == "" {
		log.Fatal("Missing 'track'-parameter, please specify to which track you want to deploy the apk.")
	}

	if p.debug {
		log.SetLevel(log.DebugLevel)
	}

	err := p.Init()
	if err != nil {
		log.Fatalf("Couldn't initialize PublishAPK service, see: %v", err)
	}

	err = p.Upload()
	if err != nil {
		log.Fatalf("Couldn't upload apk to track, see: %v", err)
	}
}
