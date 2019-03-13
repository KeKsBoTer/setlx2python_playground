package setlx2python_playground

import (
	"io/ioutil"
	"net/http"
)

func fetchURL(url string) (*string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := string(body)
	return &data, nil
}
