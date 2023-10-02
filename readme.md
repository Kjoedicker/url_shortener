<br />
<div align="center">
    <h3 align="center">URL Shortener</h3>
</div>

## About The Project
A volatile URL shortener
## API

* GET /
    * `200` = Success 
* GET /shorten/{url}
    * `201` = Created
* GET /{shortCode}
    * `302` = Redirect 
    * `404` = URL not found

## Usage

```
# Returns `201` and an JSON response with the new key. "{"Original":"github.com","ShortCode":"c49c601d","ShortenedUrl":"localhost:8000/c49c601d"}"
curl -v http://localhost:8080/shorten/github.com

#  Returns `200` and dump of the current keys. "{"c49c601d": "http://github.com"}"
curl -v http://localhost:8080/

# Return `302` indicating a redirect
curl -v http://localhost:8080/c49c601d9
```
