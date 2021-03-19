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

	populationNamesOutput, err := s.Populations(ctx, nil)
	if err != nil {
		return nil, err
	}
	popNames := populationNamesOutput.Names

	populationGridsByPopName := make(map[string][]float64)
	for _, popName := range popNames {
		popOutputStruct, err := s.CSTConfig.PopulationIncidence(ctx, &eieiorpc.PopulationIncidenceInput{
			Year:       year,
			Population: popName,
			// these two don't matter b/c we just care about population count
			// TODO: Export method that just gets pop counts, don't waste computing on incidence
			HR:         "NasariACS",
			AQM:        "isrm",
		})
		if err != nil {
			return nil, err
		}

		pop := popOutputStruct.GetPopulation()
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
			popTotals[popName] += numIndividuals
			exposureByPop[popName] += numIndividuals * concentrationAmt
			log.Printf("\t\t[Population %s] %.2f ppl --> %.2f exposure", popName, numIndividuals, numIndividuals * concentrationAmt)
		}
	}

	return &exposureByPop, nil
}