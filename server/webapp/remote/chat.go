package remote

import (
	"net/url"
	"net/http"
	"fmt"
	"strings"
	"liulishuo/somechat/server/webapp/conf"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"io/ioutil"
)

func ChatApiHtppPost(path string, form url.Values) ([]byte, error) {
	client := new(http.Client)
	remoteUrl := fmt.Sprintf(`http://%s%s`, conf.Conf.Remote.Chat, path)
	req, reqErr := http.NewRequest("POST", remoteUrl, strings.NewReader(form.Encode()))
	if reqErr != nil {
		err := fmt.Errorf("new http request failed, chat remote api, url = %s, error = %v", remoteUrl, reqErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, respErr := client.Do(req)
	if respErr != nil {
		err := fmt.Errorf("post http request failed, chat remote api, url = %s, error = %v", remoteUrl, respErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
		return nil, err
	}
	defer resp.Body.Close()
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		err := fmt.Errorf("read post http response failed, chat remote api, url = %s, error = %v", remoteUrl, bodyErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
		return nil, err
	}
	return body, nil
}
