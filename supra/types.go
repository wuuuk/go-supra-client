package supra

import "github.com/wuuuk/go-bcs/bcs"

type ProofData struct {
	Decimals   [][]uint16       `json:"decimals"`
	Pairs      [][]uint32       `json:"pairs"`
	PairMasks  [][]bool         `json:"pairMasks"`
	Timestamps [][]*bcs.Uint128 `json:"timestamps"`
	Prices     [][]*bcs.Uint128 `json:"prices"`
}
