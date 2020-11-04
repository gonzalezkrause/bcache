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

	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

var (
	// Version to be set on compilation
	Version = "0.0.1"
	// BuildDate to be set on compilation
	BuildDate = "nil"
	// BuildID to be set on compilation
	BuildID = "NoGitID"
)

var (
	// ErrBucketDoesNotExist error
	ErrBucketDoesNotExist = errors.New("Cache bucket does not exist")
	// ErrKeyDoesNotExistError error
	ErrKeyDoesNotExistError = errors.New("Key does not exist")
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
			return ErrBucketDoesNotExist
		}

		return b.Delete([]byte(key))
	})
}

func (s *server) Get(bucket, key string) (data []byte, err error) {
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketDoesNotExist
		}

		r := b.Get([]byte(key))

		if r != nil {
			data = make([]byte, len(r))
			copy(data, r)
		} else {
			return ErrKeyDoesNotExistError
		}

		return nil
	})
	return
}

func (s *server) Buckets() (data []string, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			data = append(data, string(name))
			return nil
		})
	})

	return data, err
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if r.URL.Path == "/" {
		w.Write([]byte("BCache v" + Version + "\r\n"))
		return
	}

	if r.URL.Path == "/v1/buckets" {
		bucks, err := s.Buckets()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var str string
		for i := range bucks {
			str += bucks[i] + "\n"
		}

		w.Write([]byte(str))
		return
	}

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
	flag.StringVar(&dbfile, "db", "./cache.db", "Cache file")
	flag.BoolVar(&remove, "rm", false, "Remove cache file after app closes")
	flag.Parse()

	log.Printf("Starting BCache v%s (%s - %s)\n", Version, BuildDate, BuildID)
	log.Println("Using BoltDB file:", dbfile)
	if remove {
		log.Println("Cache file persistence disabled!")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down BCache server")
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
	router.Handle("/v1/buckets", server)
	router.Handle("/", server)
	http.Handle("/", router)

	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
