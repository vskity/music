package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

var baseURL = "http://music.163.com"

func main() {
	getSongList("/discover/toplist?id=3779629")
	//getTopList()
	//resp, _ := doAction("GET", "http://music.163.com/discover/toplist?id=3779629", nil)
	//echoResp(resp)
}

func getTopList() map[string]string {
	d, err := goquery.NewDocument(baseURL + "/discover/toplist?")
	if err != nil {
		fmt.Println(err.Error())
	}
	topList := make(map[string]string, 21)
	d.Find(".s-fc0").Each(func(i int, s *goquery.Selection) {
		if id, b := s.Attr("href"); b {
			topList[s.Text()] = id
			fmt.Println(s.Text(), "  ", id)
		}
	})
	return topList
}

func getSongList(id string) []*song {
	d, err := goquery.NewDocument(baseURL + id)
	if err != nil {
		fmt.Println(err.Error())
	}

	selection := d.Find("div ul.f-hide li")
	songList := make([]*song, selection.Size())
	selection.Each(func(i int, s *goquery.Selection) {
		ss := s.Find("a")

		tmp := &song{
			id:   ss.AttrOr("href", ""),
			name: ss.Text(),
		}
		songList[i] = tmp
		fmt.Println(tmp.id, "   ", tmp.name)
	})
	return songList
}

func doAction(method, url string, data map[string]interface{}) (*http.Response, error) {

	body := encodeData(data)

	req, nrErr := http.NewRequest(method, url, body)
	if nrErr != nil {
		return nil, nrErr
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,ms;q=0.6,en;q=0.4,zh-TW;q=0.2")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "music.163.com")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.152 Safari/537.36")

	cli := http.Client{}

	return cli.Do(req)

}

func encodeData(data map[string]interface{}) io.Reader {

	return bytes.NewBufferString("")
}

func echoResp(resp *http.Response) {

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Status)
	fmt.Println(resp.Close)

	data, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	io.Copy(os.Stdout, data)

}
