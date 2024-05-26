package main

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/uncle-gua/gobinance/futures"
	"github.com/uncle-gua/log"
)

func main() {
	if err := goodluck.Go(); err != nil {
		log.Fatal(err)
	}
}

var goodluck = &GoodLuck{}

type GoodLuck struct {
	client *futures.Client
}

func (g *GoodLuck) Go() error {
	g.client = futures.NewClient(config.ApiKey, config.ApiSecret)

	info, err := g.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	symbols := make([]futures.Symbol, 0)
	for _, s := range info.Symbols {
		if s.QuoteAsset == "USDT" && s.ContractType == "PERPETUAL" && s.Status == "TRADING" {
			symbols = append(symbols, s)
		}
	}

	symbol := func() futures.Symbol {
		count := len(symbols)
		x := rand.Intn(count)
		return symbols[x]
	}()

	side1, side2, positionSide := func() (futures.SideType, futures.SideType, futures.PositionSideType) {
		x := rand.Intn(2)
		if x == 0 {
			return futures.SideTypeBuy, futures.SideTypeSell, futures.PositionSideTypeLong
		}

		return futures.SideTypeSell, futures.SideTypeBuy, futures.PositionSideTypeShort
	}()

	duration := func() time.Duration {
		x := rand.Intn(config.Duration) + 1
		return time.Duration(x) * time.Minute
	}()

	log.Infof("duration: %s", duration)

	price, err := g.getPrice(symbol)
	if err != nil {
		return err
	}

	quantity, err := func(s string) (string, error) {
		price, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return "", err
		}

		filter := symbol.LotSizeFilter()
		if filter == nil {
			return "", errors.New("can get lot size filter")
		}

		p := int(math.Log10(1 / filter.StepSize))
		return strconv.FormatFloat(config.Amount/price, 'f', p, 64), nil
	}(price)
	if err != nil {
		return err
	}

	_, err = g.client.NewCreateOrderService().
		Symbol(symbol.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side1).
		PositionSide(positionSide).
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		return err
	}

	log.Infof("BUY %s, price: %s, quantity: %s", futures.PositionSideTypeLong, price, quantity)

	time.Sleep(duration)

	price, err = g.getPrice(symbol)
	if err != nil {
		return err
	}

	_, err = g.client.NewCreateOrderService().
		Symbol(symbol.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side2).
		PositionSide(positionSide).
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		return err
	}

	log.Infof("SELL %s, price: %s, quantity: %s", futures.PositionSideTypeLong, price, quantity)

	return nil
}

func (g *GoodLuck) getPrice(symbol futures.Symbol) (string, error) {
	prices, err := g.client.NewListPricesService().Symbol(symbol.Symbol).Do(context.Background())
	if err != nil {
		return "", err
	}
	if len(prices) == 0 {
		return "", errors.New("can not get price")
	}

	return prices[0].Price, nil
}
