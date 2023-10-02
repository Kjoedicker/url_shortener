package url_shortener

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var expect = assert.Equal

func Test_RootHandler(t *testing.T) {
	urlShortener := UrlShortener{UrlMap: make(map[string]string)}
	t.Run("GET /", func(t *testing.T) {

		t.Run("When there are not shortened URLS", func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8000", nil)
			w := httptest.NewRecorder()

			urlShortener.RootHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, string(body), "{}", "Should return an empty JSON body")
		})

		t.Run("When there are shortened URLS", func(t *testing.T) {
			urlShortener.UrlMap["example"] = "example.com"

			req := httptest.NewRequest("GET", "http://localhost:8000", nil)
			w := httptest.NewRecorder()

			urlShortener.RootHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			var urlMap map[string]string
			json.Unmarshal(body, &urlMap)

			assert.Equal(t, urlMap["example"], "example.com", "Should return a map with the expected values")
		})
	})
}

func Test_URLRedirectHandler(t *testing.T) {
	urlShortener := UrlShortener{UrlMap: make(map[string]string)}
	mockShortUrl := "8675309"
	mockUrl := "http://example.com"

	t.Run("GET /{url}", func(t *testing.T) {

		t.Run("When the short URL does not exist", func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8000/"+mockShortUrl, nil)
			req = mux.SetURLVars(req, map[string]string{"url": mockShortUrl})
			w := httptest.NewRecorder()

			urlShortener.UrlRedirectHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, resp.StatusCode, 404, "Should return 404")
			assert.Equal(t, string(body), "URL not found for: "+mockShortUrl, "Should return a message with context")
		})

		t.Run("When the short URL does exist", func(t *testing.T) {

			urlShortener.UrlMap[mockShortUrl] = mockUrl

			req := httptest.NewRequest("GET", "http://localhost:8000/"+mockShortUrl, nil)
			req = mux.SetURLVars(req, map[string]string{"url": mockShortUrl})

			w := httptest.NewRecorder()

			urlShortener.UrlRedirectHandler(w, req)

			resp := w.Result()

			location, err := resp.Location()
			if err != nil {
				log.Fatal(err)
			}

			assert.Equal(t, resp.StatusCode, 302, "Should return 302, indicating redirect")
			assert.Equal(t, mockUrl, location.String(), "Should have a location of: "+mockUrl)
		})
	})
}

func Test_UrlShortenerHandler(t *testing.T) {
	urlShortener := UrlShortener{UrlMap: make(map[string]string)}
	mockShortUrl := "8675309"
	mockUrl := "http://example.com"

	t.Run("GET /shorten/{url}", func(t *testing.T) {

		t.Run("When the short URL does not exist", func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8000/shorten/"+mockShortUrl, nil)
			req = mux.SetURLVars(req, map[string]string{"url": mockUrl})
			w := httptest.NewRecorder()

			urlShortener.UrlShortenerHandler(w, req)

			resp := w.Result()

			expect(t, resp.StatusCode, 201)

			var urlResponse UrlResponse
			body, _ := io.ReadAll(resp.Body)
			json.Unmarshal(body, &urlResponse)

			expect(t, mockUrl, urlResponse.Original)
			expect(t, mockUrl, urlShortener.UrlMap[urlResponse.Shortened])
		})
	})
}
