# Aphrodite

### Requirements

- [FFMPEG](https://www.ffmpeg.org/)

This is whatsapp bot that can convert image/video to sticker and vice versa. First install ffmpeg to your computer or server to start using this bot.

```
git clone https://github.com/seynth/whatsapp-bot.git
cd whatsapp-bot
go mod tidy
go run .
```

### Scan qr when it show on your terminal.

![Scan qr code to use the bot](/assets/scan-qr.jpg)

Once the prompt say `Aphrodite ready` you can go ahead and type some command to start using it.

### Available command

- `#i2s` - image to sticker (with image attachment)
- `#v2s` - video to sticker (with video attachment)
- `#s2i` - sticker to image (reply to a sticker)

> [!CAUTION]
> Currently this bot does not support converting animated sticker to video



