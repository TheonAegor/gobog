package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/TheonAegor/gobog/internal"
	"github.com/TheonAegor/gobog/pkg/tg"
	"github.com/cucumber/godog"
	"os/exec"
)

type (
	botIdCtxKey     struct{}
	prevBotResponse struct{}
)

var prevBotResp = &tg.BotResponse{}

func thereIsBot(handler *internal.Handler) func(context.Context, int64) (context.Context, error) {
	return func(ctx context.Context, botId int64) (context.Context, error) {
		return context.WithValue(ctx, botIdCtxKey{}, botId), nil
	}
}

func iSent(handler *internal.Handler) func(context.Context, string) (context.Context, error) {
	return func(ctx context.Context, s string) (context.Context, error) {
		botId, ok := ctx.Value(botIdCtxKey{}).(int64)
		if !ok {
			return ctx, errors.New("there are no botid")
		}

		_, err := handler.TgClient.SendMessage(ctx, botId, s)
		if err != nil {
			return ctx, errors.New("cannot send message")
		}

		return ctx, nil
	}
}

func iPush(handler *internal.Handler) func(context.Context, string) (context.Context, error) {
	return func(ctx context.Context, s string) (context.Context, error) {
		prevMessage, ok := ctx.Value(prevBotResponse{}).(*tg.BotResponse)
		if !ok {
			prevMessage = prevBotResp
			if prevMessage == nil {
				return ctx, errors.New("there are no prev message")
			}
		}

		button, ok := prevMessage.ButtonsByName[s]
		if !ok {
			return ctx, errors.New(fmt.Sprintf("no button %s found", s))
		}
		_ = handler.TgClient.PushButton(ctx, button)

		return ctx, nil
	}
}

func thereShouldBeButtons(handler *internal.Handler) func(context.Context, int) (context.Context, error) {
	return func(ctx context.Context, numButton int) (context.Context, error) {
		botId, ok := ctx.Value(botIdCtxKey{}).(int64)
		if !ok {
			return ctx, errors.New("there are no botid")
		}

		m, err := handler.TgClient.GetUpdate(ctx, botId)
		if err != nil {
			return ctx, errors.New("error getting update")
		}

		menu, err := tg.ParseMessage(m.Message)
		if err != nil {
			return ctx, err
		}

		if menu.Len != numButton {
			return ctx, fmt.Errorf("expected to get %d button, but have %d", numButton, len(menu.Buttons))
		}

		prevBotResp = menu

		return context.WithValue(ctx, prevBotResponse{}, menu), nil
	}
}

func textShouldMatch(handler *internal.Handler) func(context.Context, *godog.DocString) (context.Context, error) {
	return func(ctx context.Context, body *godog.DocString) (context.Context, error) {
		prevMessage, ok := ctx.Value(prevBotResponse{}).(*tg.BotResponse)
		if !ok {
			prevMessage = prevBotResp
			if prevMessage == nil {
				return ctx, errors.New("there are no prev message")
			}
		}

		text := body.Content
		if prevMessage.Text != text {
			return ctx, errors.New(fmt.Sprintf("get unexpected text, get: %s; want: %s", prevMessage.Text, text))
		}

		// prevMessage stays the same
		return ctx, nil
	}
}

func resedDb(handler *internal.Handler) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		cmd := exec.Command("/bin/sh", "reset_db.sh")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return ctx, errors.New(fmt.Sprint(err) + ": " + stderr.String())
		}
		handler.Opts.Log.Info(ctx, "Result: "+out.String())

		return ctx, nil
	}
}

func InitializeScenario(sc *godog.ScenarioContext, handler *internal.Handler) {
	sc.Step(`^there is bot (\d+)$`, thereIsBot(handler))
	sc.Step(`^I sent "([^"]*)"$`, iSent(handler))
	sc.Step(`^there should be (\d+) buttons$`, thereShouldBeButtons(handler))
	sc.Step(`^I push "([^"]*)"$`, iPush(handler))
	sc.Step(`^the text should be:$`, textShouldMatch(handler))
	sc.Step(`^here we go again$`, resedDb(handler))
}
