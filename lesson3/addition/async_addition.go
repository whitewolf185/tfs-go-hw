package addition

import (
	"context"
	"encoding/csv"
	"hw-async/domain"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetData(ctx context.Context, prices <-chan domain.Price, wg *sync.WaitGroup) {
	defer wg.Done()

	C1m := make(chan domain.Price)
	C2m := make(chan domain.Price)
	C10m := make(chan domain.Price)

	wg.Add(3)
	go createCandle(ctx, wg, C1m, domain.CandlePeriod1m)
	go createCandle(ctx, wg, C2m, domain.CandlePeriod2m)
	go createCandle(ctx, wg, C10m, domain.CandlePeriod10m)

	dataLogger := log.New()
	for {
		select {
		case <-ctx.Done():
			dataLogger.Info("done")
			return
		case v := <-prices:
			dataLogger.Info(v)
			C1m <- v
			C2m <- v
			C10m <- v
		}
	}
}

func createCandle(ctx context.Context, wg *sync.WaitGroup,
	price <-chan domain.Price, canPer domain.CandlePeriod) {
	defer wg.Done()
	dataLogger := log.New()
	var (
		err   error
		timer time.Time
	)
	candleMap := make(map[string]domain.Candle)
	start := false

	wg.Add(1)
	canChan := make(chan domain.Candle)
	go WriteToFile(ctx, wg, canChan, canPer)

	for {
		select {
		case <-ctx.Done():
			if start {
				for _, candle := range candleMap {
					canChan <- candle
				}
			}
			Logger.Info("done")
			return
		case val := <-price:
			candle, ok := candleMap[val.Ticker]
			start = true
			if !ok {
				// только создал свечку
				candle = CreateNewCandle(val, canPer)

				candleMap[val.Ticker] = candle
			} else {
				timer, err = domain.PeriodTS(canPer, val.TS)
				if err != nil {
					Logger.Error("in PariodTS ", err)
				}
				comparison, errInt := domain.PeriodTSInt(canPer)
				if errInt != nil {
					Logger.Error("in PariodTSInt ", err)
				}

				if timer.Minute()-candle.TS.Minute() >= comparison { // закрытие свечи
					canChan <- candle
					dataLogger.Info("created new candle ", canPer)
					candle = CreateNewCandle(val, canPer)
				} else {
					candle.Low = MinVal(candle.Low, val.Value)
					candle.High = MaxVal(candle.High, val.Value)
				}

				candleMap[val.Ticker] = candle
			}
		}
	}
}

func WriteToFile(ctx context.Context, wg *sync.WaitGroup,
	value <-chan domain.Candle, canPer domain.CandlePeriod) {
	defer wg.Done()
	var (
		csvFile *os.File
		err     error
	)

	switch canPer {
	case domain.CandlePeriod1m:
		csvFile, err = os.Create("D:/Documents/tfs-go-hw/lesson3/table/Candles1m.csv")
		if err != nil {
			Logger.Error("Bad open file 1m")
			panic(err)
		}
	case domain.CandlePeriod2m:
		csvFile, err = os.Create("D:/Documents/tfs-go-hw/lesson3/table/Candles2m.csv")
		if err != nil {
			Logger.Error("Bad open file 2m")
			panic(err)
		}
	case domain.CandlePeriod10m:
		csvFile, err = os.Create("D:/Documents/tfs-go-hw/lesson3/table/Candles10m.csv")
		if err != nil {
			Logger.Error("Bad open file 10m")
			panic(err)
		}
	}

	defer func(csvFile *os.File) {
		Logger.Info("closing file ", canPer)
		err = csvFile.Close()
		if err != nil {
			Logger.Error("Bad close file ", canPer)
			panic(err)
		}
	}(csvFile)

	csvWriter := csv.NewWriter(csvFile)

	for {
		select {
		case <-ctx.Done():
			for i := 0; i < 4; i++ {
				v := <-value
				result := candleToStr(v, canPer)

				err := csvWriter.Write(result)
				if err != nil {
					panic(err)
				}
				csvWriter.Flush()
				info := log.New()
				info.Info("Writing to the file was successful")
			}
			return
		case v := <-value:
			result := candleToStr(v, canPer)

			err := csvWriter.Write(result)
			if err != nil {
				panic(err)
			}
			csvWriter.Flush()
			info := log.New()
			info.Info("Writing to the file was successful")
		}
	}
}
