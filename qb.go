package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var client *http.Client

func initQB() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client = http.DefaultClient
	client.Jar = jar
	client.Timeout = 30 * time.Second
	client.Transport = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	login()
}

func login() {
	loginURL := fmt.Sprintf("%s/api/v2/auth/login?username=%s&password=%s", qBittorrentHost, qBittorrentUsername, qBittorrentPassword)
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		log.Println(err)
		return
	}
	rsp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	if rsp.StatusCode != http.StatusOK {
		log.Println("login request fail")
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if string(body) == "Fails." {
		log.Println("login failed")
		return
	}
	log.Println("login success")
}

type Torrent struct {
	AddedOn       int64   `json:"added_on"`
	Dlspeed       int64   `json:"dlspeed"`
	Eta           int64   `json:"eta"`
	FLPiecePrio   bool    `json:"f_l_piece_prio"`
	ForceStart    bool    `json:"force_start"`
	Hash          string  `json:"hash"`
	Category      string  `json:"category"`
	Tags          string  `json:"tags"`
	Name          string  `json:"name"`
	NumComplete   int     `json:"num_complete"`
	NumIncomplete int     `json:"num_incomplete"`
	NumLeechs     int     `json:"num_leechs"`
	NumSeeds      int     `json:"num_seeds"`
	Priority      int     `json:"priority"`
	Progress      float64 `json:"progress"`
	Ratio         float64 `json:"ratio"`
	SeqDl         bool    `json:"seq_dl"`
	Size          int64   `json:"size"`
	State         string  `json:"state"`
	SuperSeeding  bool    `json:"super_seeding"`
	Upspeed       int64   `json:"upspeed"`
}

const (
	TagNotified = "Notified"
)

func (t *Torrent) notified() bool {
	s := strings.Split(t.Tags, ",")
	for _, v := range s {
		if strings.TrimSpace(v) == TagNotified {
			return true
		}
	}
	return false
}

func (t *Torrent) binarySize() string {
	const unit = 1024
	if t.Size < unit {
		return fmt.Sprintf("%d B", t.Size)
	}
	div, exp := int64(unit), 0
	for n := t.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(t.Size)/float64(div), "KMGTPE"[exp])
}

func (t *Torrent) addedTime() string {
	return time.Unix(t.AddedOn, 0).Format("2006-01-02 15:04:05")
}

func getTorrents() (*[]Torrent, error) {
	params := map[string]string{
		"filter": "completed",
	}
	result := new([]Torrent)
	err := do("torrents", "info", params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func setTag(hash, tags string) error {
	params := map[string]string{
		"hashes": hash,
		"tags":   tags,
	}
	err := do("torrents", "addTags", params, nil)
	if err != nil {
		return err
	}
	return nil
}

func unsetTag(hash, tags string) error {
	params := map[string]string{
		"hashes": hash,
		"tags":   tags,
	}
	err := do("torrents", "removeTags", params, nil)
	if err != nil {
		return err
	}
	return nil
}

func do(apiName, methodName string, params map[string]string, result interface{}) error {
	u, err := url.Parse(fmt.Sprintf("%s/api/v2/%s/%s?", qBittorrentHost, apiName, methodName))
	if err != nil {
		return err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	var rsp *http.Response
	rsp, err = client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode == http.StatusForbidden {
		log.Println("session expired")
		login()
		// retry
		rsp, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s/%s request fail, status code: %d", apiName, methodName, rsp.StatusCode)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}
