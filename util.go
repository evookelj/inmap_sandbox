package main

import (
	"github.com/BurntSushi/toml"
	"github.com/evookelj/inmap/emissions/slca/eieio"
	"github.com/evookelj/inmap/emissions/slca/eieio/eieiorpc"
	"github.com/evookelj/inmap/epi"
	"gonum.org/v1/gonum/mat"
	"os"
)

var CONFIG = os.ExpandEnv("${INMAP_SANDBOX_ROOT}/data/my_config.toml")

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

func array2vec(d []float64) *mat.VecDense {
if len(d) == 0 {
return nil
}
return mat.NewVecDense(len(d), d)
}

func rpc2vec(d *eieiorpc.Vector) *mat.VecDense {
if d == nil {
return nil
}
return array2vec(d.Data)
}

func vec2array(v *mat.VecDense) []float64 {
if v == nil {
return nil
}
return v.RawVector().Data
}

func mask2rpc(m *eieio.Mask) *eieiorpc.Mask {
if m == nil {
return nil
}
return &eieiorpc.Mask{Data: vec2array((*mat.VecDense)(m))}
}

func rpc2mask(m *eieiorpc.Mask) *eieio.Mask {
if m == nil {
return nil
}
return (*eieio.Mask)(array2vec(m.Data))
}

func vec2rpc(v *mat.VecDense) *eieiorpc.Vector {
return &eieiorpc.Vector{Data: vec2array(v)}
}

func mat2rpc(m *mat.Dense) *eieiorpc.Matrix {
r, c := m.Dims()
return &eieiorpc.Matrix{Rows: int32(r), Cols: int32(c), Data: m.RawMatrix().Data}
}

func rpc2mat(m *eieiorpc.Matrix) *mat.Dense {
return mat.NewDense(int(m.Rows), int(m.Cols), m.Data)
}