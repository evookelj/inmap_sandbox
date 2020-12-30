package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/evookelj/inmap/emissions/slca/eieio"
	"github.com/evookelj/inmap/emissions/slca/eieio/eieiorpc"
	"github.com/evookelj/inmap/epi"
	"github.com/pkg/errors"
	"log"
	"os"
)

const CONFIG = "/Users/emmavukelj/eio_reproduce/data/my_config.toml"

func getEIOServer() (*eieio.Server, error) {
	f, err := os.Open(CONFIG)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg eieio.ServerConfig
	_, err = toml.DecodeReader(f, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.Config.Years = []eieio.Year{2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015}

	return eieio.NewServer(&cfg, "", epi.NasariACS)
}

func getConsumptionBySCC(s *eieio.Server, dem eieiorpc.Demograph, year int32, loc eieiorpc.Location) ([]float64, error) {
	totalConsRPC, err := s.CES.DemographicConsumption(context.Background(), &eieiorpc.DemographicConsumptionInput{
		Year:      year,
		Demograph: dem,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error calculating demographic consumption")
	}

	consumptionBySCC := make([]float64, len(s.SCCs))
	for industryIdx, consumption := range totalConsRPC.Data {
		SCCs := s.IndustryToSCCMap[industryIdx]
		for _, sccIdx := range SCCs {
			consumptionBySCC[sccIdx] += consumption
		}
	}

	return consumptionBySCC, nil
}

func getDemand(s *eieio.Server, year int32, loc eieiorpc.Location) ([]float64, error) {
	demand, err := s.FinalDemand(context.TODO(), &eieiorpc.FinalDemandInput{
		FinalDemandType: eieiorpc.FinalDemandType_AllDemand,
		Year:            year,
		Location:        loc,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting final demand")
	}

	// Step 3: Get emissions for that demand
	emisRPC, err := s.EmissionsMatrix(context.Background(), &eieiorpc.EmissionsMatrixInput{
		Demand:               demand,
		Year:                 int32(year),
		Location:             loc,
		AQM:                  "isrm",
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting emissions matrix")
	}
	emis := rpc2mat(emisRPC)

	r, c := emis.Dims()
	if c != len(s.SCCs) {
		return nil, fmt.Errorf("expected emissions to have #SCC %d columns, got %d", len(s.SCCs), c)
	}
	emisSCC := make([]float64, c)
	for sectorIdx, _ := range s.SCCs {
		var totalEmissions float64 = 0
		for i := 0; i < r; i++ {
			totalEmissions += emis.At(i, sectorIdx)
		}
		emisSCC[sectorIdx] = totalEmissions
	}
	return emisSCC, nil
}

func demAndEmissions() error {
	s, err := getEIOServer()
	if err != nil {
		return errors.Wrap(err, "error creating EIO server")
	}

	var year int32 = 2015
	loc := eieiorpc.Location_Domestic
	emis, err := getDemand(s, year, loc)

	for _, dem := range []eieiorpc.Demograph{eieiorpc.Demograph_Black, eieiorpc.Demograph_Hispanic, eieiorpc.Demograph_WhiteOther} {
		consumption, err := getConsumptionBySCC(s, dem, year, loc)
		if err != nil {
			return errors.Wrap(err, "error getting consumption")
		}

		for sectorIdx, sector := range s.SCCs {
			log.Printf("Sector %v\tDemograph %v --> %v\n", sector, dem, consumption[sectorIdx]*emis[sectorIdx])
		}
	}

	return nil
}

func main() {
	if err := demAndEmissions(); err != nil {
		log.Fatalf(err.Error())
	}
}