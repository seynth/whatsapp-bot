package auxiliary

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func AphrCommandHandler(aphr *whatsmeow.Client, evt *events.Message) {
	extendedMsg := evt.Message.GetExtendedTextMessage()

	if extendedMsg != nil {
		caption := extendedMsg.GetText()

		switch strings.TrimSpace(caption) {
		case "#s2i":
			jpegBytes := make(chan []byte)
			defer close(jpegBytes)

			mp4Bytes := make(chan []byte)
			defer close(mp4Bytes)

			fmt.Println(extendedMsg)
			if quotedMsg := extendedMsg.ContextInfo.GetQuotedMessage(); quotedMsg.StickerMessage != nil {

				isAnim := quotedMsg.StickerMessage.IsAnimated

				replyCommand := &waE2E.ExtendedTextMessage{
					Text: proto.String("*[Aphrodite]* Converting to media wait . . ."),
					ContextInfo: &waE2E.ContextInfo{
						QuotedMessage: &waE2E.Message{
							StickerMessage: &waE2E.StickerMessage{
								URL:           quotedMsg.StickerMessage.URL,
								FileSHA256:    quotedMsg.StickerMessage.FileSHA256,
								FileEncSHA256: quotedMsg.StickerMessage.FileEncSHA256,
								FileLength:    quotedMsg.StickerMessage.FileLength,
								MediaKey:      quotedMsg.StickerMessage.MediaKey,
								DirectPath:    quotedMsg.StickerMessage.DirectPath,
							},
						},
					},
				}

				unsupportedReply := &waE2E.ExtendedTextMessage{
					Text: proto.String("*[Aphrodite]* This animated sticker is not supported . . ."),
					ContextInfo: &waE2E.ContextInfo{
						QuotedMessage: &waE2E.Message{
							StickerMessage: &waE2E.StickerMessage{
								URL:           quotedMsg.StickerMessage.URL,
								FileSHA256:    quotedMsg.StickerMessage.FileSHA256,
								FileEncSHA256: quotedMsg.StickerMessage.FileEncSHA256,
								FileLength:    quotedMsg.StickerMessage.FileLength,
								MediaKey:      quotedMsg.StickerMessage.MediaKey,
								DirectPath:    quotedMsg.StickerMessage.DirectPath,
							},
						},
					},
				}

				sticker, errSticker := aphr.DownloadAny(quotedMsg)
				if errSticker != nil {
					log.Fatalln("Download sticker err: ", errSticker)
				}

				if !*isAnim {
					aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
						ExtendedTextMessage: replyCommand,
					})

					go AphrConvertWebp2Jpeg(sticker, jpegBytes)

					resSticker, errUpload := aphr.Upload(context.Background(), <-jpegBytes, whatsmeow.GetMediaType(evt.Message.StickerMessage))
					if errUpload != nil {
						log.Fatalln("Error upload sticker: ", errUpload)
					}

					imgMsg := &waE2E.ImageMessage{
						URL:           &resSticker.URL,
						DirectPath:    &resSticker.DirectPath,
						FileSHA256:    resSticker.FileSHA256,
						FileLength:    &resSticker.FileLength,
						FileEncSHA256: resSticker.FileEncSHA256,
						Mimetype:      proto.String("image/jpeg"),
						Caption:       proto.String("*[Aphrodite]* Done, converted to image"),
						MediaKey:      resSticker.MediaKey,
						ContextInfo: &waE2E.ContextInfo{
							QuotedMessage: &waE2E.Message{
								StickerMessage: &waE2E.StickerMessage{
									URL:           quotedMsg.StickerMessage.URL,
									FileSHA256:    quotedMsg.StickerMessage.FileSHA256,
									FileEncSHA256: quotedMsg.StickerMessage.FileEncSHA256,
									FileLength:    quotedMsg.StickerMessage.FileLength,
									MediaKey:      quotedMsg.StickerMessage.MediaKey,
									DirectPath:    quotedMsg.StickerMessage.DirectPath,
								},
							},
						},
					}

					aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
						ImageMessage: imgMsg,
					})
				} else {
					aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
						ExtendedTextMessage: unsupportedReply,
					})

				}

			} else {
				aphr.SendMessage(context.Background(), evt.Info.Sender, &waE2E.Message{
					Conversation: proto.String("Error: Not a sticker"),
				})
			}
		}
	}

}
