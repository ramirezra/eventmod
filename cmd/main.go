package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/inklabs/rangedb/provider/leveldbstore"
)

const (
	httpTimeout = 10 * time.Second
)

func main() {
	fmt.Println("RangeDB API")

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	port := flag.Int("port", 3000, "port")
	// baseURI := flag.String("baseUri", "http://0.0.0.0:3000", "")
	dbPath := flag.String("dbPath", ".leveldb", "path to LevelDB directory")
	// templatesPath := flag.String("templates", "", "optionsal templates path")
	// gRPCPort := flag.Int("gRPCPort", 3001, "gRPC port")
	flag.Parse()

	httpAddress := fmt.Sprintf("0.0.0.0:%d", *port)

	logger := log.New(os.Stderr, "", 0)
	leveldbStore, err := leveldbstore.New(*dbPath, leveldbstore.WithLogger(logger))
	if err != nil {
		log.Fatalf("Unable to load db (%s): %v", *dbPath, err)
	}

	muxServer := http.NewServeMux()
	muxServer.Handle("/", ui)

	httpServer := &http.Server{
		Addr:         httpAddress,
		ReadTimeout:  httpTimeout + time.Second,
		WriteTimeout: httpTimeout + time.Second,
		Handler:      muxServer,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go serveHTTP(httpServer, httpAddress)
	<-stop

}

func serveHTTP(srv *http.Server, addr string) {
	fmt.Printf("Listening: http://%s/\n", addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
