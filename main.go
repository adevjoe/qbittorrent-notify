package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	// BotToken define bot token
	BotToken string
	// ChatID define chat id
	ChatID int64
	// qBittorrent host
	qBittorrentHost string
	// qBittorrent username
	qBittorrentUsername string
	// qBittorrent password
	qBittorrentPassword string
	// tick time
	tickTime time.Duration = 10 * time.Second
)

func init() {
	chatDefault, _ := strconv.Atoi(os.Getenv("CHAT_ID"))
	// init args
	flag.StringVar(&BotToken, "botToken", os.Getenv("BOT_TOKEN"), "BotToken define bot token")
	flag.Int64Var(&ChatID, "chatID", int64(chatDefault), "ChatID define chat id")
	flag.StringVar(&qBittorrentHost, "qbHost", os.Getenv("QBITTORRENT_HOST"), "qBittorrent host")
	flag.StringVar(&qBittorrentUsername, "qbUsername", os.Getenv("QBITTORRENT_USERNAME"), "qBittorrent username")
	flag.StringVar(&qBittorrentPassword, "qbPassword", os.Getenv("QBITTORRENT_PASSWORD"), "qBittorrent password")
	flag.Parse()

	// init bot
	initBot()

	// initQb
	initQB()
}

func main() {
	for {
		torrents, err := getTorrents()
		if err != nil {
			log.Println(err)
			time.Sleep(tickTime)
			continue
		}
		for _, torrent := range *torrents {
			if torrent.notified() {
				continue
			}
			err := setTag(torrent.Hash, TagNotified)
			if err != nil {
				log.Printf("set Tag %s err, %v", torrent.Name, err)
				continue
			}
			err = msg("qBittorrent Downloaded",
				fmt.Sprintf("```\nName: %s\nSize: %s\nCategory: %s\nAdded On: %s\n```",
					torrent.Name, torrent.binarySize(), torrent.Category, torrent.addedTime()))
			if err != nil {
				_ = unsetTag(torrent.Hash, TagNotified)
				log.Printf("notify %s err, %v", torrent.Name, err)
			}
		}
		time.Sleep(tickTime)
	}
}
