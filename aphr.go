package main

import (
	"aphrodite/auxiliary"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
)

type Aphrodite struct {
	aphrClient *whatsmeow.Client
	evtHandler uint32
}

func main() {

	dbContainer, err := sqlstore.New("sqlite3", "file:aphrodite.db?_foreign_keys=on", nil)
	if err != nil {
		log.Fatalln("error creating database: ", err)
	}

	storeDevice, errDevice := dbContainer.GetFirstDevice()
	if errDevice != nil {
		panic(errDevice)
	}

	// debuging := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(storeDevice, nil)
	aphrClient := Aphrodite{
		aphrClient: client,
	}
	defer client.Disconnect()

	fmt.Println("[+] Setting up Aphrodite please wait . . .")
	aphrClient.aphrSetup()

	if client.Store.ID == nil {
		qrc, _ := client.GetQRChannel(context.Background())
		fmt.Println("[Aphrodite] Connecting to server . . .")
		err := client.Connect()
		if err != nil {
			panic(err)
		}

		for e := range qrc {
			if e.Event == "code" {
				fmt.Println("Scan QR dibawah ini untuk menggunakan Aphrodite . . .")
				qrterminal.Generate(e.Code, qrterminal.H, os.Stdout)
			} else {

				fmt.Println("login event: ", e.Event)
			}
		}
		fmt.Println("Aphrodite ready . . .")
	} else {
		fmt.Println("[Aphrodite] Connecting to server . . .")
		errClient := client.Connect()
		if errClient != nil {
			log.Fatalln("Error client: ", errClient)
		}
		fmt.Println("Aphrodite ready . . .")
	}

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
	<-ctrlC

}

func (aphr *Aphrodite) aphrSetup() {
	aphr.evtHandler = aphr.aphrClient.AddEventHandler(aphr.aphrHandler)
}

func (aphr *Aphrodite) aphrHandler(ev interface{}) {
	switch event := ev.(type) {
	case *events.Message:

		if event.Info.Type == "text" {
			// fmt.Println("text type", event.Message)
			go auxiliary.AphrCommandHandler(aphr.aphrClient, event)
		} else if event.Info.Type == "media" {
			// fmt.Println("media type", event.Message)
			go auxiliary.AphrMediaHandler(event.Info.MediaType, event, aphr.aphrClient)
		}

	}
}
