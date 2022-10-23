# qBittorrent notify

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/adevjoe/qbittorrent-notify/latest?style=for-the-badge)
[![Docker Pulls](https://img.shields.io/docker/pulls/adevjoe/qbittorrent-notify?label=qbittorrent-notify%20pulls&style=for-the-badge)](https://hub.docker.com/repository/docker/adevjoe/qbittorrent-notify)
[![Docker Image Version (latest by semver)](https://img.shields.io/docker/v/adevjoe/qbittorrent-notify?sort=semver&style=for-the-badge)](https://hub.docker.com/repository/docker/adevjoe/qbittorrent-notify)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/adevjoe/qbittorrent-notify?sort=semver?style=for-the-badge)](https://hub.docker.com/repository/docker/adevjoe/qbittorrent-notify)

qBittorrent notify when torrents complete

Only run on qBittorrent v4.1+

### Usage

#### Binary

```sh
go build
./qbittorrent-notify -botToken=your_token -chatID=your_chat_id -qbHost=your_qb_host -qbUsername=your_qb_username -qbPassword=your_qb_password
```

#### Docker

```sh
docker run -d --restart=always \
    -e BOT_TOKEN=your_token \
    -e CHAT_ID=your_chat_id \
    -e QBITTORRENT_HOST=your_qb_host \
    -e QBITTORRENT_USERNAME=your_qb_username \
    -e QBITTORRENT_PASSWORD=your_qb_password \
    adevjoe/qbittorrent-notify:latest
```

#### Docker compose

Download docker-compose.yaml file, and run with docker-compose in same dir.

```
docker-compose up -d
```

### Telegram Bot

How can i get telegram bot chat id?

Follow [this link](https://stackoverflow.com/questions/32423837/telegram-bot-how-to-get-a-group-chat-id)
