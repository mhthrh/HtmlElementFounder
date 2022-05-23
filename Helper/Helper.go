package Helper

import (
	"io/ioutil"
	"net/http"
)

func ReadPage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
func LinkValid(string2 string) bool {
	return true
}
