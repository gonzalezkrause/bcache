package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var (
	Version   = "0.0.1"
	BuildDate = "nil"
	BuildId   = "NoGitID"
)

var (
	BucketDoesNotExist   = errors.New("Cache bucket does not exist")
	KeyDoesNotExistError = errors.New("Key does not exist")
)

type server struct {
	db *bolt.DB
}

func newServer(filename string) (s *server, err error) {
	s = new(server)
	s.db, err = bolt.Open(filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

func (s *server) Put(bucket, key string, val []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return b.Put([]byte(key), val)
	})
}

func (s *server) Del(bucket, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketDoesNotExist
		}

		return b.Delete([]byte(key))
	})
}

func (s *server) Get(bucket, key string) (data []byte, err error) {
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketDoesNotExist
		}

		r := b.Get([]byte(key))

		if r != nil {
			data = make([]byte, len(r))
			copy(data, r)
		} else {
			return KeyDoesNotExistError
		}

		return nil
	})
	return
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["bucket"] == "" || vars["key"] == "" {
		http.Error(w, "Missing bucket or key", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "POST", "PUT":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.Put(vars["bucket"], vars["key"], data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":
		if err := s.Del(vars["bucket"], vars["key"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "GET":
		data, err := s.Get(vars["bucket"], vars["key"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

func main() {
	var (
		addr   string
		dbfile string
		remove bool
	)

	flag.StringVar(&addr, "listen", ":3017", "Address to listen on")
	flag.StringVar(&dbfile, "db", "/var/boltcache/cache.db", "Cache file")
	flag.BoolVar(&remove, "rm", false, "Remove cache file after app closes")
	flag.Parse()

	log.Printf("Starting BoltCache v%s (%s - %s)\n", Version, BuildDate, BuildId)
	log.Println("Using Bolt DB file:", dbfile)
	if remove {
		log.Println("Cache file persistence disabled!")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down BoltCache server")
		if remove {
			os.Remove(dbfile)
		}
		os.Exit(0)
	}()

	server, err := newServer(dbfile)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	router := mux.NewRouter()
	router.Handle("/v1/{bucket}/{key}", server)
	http.Handle("/", router)

	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
