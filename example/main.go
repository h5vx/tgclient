package main

import (
	"flag"
	"log"
	"os"

	"github.com/ansel1/merry"

	"../../mtproto"
)

func main() {
	appID := flag.Int("app_id", 0, "app id")
	appHash := flag.String("app_hash", "", "app hash")
	flag.Parse()

	if *appID == 0 || *appHash == "" {
		println("App ID and hash are required!")
		flag.Usage()
		os.Exit(2)
	}

	if err := start(int32(*appID), *appHash); err != nil {
		log.Fatal(merry.Details(err))
	}
}

func start(appID int32, appHash string) error {
	m, err := mtproto.NewMTProto(appID, appHash)
	if err != nil {
		return merry.Wrap(err)
	}
	if err := m.Connect(); err != nil {
		return merry.Wrap(err)
	}

	for {
		res := m.SendSync(mtproto.TL_updates_getState{})
		if mtproto.IsErrorType(res, mtproto.TL_ErrUnauthorized) { //AUTH_KEY_UNREGISTERED SESSION_REVOKED SESSION_EXPIRED
			if err := m.Auth(); err != nil {
				return merry.Wrap(err)
			}
			continue
		}
		_, ok := res.(mtproto.TL_updates_state)
		if !ok {
			return mtproto.WrongRespError(res)
		}
		break
	}
	log.Println("Seems authed.")

	if err := m.GetContacts(); err != nil {
		return merry.Wrap(err)
	}

	<-chan bool(nil) //pausing forever
	return nil
}