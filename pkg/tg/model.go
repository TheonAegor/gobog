package tg

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
)

var (
	NoButtonsErr              = errors.New("no button is message")
	UnknownButtonTypeErr      = errors.New("unknown button type")
	UnknownReplyMarkupTypeErr = errors.New("unknown reply markup type")
)

type CallBack struct {
	Text      string
	MessageId int64
	ChatId    int64
	Data      []byte
}

func (c CallBack) ToGetCallbackQueryAnswerRequest() *client.GetCallbackQueryAnswerRequest {
	payload := &client.CallbackQueryPayloadData{Data: c.Data}

	return &client.GetCallbackQueryAnswerRequest{
		ChatId:    c.ChatId,
		MessageId: c.MessageId,
		Payload:   payload,
	}
}

type PushButtonResponse struct {
	callbackQueryAnswer *client.CallbackQueryAnswer
	err                 error
}

type UserButton struct {
	Text   string
	UserId int64
}

type UrlButton struct {
	Text string
	Url  string
}

type BotResponse struct {
	Text          string               `json:"text"`
	Buttons       map[int]*CallBack    `json:"buttons"`
	ButtonsByName map[string]*CallBack `json:"buttons_by_name"`
	UserButtons   map[int]*UserButton  `json:"user_buttons"`
	UrlButtons    []*UrlButton         `json:"url_buttons"`

	ButtonNames []string `json:"button_names"`
	Len         int
}

// TODO: rename func
func ParseMessage(msg *client.Message) (*BotResponse, error) {
	if msg.ReplyMarkup == nil {
		return nil, NoButtonsErr
	}

	msgText, ok := msg.Content.(*client.MessageText)
	if !ok {
		// log warning
	}

	var msgTxt string
	if msgText != nil && msgText.Text != nil {
		msgTxt = msgText.Text.Text
	}

	menu := &BotResponse{
		Buttons:       make(map[int]*CallBack),
		UserButtons:   make(map[int]*UserButton),
		ButtonsByName: make(map[string]*CallBack),
		Text:          msgTxt,
	}

	buttons, ok := msg.ReplyMarkup.(*client.ReplyMarkupInlineKeyboard)
	if !ok {
		return nil, UnknownReplyMarkupTypeErr
	}

	var (
		button     *client.InlineKeyboardButtonTypeCallback
		userButton *client.InlineKeyboardButtonTypeUser
		urlButton  *client.InlineKeyboardButtonTypeUrl
	)

	cnt := 0
	cntUser := 0
	for _, row := range buttons.Rows {
		for _, ik := range row {
			switch ik.Type.InlineKeyboardButtonTypeType() {
			case client.TypeInlineKeyboardButtonTypeCallback:
				button = ik.Type.(*client.InlineKeyboardButtonTypeCallback)
				callBack := &CallBack{
					Text:      ik.Text,
					MessageId: msg.Id,
					ChatId:    msg.ChatId,
					Data:      button.Data,
				}

				menu.Buttons[cnt] = callBack
				menu.ButtonsByName[callBack.Text] = callBack
				cnt++
			case client.TypeInlineKeyboardButtonTypeUser:
				userButton = ik.Type.(*client.InlineKeyboardButtonTypeUser)
				menu.UserButtons[cntUser] = &UserButton{UserId: userButton.UserId, Text: ik.Text}
				cntUser++
			case client.TypeInlineKeyboardButtonTypeUrl:
				urlButton = ik.Type.(*client.InlineKeyboardButtonTypeUrl)
				menu.UrlButtons = append(menu.UrlButtons, &UrlButton{
					Text: ik.Text,
					Url:  urlButton.Url,
				})
			default:
				return nil, UnknownButtonTypeErr
			}
			menu.ButtonNames = append(menu.ButtonNames, ik.Text)
		}
	}

	menu.Len = len(menu.ButtonNames)

	return menu, nil
}
