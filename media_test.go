package onvif

import (
	"fmt"
	"log"
	"testing"
)

func TestGetProfiles(t *testing.T) {
	log.Println("Test GetProfiles")

	res, err := testDevice.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetStreamURI(t *testing.T) {
	var res MediaURI 

	log.Println("Test GetStreamURI")

	profiles, err := testDevice.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	res, err = testDevice.GetStreamURI(profiles[0].Token, "UDP")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetSnapshotURI(t *testing.T) {
	var res MediaURI 

	log.Println("Test GetSnapshotURI")

	profiles, err := testDevice.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	res, err = testDevice.GetSnapshotURI(profiles[0].Token)
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}
