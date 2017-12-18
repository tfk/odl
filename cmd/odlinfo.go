package main

import (
	"flag"
	"fmt"
	"odl"
	"os"
)

func printDefaultAndExit() {
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	userPtr := flag.String("username", "", "Username for the ODL-Service (Required)")
	passwordPtr := flag.String("password", "", "Password for the ODL-Service (Required)")
	actionPtr := flag.String("action", "", "Action to be executed (list, detail) (Required)")

	flag.Parse()

	if *userPtr == "" || *passwordPtr == "" || *actionPtr == "" {
		printDefaultAndExit()
	}

	info := odl.NewInfo(*userPtr, *passwordPtr)

	switch *actionPtr {
	case "list":
		stations := info.ListStations()
		for _, station := range stations {
			fmt.Printf("ID: %s\tPlace: %s(%s)\tRadiation: %f\n", station.ID, station.Place, station.Zip, station.Radiation)
		}
	case "detail":
		args := flag.Args()
		for _, id := range args {
			s := info.GetStation(id)
			fmt.Printf("ID: %s\tPlace: %s(%s)\tRadiation: %f\n", s.Info.ID, s.Info.Place, s.Info.Zip, s.Info.Radiation)
		}
	default:
		printDefaultAndExit()
	}
}
