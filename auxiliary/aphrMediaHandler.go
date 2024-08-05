package auxiliary

import (
	"context"
	"log"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func AphrMediaHandler(mediaType string, evt *events.Message, aphr *whatsmeow.Client) {
	switch mediaType {
	case "image":
		imgMsg := evt.Message.GetImageMessage()
		cmd := strings.TrimSpace(imgMsg.GetCaption())

		if cmd == "#i2s" || cmd == "#image2sticker" {
			imgByte, err := aphr.DownloadAny(evt.Message)
			if err != nil {
				log.Fatalln("error download image: ", err)
			}

			webpChan := make(chan []byte)
			defer close(webpChan)

			go AphrConvertJpeg2WebP(imgByte, webpChan)

			res, errUp := aphr.Upload(context.Background(), <-webpChan, whatsmeow.GetMediaType(evt.Message.ImageMessage))
			if errUp != nil {
				log.Fatalln("Error upload image: ", errUp)
			}

			stickerMsgStruct := &waE2E.StickerMessage{
				URL:           &res.URL,
				FileSHA256:    res.FileSHA256,
				FileEncSHA256: res.FileEncSHA256,
				FileLength:    &res.FileLength,
				MediaKey:      res.MediaKey,
				Mimetype:      proto.String("image/webp"),
				DirectPath:    &res.DirectPath,
				Height:        proto.Uint32(512),
				Width:         proto.Uint32(512),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: &waE2E.Message{
						ImageMessage: &waE2E.ImageMessage{
							URL:           imgMsg.URL,
							FileSHA256:    imgMsg.FileSHA256,
							FileLength:    imgMsg.FileLength,
							FileEncSHA256: imgMsg.FileEncSHA256,
							DirectPath:    imgMsg.DirectPath,
							MediaKey:      imgMsg.MediaKey,
						},
					},
				},
			}

			aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
				StickerMessage: stickerMsgStruct,
			})

		}
	case "video":
		vidMsg := evt.Message.GetVideoMessage()
		cmd := strings.TrimSpace(vidMsg.GetCaption())

		if cmd == "#v2s" || cmd == "#video2sticker" {
			stickerBytes := make(chan []byte)
			defer close(stickerBytes)

			mp4Bytes, errBytes := aphr.DownloadAny(evt.Message)
			if errBytes != nil {
				log.Fatalln("Error download mp4: ", errBytes)
			}

			replyCommand := &waE2E.ExtendedTextMessage{
				Text: proto.String("*[Aphrodite]* Converting to *sticker* wait . . ."),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: &waE2E.Message{
						VideoMessage: &waE2E.VideoMessage{
							URL:           vidMsg.URL,
							FileSHA256:    vidMsg.FileSHA256,
							FileEncSHA256: vidMsg.FileEncSHA256,
							FileLength:    vidMsg.FileLength,
							MediaKey:      vidMsg.MediaKey,
							DirectPath:    vidMsg.DirectPath,
						},
					},
				},
			}

			aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
				ExtendedTextMessage: replyCommand,
			})
			go AphrConvertMp4toWebp(mp4Bytes, stickerBytes)

			resSticker, errUpSticker := aphr.Upload(context.Background(), <-stickerBytes, whatsmeow.GetMediaType(evt.Message.StickerMessage))
			if errUpSticker != nil {
				log.Fatalln("Error upload sticker bytes: ", errUpSticker)
			}

			isAnim := true
			stickerMsgStruct := &waE2E.StickerMessage{
				URL:           &resSticker.URL,
				FileSHA256:    resSticker.FileSHA256,
				FileEncSHA256: resSticker.FileEncSHA256,
				FileLength:    &resSticker.FileLength,
				MediaKey:      resSticker.MediaKey,
				Mimetype:      proto.String("image/webp"),
				DirectPath:    &resSticker.DirectPath,
				Height:        proto.Uint32(512),
				Width:         proto.Uint32(512),
				IsAnimated:    &isAnim,
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: &waE2E.Message{
						VideoMessage: &waE2E.VideoMessage{
							URL:           vidMsg.URL,
							Caption:       vidMsg.Caption,
							FileSHA256:    vidMsg.FileSHA256,
							FileEncSHA256: vidMsg.FileEncSHA256,
							FileLength:    vidMsg.FileLength,
							MediaKey:      vidMsg.MediaKey,
							DirectPath:    vidMsg.DirectPath,
						},
					},
				},
			}
			aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
				StickerMessage: stickerMsgStruct,
			})
		}
	}

}
