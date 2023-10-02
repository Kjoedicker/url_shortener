package url_shortener

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
)

var (
	ServerAddr = "localhost:8000"
)

type UrlShortener struct {
	UrlMap map[string]string
}
type UrlResponse struct {
	Original     string
	ShortCode    string
	ShortenedUrl string
}

func (u UrlShortener) RootHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(u.UrlMap)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (u UrlShortener) UrlRedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortenedURL, ok := u.getUrl(shortCode); ok {
		fmt.Println("Attempting to redirect", shortCode, u.UrlMap[shortCode])
		http.Redirect(w, r, shortenedURL, http.StatusFound)
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "URL not found for: %s", shortCode)
}

func hashUrl(url string) (hashedUrl string) {
	hash := fnv.New32()
	hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum32())
}

func (u UrlShortener) getUrl(hash string) (string, bool) {
	url, ok := u.UrlMap[hash]
	return url, ok
}

func (u UrlShortener) UrlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]

	hashedUrl := hashUrl(url)
	u.UrlMap[hashedUrl] = url

	responseObject := UrlResponse{
		Original:     url,
		ShortCode:    hashedUrl,
		ShortenedUrl: path.Join(ServerAddr, hashedUrl),
	}
	jsonResponse, err := json.Marshal(&responseObject)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func buildRouting() *mux.Router {
	urlShortener := UrlShortener{UrlMap: make(map[string]string)}

	router := mux.NewRouter()

	router.HandleFunc("/", urlShortener.RootHandler)
	router.HandleFunc("/{shortCode}", urlShortener.UrlRedirectHandler)
	router.HandleFunc("/shorten/{url}", urlShortener.UrlShortenerHandler)

	return router
}

func RunServer() {
	router := buildRouting()

	server := http.Server{
		Handler:      router,
		Addr:         ServerAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on:", server.Addr)
	log.Fatal(server.ListenAndServe())
}
