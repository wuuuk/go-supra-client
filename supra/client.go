package supra

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/wuuuk/go-bcs/bcs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SupraClient struct {
	grpcClient        *grpc.ClientConn
	PullServiceClient PullServiceClient
	Context           context.Context
}

func InitSupraClient(ctx context.Context, grpcAddress string) (*SupraClient, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial %v", err)
	}

	pullServiceClient := NewPullServiceClient(conn)

	return &SupraClient{
		grpcClient:        conn,
		PullServiceClient: pullServiceClient,
		Context:           ctx,
	}, nil
}

func (supra *SupraClient) GetProofOnSui(indexes []uint32) (*PullResponseSui, error) {
	result, err := supra.PullServiceClient.GetProof(supra.Context, &PullRequest{
		PairIndexes: indexes,
		ChainType:   "sui",
	})
	if err != nil {
		return nil, fmt.Errorf("pullServiceClient.GetProof %v", err)
	}

	data, ok := result.Resp.(*PullResponse_Sui)
	if !ok {
		return nil, fmt.Errorf("response type not match")
	}
	return data.Sui, nil
}

func (response *PullResponseSui) DecodeProofData() (data *ProofData, err error) {
	data = new(ProofData)
	_, err = bcs.Unmarshal(response.SccDecimals, &data.Decimals)
	if err != nil {
		return nil, fmt.Errorf("bcs.Unmarshal SccDecimals %v", err)
	}
	_, err = bcs.Unmarshal(response.SccPair, &data.Pairs)
	if err != nil {
		return nil, fmt.Errorf("bcs.Unmarshal SccPair %v", err)
	}

	_, err = bcs.Unmarshal(response.PairMask, &data.PairMasks)
	if err != nil {
		return nil, fmt.Errorf("bcs.Unmarshal PairMask %v", err)
	}

	_, err = bcs.Unmarshal(response.SccTimestamp, &data.Timestamps)
	if err != nil {
		return nil, fmt.Errorf("bcs.Unmarshal SccTimestamp %v", err)
	}

	_, err = bcs.Unmarshal(response.SccPrices, &data.Prices)
	if err != nil {
		return nil, fmt.Errorf("bcs.Unmarshal  SccPrices %v", err)
	}

	return
}
