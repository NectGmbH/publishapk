package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	publisher "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
)

// PublishAPK represents the apk publisher.
type PublishAPK struct {
	apkPath string
	pkgName string
	email   string
	keyPath string
	track   string
	debug   bool

	svc *publisher.Service
}

// Init initializes the current PublishAPK service.
func (p *PublishAPK) Init() error {
	var key []byte

	if _, err := os.Stat(p.keyPath); os.IsNotExist(err) {
		key, err = base64.StdEncoding.DecodeString(p.keyPath)
		if err != nil {
			return fmt.Errorf("`%s` is neither a file nor a base64 encoded private key, see decode error: %v", p.keyPath, err)
		}
	} else {
		key, err = ioutil.ReadFile(p.keyPath)
		if err != nil {
			return fmt.Errorf("couldn't read key from `%s`, see: %v", p.keyPath, err)
		}
	}

	conf := &jwt.Config{
		Email:      p.email,
		PrivateKey: key,
		Scopes:     []string{publisher.AndroidpublisherScope},
		TokenURL:   google.JWTTokenURL,
	}

	client := conf.Client(oauth2.NoContext)

	svc, err := publisher.New(client)
	if err != nil {
		return fmt.Errorf("couldn't create new publisher, see: %v", err)
	}

	p.svc = svc

	return nil
}

// Upload uploads the apk to the specified track.
func (p *PublishAPK) Upload() error {
	editSvc := p.svc.Edits

	insert := editSvc.Insert(p.pkgName, nil)
	edit, err := insert.Do()
	if err != nil {
		return fmt.Errorf("couldn't create insert, see: %v", err)
	}
	log.Debugf("created new edit with id %s", edit.Id)

	apkReader, err := os.Open(p.apkPath)
	if err != nil {
		return fmt.Errorf("couldn't open APK from `%s`, see: %v", p.apkPath, err)
	}
	defer apkReader.Close()

	upload := editSvc.Apks.Upload(p.pkgName, edit.Id)
	upload.Media(apkReader, googleapi.ContentType("application/vnd.android.package-archive"))
	apk, err := upload.Do()
	if err != nil {
		return fmt.Errorf("couldn't upload APK, see: %v", err)
	}

	log.Debugf("uploaded apk, new version code: %d", apk.VersionCode)

	trackChange := &publisher.Track{
		Releases: []*publisher.TrackRelease{
			{
				VersionCodes: googleapi.Int64s{apk.VersionCode},
				ReleaseNotes: []*publisher.LocalizedText{},
				Status:       "completed",
			},
		},
	}

	updateReq := editSvc.Tracks.Update(p.pkgName, edit.Id, p.track, trackChange)
	_, err = updateReq.Do()
	if err != nil {
		return fmt.Errorf("couldn't update track `%s`, see: %v", p.track, err)
	}

	committedEdit, err := editSvc.Commit(p.pkgName, edit.Id).Do()
	if err != nil {
		return fmt.Errorf("couldn't commit change for edit id `%s`, see: %v", committedEdit.Id, err)
	}

	log.Debugf("committed app edit with id %s", committedEdit.Id)

	return nil
}
