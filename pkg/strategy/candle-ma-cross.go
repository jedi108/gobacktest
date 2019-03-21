package strategy

import (
	"fmt"
	bt "github.com/jedi108/gobacktest/pkg/backtest"
	"log"
)

// MovingAverageCross is a test strategy, which interprets the SMA on a series of data events
// specified by ShortWindow (SW) and LongWindow (LW).
// If SW bigger tha LW and there is not already an invested BOT position, the strategy creates a buy signal.
// If SW falls below LW and there is an invested BOT position, the strategy creates an exit signal.
type MovingAverageCrossWithCandle struct {
	ShortWindow int
	LongWindow  int
}

// CalculateSignal handles the single Event
func (s *MovingAverageCrossWithCandle) CalculateSignal(e bt.DataEventHandler, data bt.DataHandler, p bt.PortfolioHandler) (bt.SignalEvent, error) {
	// create empty Signal
	se := &bt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case *bt.Bar:
		// calculate and set SMA for short window
		smaShort, err := bt.CalculateSMA(s.ShortWindow, data.List(e.GetSymbol()))
		if err != nil {
			return se, err
		}
		e.Metrics[fmt.Sprintf("SMA%d", s.ShortWindow)] = smaShort

		// calculate and set SMA for long window
		smaLong, err := bt.CalculateSMA(s.LongWindow, data.List(e.GetSymbol()))
		if err != nil {
			return se, err
		}
		e.Metrics[fmt.Sprintf("SMA%d", s.LongWindow)] = smaLong

		// check if already invested
		_, invested := p.IsInvested(e.GetSymbol())

		if (smaShort > smaLong) && invested {
			return se, fmt.Errorf("buy signal but already invested in %v, no signal created,", e.GetSymbol())
		}

		if IsSignalLong(smaShort, smaLong, invested, data.List(e.GetSymbol())) {
			// buy signal, populate the signal event
			se.Event = bt.Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
			se.Direction = "long"
		}

		if IsSignalLongExit(smaShort, smaLong, invested, data.List(e.GetSymbol())) {
			//if (smaShort <= smaLong) && !invested {
			return se, fmt.Errorf("sell signal but not invested in %v, no signal created,", e.GetSymbol())
		}

		if (smaShort <= smaLong) && invested {
			// sell signal, populate the signal event
			se.Event = bt.Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
			se.Direction = "exit"
		}

	}
	return se, nil
}

func IsSignalLongExit(smaShort, smaLong float64, invested bool, da []bt.DataEventHandler) (hasSignal bool) {
	if !invested {
		return false
	}

	//if smaShort <= smaLong {
	//	return true
	//}

	lenData := len(da)
	if lenData < 2 {
		return false
	}

	_, last := firstLast(lenData, da)

	if !bullish(last) {
		return true
	}

	return false
}

func IsSignalLong(smaShort, smaLong float64, invested bool, da []bt.DataEventHandler) (hasSignal bool) {
	if invested {
		return false
	}

	//hasSignal = shortMoreLong(smaShort, smaLong)
	//
	//if !hasSignal {
	//	return false
	//}

	return isAbsorptionPattern(da)
}

func shortMoreLong(smaShort, smaLong float64) bool {
	return smaShort > smaLong
}

func firstLast(lenData int, da []bt.DataEventHandler) (first barses, last barses) {
	return da[lenData-2].(barses), da[lenData-1].(barses)
}

func isAbsorptionPattern(da []bt.DataEventHandler) bool {
	lenData := len(da)
	if lenData < 2 {
		return false
	}

	first, last := firstLast(lenData, da)

	if !bullish(last) {
		return false
	}

	return true

	if !bullishAbsorptionPattern(first, last) {
		return false
	}

	return true
}

func bullish(last barses) bool {

	return true

	return last.BarOpen() < last.BarClose()
}

func bullishAbsorptionPattern(first, last barses) (result bool) {

	return true

	result = last.BarOpen() <= first.BarClose() && last.BarClose() >= first.BarOpen()
	if result {
		log.Println("result true==============")
		log.Printf("last.BarOpen %f <= %f first.BarClose", last.BarOpen(), first.BarClose())
		log.Printf("last.BarClose %f >= %f first.BarOpen", last.BarClose(), first.BarOpen())
	}
	return result
}

type barses interface {
	BarClose() float64
	BarOpen() float64
	BarHigh() float64
	BarLow() float64
}
