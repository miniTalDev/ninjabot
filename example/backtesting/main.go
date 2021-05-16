package main

import (
	"context"

	"github.com/rodrigo-brito/ninjabot"
	"github.com/rodrigo-brito/ninjabot/example"
	"github.com/rodrigo-brito/ninjabot/pkg/exchange"
	"github.com/rodrigo-brito/ninjabot/pkg/model"
	"github.com/rodrigo-brito/ninjabot/pkg/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	settings := model.Settings{
		Pairs: []string{
			"BTCUSDT",
		},
	}

	strategy := new(example.MyStrategy)

	csvFeed, err := exchange.NewCSVFeed(
		strategy.Timeframe(),
		exchange.PairFeed{
			Pair:      "BTCUSDT",
			File:      "testdata/btc-1h.csv",
			Timeframe: "1h",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.New("backtest.db")
	if err != nil {
		log.Fatal(err)
	}

	wallet := exchange.NewPaperWallet(
		ctx,
		"USDT",
		exchange.WithPaperAsset("USDT", 10000),
		exchange.WithDataFeed(csvFeed),
	)

	bot, err := ninjabot.NewBot(
		ctx,
		settings,
		wallet,
		strategy,
		ninjabot.WithStorage(storage),
		ninjabot.WithCandleSubscription(wallet),
		ninjabot.WithLogLevel(log.ErrorLevel),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print bot results
	bot.Summary()
	wallet.Summary()
}
