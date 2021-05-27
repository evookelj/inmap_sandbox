package main

import (
	"context"
	"fmt"
	"github.com/evookelj/inmap/emissions/slca/eieio"
	"github.com/evookelj/inmap/emissions/slca/eieio/eieiorpc"
	"log"
)

func getExposureByPopulation(ctx context.Context, s *eieio.Server, year int32, loc eieiorpc.Location, demand *eieiorpc.Vector) (*map[string]float64, error) {
	vec, err := s.SpatialEIO.Concentrations(ctx, &eieiorpc.ConcentrationInput{
		Demand:    demand,
		Pollutant: eieiorpc.Pollutant_TotalPM25,
		Year:      year,
		Location:  loc,
		AQM:       "isrm",
	})
	conc := vec.Data
	if err != nil {
		return nil, err
	}

	popNames := append(s.CSTConfig.CensusPopColumns, s.CSTConfig.CensusIncomeDecileNames...)
	populationGridsByPopName := make(map[string][]float64)
	for i, popName := range popNames {
			pop, err := s.CSTConfig.PopulationCount(ctx, &eieiorpc.PopulationCountInput{
				Year:        2014, // year,
				Population:  popName,
				AQM:         "isrm",
				IsIncomePop: i >= len(s.CSTConfig.CensusPopColumns), // based off gen of popNames above
			})
			if err != nil {
				return nil, err
			}

			if len(pop) != len(conc) {
				return nil, fmt.Errorf("expected len(population)=len(concentrations); got %d != %d", len(pop), len(conc))
			}
			populationGridsByPopName[popName] = pop
	}

	popTotals := make(map[string]float64)
	for _, pop := range popNames {
		popTotals[pop] = 0
	}

	exposureByPop := make(map[string]float64)
	for gridIdx, concentrationAmt := range conc {
		log.Printf("\t[Grid %d] [Concentration=%.2f]", gridIdx, concentrationAmt)
		for _, popName := range popNames {
			numIndividuals := populationGridsByPopName[popName][gridIdx]
			exposureByPop[popName] += numIndividuals * concentrationAmt
			if popName != s.CSTConfig.CensusTotalPopColumn {
				popTotals[popName] += numIndividuals
			}
			log.Printf("\t\t[Population %s] %.2f ppl --> %.2f exposure", popName, numIndividuals, numIndividuals*concentrationAmt)
		}
	}

	return &exposureByPop, nil
}