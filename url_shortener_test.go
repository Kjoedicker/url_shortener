package url_shortener

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	expect = assert.Equal

	baseUrl = "localhost:8000"

	mockUrl       = "http://example.com"
	mockShortCode = "8675309"
)

func buildUrlShortener() (urlShortener UrlShortener) {
	return UrlShortener{UrlMap: make(map[string]string)}
}

func Test_RootHandler(t *testing.T) {
	urlShortener := buildUrlShortener()

	t.Run("GET /", func(t *testing.T) {
		requestUrl := path.Join(baseUrl, "/")

		t.Run("When there are not shortened URLS", func(t *testing.T) {
			req := httptest.NewRequest("GET", requestUrl, nil)
			w := httptest.NewRecorder()

			urlShortener.RootHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, string(body), "{}", "Should return an empty JSON body")
		})

		t.Run("When there are shortened URLS", func(t *testing.T) {
			urlShortener.UrlMap["example"] = "example.com"

			req := httptest.NewRequest("GET", requestUrl, nil)
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
	urlShortener := buildUrlShortener()
	requestUrl := path.Join(baseUrl, mockShortCode)

	t.Run("GET /{shortCode}", func(t *testing.T) {

		t.Run("When the short URL does not exist", func(t *testing.T) {
			req := httptest.NewRequest("GET", requestUrl, nil)
			req = mux.SetURLVars(req, map[string]string{"shortCode": mockShortCode})
			w := httptest.NewRecorder()

			urlShortener.UrlRedirectHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, resp.StatusCode, 404, "Should return 404")
			assert.Equal(t, string(body), "URL not found for: "+mockShortCode, "Should return a message with context")
		})

		t.Run("When the short URL does exist", func(t *testing.T) {
			urlShortener.UrlMap[mockShortCode] = mockUrl

			req := httptest.NewRequest("GET", requestUrl, nil)
			req = mux.SetURLVars(req, map[string]string{"shortCode": mockShortCode})

			w := httptest.NewRecorder()

			urlShortener.UrlRedirectHandler(w, req)

			resp := w.Result()

			location, err := resp.Location()
			if err != nil {
				log.Fatal(err)
			}

			expect(t, resp.StatusCode, 302, "Should return 302, indicating redirect")
			expect(t, mockUrl, location.String(), "Should have a location of: "+mockUrl)
		})
	})
}

func Test_UrlShortenerHandler(t *testing.T) {
	urlShortener := UrlShortener{UrlMap: make(map[string]string)}

	t.Run("GET /shorten/{url}", func(t *testing.T) {
		requestUrl := path.Join(baseUrl, "shorten", mockUrl)

		req := httptest.NewRequest("GET", requestUrl, nil)
		req = mux.SetURLVars(req, map[string]string{"url": mockUrl})
		w := httptest.NewRecorder()

		urlShortener.UrlShortenerHandler(w, req)

		resp := w.Result()

		expect(t, resp.StatusCode, 201)

		var urlResponse UrlResponse
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &urlResponse)

		expect(t, mockUrl, urlResponse.Original, "The response should have returned the original URL")
		expect(t, path.Join(baseUrl, urlResponse.ShortCode), urlResponse.ShortenedUrl, "The response should have returned a shortened URL")
		expect(t, mockUrl, urlShortener.UrlMap[urlResponse.ShortCode], "The shortened URL should map to the original URL")
	})
}
