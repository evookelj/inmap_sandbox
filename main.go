package main

import (
	"context"
	"github.com/evookelj/inmap/emissions/slca/eieio/eieiorpc"
	"github.com/pkg/errors"
	"log"
)

const YEAR int32 = 2015
const LOC = eieiorpc.Location_Domestic

func mainHelper() error {
	ctx := context.Background()

	s, err := getEIOServer()
	if err != nil {
		return errors.Wrap(err, "error creating EIO server")
	}

	demand, err := s.FinalDemand(ctx, &eieiorpc.FinalDemandInput{
		FinalDemandType: eieiorpc.FinalDemandType_AllDemand,
		Year:            YEAR,
		Location:        LOC,
	})
	if err != nil {
		return errors.Wrap(err, "error getting final demand")
	}

	/*err = contributionSideTest(ctx, s, YEAR, LOC, demand)
	if err != nil {
		return err
	}*/

	exposureByPop, err := getExposureByPopulation(ctx, s, YEAR, LOC, demand)
	if err != nil {
		return err
	}
	for popName, exposure := range *exposureByPop {
		log.Printf("Pop name: %s\tExposure: %.2f", popName, exposure)
	}

	return nil
}

func main() {
	err := mainHelper()
	if err != nil {
		log.Fatalf(err.Error())
	}
}