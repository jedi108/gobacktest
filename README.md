[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](/LICENSE.md)
[![Travis](https://img.shields.io/travis/mamady/gobacktest.svg?style=flat-square)](https://travis-ci.org/mamady/gobacktest)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/jedi108/gobacktest)
[![Coverage Status](https://img.shields.io/coveralls/mamady/gobacktest/master.svg?style=flat-square)](https://coveralls.io/github/mamady/gobacktest?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedi108/gobacktest?style=flat-square)](https://goreportcard.com/report/github.com/jedi108/gobacktest)

---

# gobacktest - Fundamental stock analysis backtesting

An event-driven backtesting framework to test trading strategies based on fundamental analysis.

## Usage

Example tests are in the `/examples` folder.

```golang
package main

import (
    "github.com/jedi108/gobacktest/pkg/backtest"
    "github.com/jedi108/gobacktest/pkg/data"
    "github.com/jedi108/gobacktest/pkg/strategy"
)

func main() {
    // we need a new blanc backtester
    test := backtest.New()

    // define the symbols to be tested and load them into the backtest
    symbols := []string{"TEST.DE"}
    test.SetSymbols(symbols)

    // create a data provider and load the data into the backtest
    data := &data.BarEventFromCSVFile{FileDir: "../testdata/test/"}
    data.Load(symbols)
    test.SetData(data)

    // set the portfolio with initial cash and a default size and risk manager
    portfolio := &backtest.Portfolio{}
    portfolio.SetInitialCash(10000)

    sizeManager := &backtest.Size{DefaultSize: 100, DefaultValue: 1000}
    portfolio.SetSizeManager(sizeManager)

    riskManager := &backtest.Risk{}
    portfolio.SetRiskManager(riskManager)

    test.SetPortfolio(portfolio)

    // create a strategy provider and load it into the backtest
    strategy := &strategy.Basic{}
    test.SetStrategy(strategy)

    // create an execution provider and load it into the backtest
    exchange := &backtest.Exchange{}
    test.SetExchange(exchange)

    // choose a statistic and load into the backtest
    statistic := &backtest.Statistic{}
    test.SetStatistic(statistic)

    // run the backtest
    test.Run()

    // print the result of the test
    test.Stats().TotalEquityReturn()
}
```

## Dependencies

The internal calculations use the [github.com/shopspring/decimal](https://github.com/shopspring/decimal) package for arbitrary-precision fixed-point decimals.

Make sure to install it into your `$GOPATH` with

    go get github.com/shopspring/decimal

## Basic components

These are the basic components of an event-driven framework.

1. BackTester - general test case, bundles the following elements into a single test
2. EventHandler - the different types of events, which travel through this system - data event, signal event, order event and fill event
3. DataHandler - interface to a set of data, e.g historical quotes, fundamental data, dividends etc.
4. StrategyHandler - generates a buy/sell signal based on the data
5. PortfolioHandler - generates orders and manages profit & loss
    + (SizeHandler) - manages the size of an order
    + (RiskHandler) - manages the risk allocation of a portfolio
6. ExecutionHandler - sends orders to the broker and receives the “fills” or signals that the stock has been bought or sold
7. StatisticHandler - tracks all events during the backtests and calculates useful statistics like equity return, drawdown or sharp ratio etc., could be used to replay the complete backtest for later reference
   + (ComplianceHandler) - tracks and documents all trades to the portfolio for compliance reasons

## Infrastructure example

An overviev of the infrastructure of a complete backtesting and trading environment. Taken from the production roadmap of [QuantRocket](https://www.quantrocket.com/#product-roadmap).

- General
  + API gateway
  + configuration loader
  + logging service
  + cron service
- Data
  + database backup and download service
  + securities master services
  + historical market data service
  + fundamental data service
  + earnings data service
  + dividend data service
  + real-time market data service
  + exchange calendar service
- Strategy
  + performance analysis service - tearsheet
- Portfolio
  + account and portfolio service
  + risk management service
- Execution
  + trading platform gateway service
  + order management and trade ledger service
  + backtesting and trading engine

---

## Resources

### Articles

These links to articles are a good starting point to understand the intentions and basic functions of an event-driven backtesting framework.

- Initial idea via a blog post [Python For Finance: Algorithmic Trading](https://www.datacamp.com/community/tutorials/finance-python-trading#backtesting) by Karlijn Willems [@willems_karlijn](https://twitter.com/willems_karlijn).
- Very good explanation of the internals of a backtesting system by Michael Halls-Moore [@mhallsmoore](https://twitter.com/mhallsmoore) in the blog post series [Event-Driven-Backtesting-with-Python](https://www.quantstart.com/articles/Event-Driven-Backtesting-with-Python-Part-I).

### Other backtesting frameworks

- [QuantConnect](https://www.quantconnect.com)
- [Quantopian](https://www.quantopian.com)
- [QuantRocket](https://www.quantrocket.com) - in development, available Q2/2018
- [Quandl](https://www.quandl.com) - financial data
- [QSTrader](https://www.quantstart.com/qstrader) - open-source backtesting framework from [QuantStart](https://www.quantstart.com)

### General information on Quantitative Finance

- [Quantocracy](http://quantocracy.com) - forum for quant news
- [QuantStart](https://www.quantstart.com) - articels and tutorials about quant finance
