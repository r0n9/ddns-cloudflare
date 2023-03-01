package notify

import (
	"errors"
	"fmt"
	"net/http"
)

func Send(sendKey, content string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://sctapi.ftqq.com/%s.send", sendKey), nil)
	if err != nil {
		return errors.New("failed to notify: " + content)
	}

	q := req.URL.Query()
	q.Add("title", "Public IP Address: "+content)
	q.Add("desp", "Public IP Address: "+content)
	req.URL.RawQuery = q.Encode()

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
