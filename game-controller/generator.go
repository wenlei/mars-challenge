package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	maxTemperature       = 35.00
	minTemperature       = -142.00
	maxRadiation         = 1000
	minRadiation         = 0
	variationTemperature = 5.00
	variationRadiation   = 20
	maxTrendSeconds      = 20
	minTrendSeconds      = 5
)

// Reading contains the current sensor readings
type Reading struct {
	SolarFlare         bool    `json:"solarFlare"`
	Temperature        float64 `json:"temperature"`
	Radiation          int     `json:"radiation"`
	temperatureUptrend bool
	radiationUptrend   bool
}

func (s *Reading) updateSolarFlare() {
	x := rand.Intn(2)
	if x != 0 {
		s.SolarFlare = true
	} else {
		s.SolarFlare = false
	}
}

func (s *Reading) updateTemperature() {
	//TODO: consider solar Flare
	var min float64
	var max float64

	if s.temperatureUptrend {
		max = s.Temperature + variationTemperature
		min = s.Temperature
	} else {
		max = s.Temperature
		min = s.Temperature - variationTemperature
	}

	temperature := (rand.Float64() * (max - min)) + min
	if temperature < minTemperature {
		temperature = minTemperature
	} else if temperature > maxTemperature {
		temperature = maxTemperature
	}
	s.Temperature = temperature
}

func (s *Reading) updateTemperatureTrend() {
	ratio := (s.Temperature - minTemperature) / (maxTemperature - minTemperature)
	chance := rand.Float64()
	s.temperatureUptrend = chance > ratio
	log.Printf("[Temperature] Ratio: %.2f, Change: %.2f, Uptrend: %t\n", ratio, chance, s.temperatureUptrend)
}

func (s *Reading) updateRadiation() {
	//TODO: consider solar Flare
	var min int
	var max int

	if s.radiationUptrend {
		max = s.Radiation + variationRadiation
		min = s.Radiation
	} else {
		max = s.Radiation
		min = s.Radiation - variationRadiation
	}

	radiation := rand.Intn(max-min) + min
	if radiation < minRadiation {
		radiation = minRadiation
	} else if radiation > maxRadiation {
		radiation = maxRadiation
	}
	s.Radiation = radiation
}

func (s *Reading) updateRadiationTrend() {
	ratio := (float64)(s.Radiation-minRadiation) / (float64)(maxRadiation-minRadiation)
	chance := rand.Float64()
	s.radiationUptrend = chance > ratio
	log.Printf("[Radiation] Ratio: %.2f, Change: %.2f, Uptrend: %t\n", ratio, chance, s.radiationUptrend)
}

func solarFlareRoutine(wg *sync.WaitGroup, game *GameInfo) {
	timer := time.NewTimer(0)
	for game.Running {
		select {
		case <-timer.C:
			game.Reading.updateSolarFlare()
			if game.Reading.SolarFlare == true {
				timer.Reset(10 * time.Second)
			} else {
				timer.Reset(30 * time.Second)
			}
		}
	}
	wg.Done()
}

func temperatureRoutine(wg *sync.WaitGroup, game *GameInfo) {
	tickerUpdate := time.NewTicker(1 * time.Second)
	timerTrend := time.NewTimer(0)
	for game.Running {
		select {
		case <-tickerUpdate.C:
			game.Reading.updateTemperature()
		case <-timerTrend.C:
			game.Reading.updateTemperatureTrend()
			timerTrend.Reset(time.Duration(rand.Intn(maxTrendSeconds-minTrendSeconds)+minTrendSeconds) * time.Second)
		}
	}
	wg.Done()
}

func radiationRoutine(wg *sync.WaitGroup, game *GameInfo) {
	tickerUpdate := time.NewTicker(1 * time.Second)
	timerTrend := time.NewTimer(0)
	for game.Running {
		select {
		case <-tickerUpdate.C:
			game.Reading.updateRadiation()
		case <-timerTrend.C:
			game.Reading.updateRadiationTrend()
			timerTrend.Reset(time.Duration(rand.Intn(maxTrendSeconds-minTrendSeconds)+minTrendSeconds) * time.Second)
		}
	}
	wg.Done()
}