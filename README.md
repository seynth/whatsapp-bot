# Aphrodite

### Requirements

- [FFMPEG](https://www.ffmpeg.org/)
- [Golang](https://go.dev/)

This is whatsapp bot that can convert image/video to sticker and vice versa. First install [ffmpeg](https://www.ffmpeg.org/) to your computer or server to start using this bot.

```
git clone https://github.com/seynth/whatsapp-bot.git
cd whatsapp-bot
go mod tidy
go run .
```

### Scan qr when it show on your terminal.


![Scan qr code to use the bot](/assets/scan-qr.jpg)

Once the prompt say `Aphrodite ready` you can go ahead and type some command to start using it.

## Available command

### `#i2s` - image to sticker (with image attachment)

https://github.com/user-attachments/assets/6d3df877-3d24-42af-bec5-67e1939bef64

### `#v2s` - video to sticker (with video attachment)

https://github.com/user-attachments/assets/105ea3b7-20cc-4029-97bc-31656f8d169c

### `#s2i` - sticker to image (reply to a sticker)

https://github.com/user-attachments/assets/db04b818-96e5-4868-a970-ef4983a1dd95



> [!CAUTION]
> Currently this bot does not support convert animated sticker to video


## What i learn?..

- channel in golang
- goroutine 
- how whatsapp API works
- send whatsapp message with golang
- convert webp to jpg and vice versa
- convert video to webp 
- ~~convert webp to video/gif~~

# Support this repo

Feel free to pull request or open an issue

<a href='https://ko-fi.com/F1F611FQO4' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://storage.ko-fi.com/cdn/kofi1.png?v=3' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>

## Thanks to 
- [whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow)
