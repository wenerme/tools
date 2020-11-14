//go:generate rice embed-go

package teleattrdata

import (
	"io/ioutil"
	"log"

	"github.com/wenerme/tools/pkg/teleattr"

	rice "github.com/GeertJohan/go.rice"
)

var instance *teleattr.PhoneData

func PhoneData() (*teleattr.PhoneData, error) {
	if instance != nil {
		return instance, nil
	}
	f, err := riceBox.Open("phone.dat")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	instance, err = teleattr.LoadBytes(b)
	return instance, err
}

var riceBox *rice.Box

func init() {
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateFS, rice.LocateWorkingDirectory, rice.LocateEmbedded, rice.LocateAppended},
	}
	box, err := conf.FindBox("../testdata")
	if err != nil {
		log.Fatalf("error opening rice.Box: %s\n", err)
	}
	riceBox = box
}
