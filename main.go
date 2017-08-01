package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

var (
	baseURL = "http://music.163.com"
	jar     http.CookieJar
	nonce   = "0CoJUm6Qyw8W8jud"
)

type song struct {
	id     string
	name   string
	time   string
	artist string
}

func main() {

	//fmt.Println(aesEnv("exampleplaintext")) //e1cdb90013f76bdf10c3d76b40e5e1643acf2e1c6787459fad7311815e9f0af8
	a := []byte("exampleplaintext") //d20e0180d158f93a4d748e7e8455789f4cd913385d96c032fdcb4fad5b3c5d4c
	fmt.Println(encodeData(a))      //ec85f86eebc62c4d63341ba2f1bb380b4cd913385d96c032fdcb4fad5b3c5d4c

	// fmt.Println(encryptSong("491228747"))
	// getSongList("/discover/toplist?id=3779629")
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

	body := encodeData([]byte(""))

	req, nrErr := http.NewRequest(method, url, bytes.NewBufferString(body))
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

	cli := http.Client{
		Jar: jar,
	}
	return cli.Do(req)
}

func encodeData(text, key []byte) string {

	// PKCS5Padding
	padding := 16 - len(text)%16
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	text = append(text, padtext...)

	block, err := aes.NewCipher([]byte(""))
	if err != nil {
		panic(err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("read full error ", err)
		return ""
	}

	dst := make([]byte, len(text))

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, text)

	return fmt.Sprintf("%x", dst)
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

func encryptSong(id string) string {
	m := []byte("3go8&$8*3")
	song := []byte(id)
	for k, v := range song {
		song[k] = v ^ m[k]
	}
	hash := md5.Sum(song)
	return base64.URLEncoding.EncodeToString(hash[:])
}

func reverse(r []byte) []byte {
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}
