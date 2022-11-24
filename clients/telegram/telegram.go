package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Nau077/golang-tg-bot/lib/e"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// клиент будет заниматься 2 вещами
// принимать сообщения и отправлять пользователю

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	// сформируем параметры запроса
	q := url.Values{}
	// Int to ASCII
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	// отправить запрос
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	// результат в jsonб его нужно распарсить
	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil

}

func (c Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("Не могу отправить запрос:", err)
	}

	return nil
}

func (c Client) doRequest(method string, query url.Values) (data []byte, err error) {
	const errMsg = "Не смогли выполнить request"
	defer func() { err = e.WrapIfErr(errMsg, err) }()
	// сформируем url, на который будет отправляться запрос
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// связано с error.As и erros.Is
		return nil, err
	}
	// передать объект req в параметры запроса из аргумента
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		// связано с error.As и erros.Is
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
