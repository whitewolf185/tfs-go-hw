package addition

import (
	"encoding/csv"
	"fmt"
	"hw-async/domain"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetData(wg *sync.WaitGroup, prices <-chan domain.Price) (<-chan domain.Price, <-chan domain.Price, <-chan domain.Price) {
	C1m := make(chan domain.Price)
	C2m := make(chan domain.Price)
	C10m := make(chan domain.Price)

	wg.Add(1)
	go func() {
		defer close(C1m)
		defer close(C2m)
		defer close(C10m)
		defer wg.Done()

		for v := range prices {
			Logger.Info(v)
			C1m <- v
			C2m <- v
			C10m <- v
		}
	}()

	return C1m, C2m, C10m
}

func CreateCandle(wg *sync.WaitGroup,
	price <-chan domain.Price, canPer domain.CandlePeriod) <-chan domain.Candle {
	var (
		err   error
		timer time.Time
	)
	candleMap := make(map[string]domain.Candle)

	canChan := make(chan domain.Candle)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for val := range price {
			candle, ok := candleMap[val.Ticker]
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
					Logger.Info("created new candle ", canPer)
					candle = CreateNewCandle(val, canPer)
				} else {
					candle.Low = MinVal(candle.Low, val.Value)
					candle.High = MaxVal(candle.High, val.Value)
				}

				candleMap[val.Ticker] = candle
			}
		}

		for _, candle := range candleMap {
			canChan <- candle
		}

		Logger.Info("done")
		close(canChan)
	}()

	return canChan
}

func WriteToFile(wg *sync.WaitGroup,
	value <-chan domain.Candle, canPer domain.CandlePeriod) {
	var (
		csvFile *os.File
		err     error
	)

	path := fmt.Sprintf("D:/GO/tfs-go-hw/lesson3/table/candles_%s.csv", canPer)
	csvFile, err = os.Create(path)
	if err != nil {
		Logger.Error("Bad open file ", canPer, " Error: ", err)
		return
	}

	csvWriter := csv.NewWriter(csvFile)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			Logger.Info("closing file ", canPer)
			err = csvFile.Close()
			if err != nil {
				Logger.Error("Bad close file ", canPer, " Error: ", err)
				return
			}
		}()
		for v := range value {
			var strPer string
			strPer, err = GetCanPerStr(canPer)
			if err != nil {
				Logger.Error(canPer, " Error: ", err)
				return
			}

			err = csvWriter.Write([]string{
				v.Ticker,
				fmt.Sprintf("%s", strPer),
				fmt.Sprintf("%f", v.Open),
				fmt.Sprintf("%f", v.High),
				fmt.Sprintf("%f", v.Low),
				fmt.Sprintf("%f", v.Close),
				v.TS.Format(time.RFC3339),
			})
			if err != nil {
				Logger.Error("Bad write file ", canPer, " Error: ", err)
				return
			}
			csvWriter.Flush()
			info := log.New()
			info.Info("Writing to the file was successful")
		}
	}()
}
