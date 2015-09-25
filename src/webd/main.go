package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"configflag"
)

var (
	configFlags  *flag.FlagSet = flag.NewFlagSet("config flags", flag.ExitOnError)
	configList   configflag.FileList
	listenAddr   string
	redisAddr    string
	basePath     string
	templatePath string
	staticMap    string
)

// Wrapper type for Handler
type StaticHandler struct {
	innerHandler http.Handler // save trapped/wrapped Handler
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("StaticHandler.ServeHTTP url=%s", r.URL.Path)
	handler.innerHandler.ServeHTTP(w, r) // call trapped/wrapped Handler
}

// Initialize package main
func init() {
	configFlags.Var(&configList, "config", "load config flags from this file")
	configFlags.StringVar(&listenAddr, "listenOn", ":8080", "http listen address [addr]:port")
	configFlags.StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis server address")
	configFlags.StringVar(&basePath, "path", "/quixdate", "www base path")
	configFlags.StringVar(&templatePath, "template", "", "template root path")
	configFlags.StringVar(&staticMap, "static", "", "www static mapping")
}

func main() {
	log.Println("webd booting")

	// Parse flags from command-line
	configFlags.Parse(os.Args[1:])

	// Parse flags from files
	log.Printf("config files: %d", len(configList))
	if err := configflag.Load(configFlags, configList); err != nil {
		log.Printf("failure loading config flags: %s", err)
		return
	}

	// Base path
	log.Printf("www base path: %s", basePath)

	// Set template path
	log.Printf("template root path: %s", templatePath)
	if templatePath == "" {
		log.Printf("template root path is required")
		return
	}

	// Set static content paths
	staticMap := strings.TrimSpace(staticMap)
	if staticMap == "" {
		log.Printf("no static map defined")
	} else {
		for _, pair := range strings.Split(staticMap, ",") {
			pDir := strings.Split(pair, ":")
			if len(pDir) != 2 {
				log.Printf("bad static map pair: pair=[%s] map=[%s]", pair, staticMap)
				return
			}
			p := pDir[0]
			if p[0] != '/' {
				p = basePath + "/" + p // prepend base path
			}
			dir := pDir[1]
			log.Printf("installing static handler from path %s to directory %s", p, dir)
			http.Handle(p, StaticHandler{http.StripPrefix(p, http.FileServer(http.Dir(dir)))})
		}
	}

	// Set dynamic content paths
	var homePath string
	if basePath == "" {
		homePath = "/"
	} else {
		homePath = basePath
	}
	log.Printf("installing dynamic handler for path: home=%s", homePath)
	http.HandleFunc(homePath, func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handlerHome) })

	if homePath != "/" {
		log.Printf("installing not-found handler on root path")
		http.HandleFunc("/", notFound) // trap not-found handler
	}

	log.Printf("webd boot complete")

	serve(listenAddr)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("notFound: url=%s", r.URL.Path)
	http.NotFound(w, r) // default not-found handler
}

type session struct {
}

// Add session parameter to handle
func trapHandle(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request, *session)) {
	s := sessionGet(r)
	handler(w, r, s)
}

func handlerHome(w http.ResponseWriter, r *http.Request, s *session) {
	path := r.URL.Path
	log.Printf("dynamic handlerHome url=%s", path)
}

func sessionGet(r *http.Request) *session {
	return &session{}
}

func serve(addr string) {
	if addr == "" {
		log.Printf("server starting on :http (empty address)")
	} else {
		log.Printf("server starting on " + addr)
	}

	err := http.ListenAndServe(addr, nil)
	/*
		s := &http.Server{
			Addr:           addr,
			Handler:        nil,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		err := s.ListenAndServe()
	*/
	if err != nil {
		log.Printf("ListenAndServe: %s: %s", addr, err)
	}
}
