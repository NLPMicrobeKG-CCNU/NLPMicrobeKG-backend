package service

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var GraphDBURL = viper.GetString("graphdb.addr")
var Repository = viper.GetString("graphdb.repository")
var BaseQueryURL = "http://" + GraphDBURL + "/repositories/" + Repository

func QueryInfo(query string, limit, offset int) ([]byte, error) {
	method := "POST"
	data := url.Values{}
	data.Set("query", query)
	data.Set("infer", "true")
	data.Set("sameAs", "true")
	data.Set("limit", strconv.Itoa(limit))
	data.Set("offset", strconv.Itoa(offset))
	client := &http.Client{}

	req, err := http.NewRequest(method, BaseQueryURL, strings.NewReader(data.Encode()))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Host", "localhost:7200")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "110")
	req.Header.Add("Accept", "application/x-sparqlstar-results+json, application/sparql-results+json;q=0.9, */*;q=0.8")
	req.Header.Add("X-GraphDB-Track-Alias", "query-editor-889628.629999992-1606458002690")
	req.Header.Add("X-GraphDB-Catch", "1000; throw")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.193 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Origin", "http://localhost:7200")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "http://localhost:7200/sparql")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7,en-US;q=0.6")
	req.Header.Add("Cookie", "Hm_lvt_b156166357d4f3a5cc8376904f21b47d=1604472565; _ga=GA1.1.1213387063.1604472565; _hjid=4c49e293-3b30-4911-a6a5-c3ad968a8e7d; Hm_lpvt_b156166357d4f3a5cc8376904f21b47d=1604477003; Hm_lvt_a917013a276e2f543da604a6f2bba5a5=1605373902; Hm_lpvt_a917013a276e2f543da604a6f2bba5a5=1605374426; io=q7B3NJ9fleipoWBLAAAE; company_id=9; elite_user_type=student; com.ontotext.graphdb.repository7200=testk; elite_token=9bdebff1a37594463936342b49110b416e1ad0cf2ccb107282265bcaf7429e0eff81719ab7b6ab5c2a21bd720a04ce3155c3fb1f0c19b9163c1ebd45501e8e")

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
