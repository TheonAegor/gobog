package tg

import (
	"context"
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"path/filepath"
	"time"
)

const (
	responseTtl = 60
)

//var (
//	DefaultTgClient *client.Client
//)
//
//func init() {
//	log := logger.DefaultLogger
//	ctx := context.Background()
//	cfg := config.DefaultConfig
//
//	log.Info(ctx, "Start telegram client authorization")
//	authorizer := client.ClientAuthorizer()
//	go client.CliInteractor(authorizer)
//
//	authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
//		UseTestDc:              false,
//		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
//		FilesDirectory:         filepath.Join(".tdlib", "files"),
//		UseFileDatabase:        true,
//		UseChatInfoDatabase:    true,
//		UseMessageDatabase:     true,
//		UseSecretChats:         false,
//		ApiId:                  cfg.ApiId,
//		ApiHash:                cfg.ApiHash,
//		SystemLanguageCode:     "en",
//		DeviceModel:            "Server",
//		SystemVersion:          "1.0.0",
//		ApplicationVersion:     "1.0.0",
//		EnableStorageOptimizer: true,
//		IgnoreFileNames:        false,
//	}
//
//	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
//		NewVerbosityLevel: 1,
//	})
//	if err != nil {
//		log.Fatal(ctx, "SetLogVerbosityLevel error", "error", err)
//	}
//	tdlibClient, err := client.NewClient(authorizer)
//	if err != nil {
//		log.Fatal(ctx, "NewClient error", "error", err)
//	}
//
//	DefaultTgClient = tdlibClient
//}

type TgClient struct {
	tdClient *client.Client
	listner  *client.Listener
	opts     *Options

	responseTtl int
}

func (c *TgClient) SendMessage(ctx context.Context, chatId int64, text string) (*client.Message, error) {
	sentMsg, err := c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId:          chatId,
		MessageThreadId: 0,
		ReplyTo:         nil,
		Options:         nil,
		ReplyMarkup:     nil,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text:     text,
				Entities: nil,
			},
			ClearDraft: false,
		},
	})
	if err != nil {
		c.opts.log.Error(ctx, "error sending message")
		return nil, err
	}

	return sentMsg, err
}

func (c *TgClient) PushButton(ctx context.Context, callback *CallBack) <-chan PushButtonResponse {
	qaChan := make(chan PushButtonResponse, 1)

	go func() {
		defer close(qaChan)

		qa, err := c.tdClient.GetCallbackQueryAnswer(callback.ToGetCallbackQueryAnswerRequest())
		if err != nil {
			c.opts.log.Error(ctx, "error pushing button")
		}

		qaChan <- PushButtonResponse{
			callbackQueryAnswer: qa,
			err:                 err,
		}
	}()

	return qaChan
}

func (c *TgClient) GetUpdate(ctx context.Context, botId int64) (*client.UpdateNewMessage, error) {
	updates := c.listner.Updates

	maxWait := time.After(time.Second * responseTtl)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-maxWait:
			return nil, fmt.Errorf("ttl end")
		case u := <-updates:
			switch u.GetType() {
			// UpdateChatLastMessage
			case client.TypeUpdateNewMessage:
				message := u.(*client.UpdateNewMessage)

				if message.Message.ChatId != botId {
					continue
				}

				if message.Message.SenderId.MessageSenderType() == client.TypeMessageSenderUser {
					sender := message.Message.SenderId.(*client.MessageSenderUser)
					if sender.UserId != botId {
						c.opts.log.Debug(ctx, "get message not from chat", "sender", sender)
						continue
					}
				} else {
					c.opts.log.Debug(ctx, "get message with unknown senderId type",
						"type", message.Message.SenderId.MessageSenderType())
					continue
				}

				var content string
				if message.Message.Content.MessageContentType() == client.TypeMessageText {
					messageContent := message.Message.Content.(*client.MessageText)

					content = messageContent.Text.Text
				}

				c.opts.log.Info(ctx, "Get update", "type", client.TypeUpdateNewMessage, "content", content)
				//we need only one message from chat - bot reply
				return message, nil
			}
		}
	}
}

func NewTgClient(apiId int32, apiHash string, opts ...Option) (*TgClient, error) {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  apiId,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		return nil, err
	}
	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		return nil, err
	}

	listner := tdlibClient.GetListener()
	tgClient := TgClient{tdClient: tdlibClient, listner: listner, responseTtl: responseTtl, opts: &Options{}}

	for _, o := range opts {
		o(tgClient.opts)
	}

	return &tgClient, err
}
