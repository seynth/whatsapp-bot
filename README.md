# Aphrodite


### Requirements

- [FFMPEG](https://www.ffmpeg.org/)

This is whatsapp bot that can convert image/video to sticker and vice versa. First install ffmpeg to your computer to start using this bot.

```
git clone https://github.com/seynth/whatsapp-bot.git
cd whatsapp-bot
go mod tidy
go run .
```


Then you need to scan qr-code to be able to use this bot, once the prompt say `Bot ready` you can go ahead and type some command to start using it.

### Available command

- `#i2s` - image to sticker (with image attachment)
- `#v2s` - video to sticker (with video attachment)
- `#s2i` - sticker to image (reply to a sticker)

> [!CAUTION]
> Currently this bot does not support converting animated sticker to video


