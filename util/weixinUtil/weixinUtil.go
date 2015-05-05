package weixinUtil

import (
	"io/ioutil"
	"kaixin/const/weixin"
	"net/http"
	"strings"
)

func FetchAccessToken() (string, error) {
	requestUrl := strings.Join([]string{weixin.AccessTokenFetchUrl,
		"?grant_type=client_credential&appid=",
		weixin.AppId, "&secret=", weixin.Appsecret}, "")

	response, err := http.Get(requestUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
