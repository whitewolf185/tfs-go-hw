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
	start := false

	canChan := make(chan domain.Candle)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for val := range price {
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
					Logger.Info("created new candle ", canPer)
					candle = CreateNewCandle(val, canPer)
				} else {
					candle.Low = MinVal(candle.Low, val.Value)
					candle.High = MaxVal(candle.High, val.Value)
				}

				candleMap[val.Ticker] = candle
			}
		}
		if start {
			for _, candle := range candleMap {
				canChan <- candle
			}
		}
		Logger.Info("done")
		close(canChan)
	}()

	return canChan
}

func WriteToFile(wg *sync.WaitGroup,
	value <-chan domain.Candle, canPer domain.CandlePeriod) error {
	defer wg.Done()

	var (
		csvFile *os.File
		err     error
	)

	path := fmt.Sprintf("D:/Documents/tfs-go-hw/lesson3/table/Candles%s.csv", canPer)
	csvFile, err = os.Create(path)
	if err != nil {
		Logger.Error("Bad open file ", canPer)
		return badOpenErr(canPer)
	}

	csvWriter := csv.NewWriter(csvFile)

	for v := range value {
		var strPer string
		strPer, err = GetCanPerStr(canPer)

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
			return showWriterErr(canPer)
		}
		csvWriter.Flush()
		info := log.New()
		info.Info("Writing to the file was successful")
	}

	Logger.Info("closing file ", canPer)
	err = csvFile.Close()
	if err != nil {
		return badCloseErr(canPer)
	}
	return nil
}
