package zbxgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/0x1un/go-zabbix"

	"github.com/sirupsen/logrus"
)

type session struct {
	Host       string
	Username   string
	Password   string
	client     *http.Client
	zbxSession *zabbix.Session
}

type loginParameter struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	AutoLogin string `json:"autologin"`
	Enter     string `json:"enter"`
}

type zbxGraphParmas struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Graphid    int    `json:"graphid"`
	From       string `json:"from"`
	To         string `json:"to"`
	ProfileIdx string `json:"profileIdx"`
}

func NewRequest(method string, url string, params interface{}) (*http.Request, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, bytes.NewReader(data))
}

func NewSeesion(host, username, passwd string) *session {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Println(err)
		return nil
	}
	rpcUrl := "http://" + host + "/" + "api_jsonrpc.php"
	ses, err := zabbix.NewSession(rpcUrl, username, passwd)
	if err != nil {
		logrus.Fatal(err)
	}
	return &session{
		Host:     host,
		Username: username,
		Password: passwd,
		client: &http.Client{
			Jar: jar,
		},
		zbxSession: ses,
	}
}

// 下载流量图，返回图片流
func (c *session) DownloadTrafficGraph(graphid, from, to string) ([]byte, error) {
	graUrl := "http://" + c.Host + "/chart2.php"
	req, err := http.NewRequest(http.MethodGet, graUrl, nil)
	if err != nil {
		return nil, err
	}
	reqQuery := req.URL.Query()
	reqQuery.Set("width", "1920")
	reqQuery.Set("height", "201")
	reqQuery.Set("graphid", graphid)
	reqQuery.Set("from", from)
	reqQuery.Set("to", to)
	reqQuery.Set("profileIdx", "web.graphs.filter")
	req.URL.RawQuery = reqQuery.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	params := zabbix.GraphGetParameters{}
	params.Graphids = []string{graphid}
	params.Output = "extend"
	params.SortField = []string{"name"}
	res, err := c.zbxSession.GraphGet(params)
	if err != nil {
		logrus.Fatal(err)
	}
	name := func() string {
		if len(res) > 0 {
			return res[0].Name
		}
		return ""
	}()
	logrus.Printf("download traffic graph: [%s:%s]\n", name, graphid)
	return data, nil
}

func (c *session) Login() error {
	loginUrl := "http://" + c.Host
	values := url.Values{}
	values.Add("name", c.Username)
	values.Add("password", c.Password)
	values.Add("autologin", "1")
	values.Add("enter", "Sign in")
	req, err := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", c.Host)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login zbx: %d\n", resp.StatusCode)
	}

	return nil
}
