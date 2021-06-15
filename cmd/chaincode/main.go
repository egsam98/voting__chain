package main

import (
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/egsam98/voting/chain/services/smart"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := run(); err != nil {
		log.Fatal().Stack().Err(err).Msg("main: Fatal error")
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("main: Read ENVs from .env file")
	}

	cc, err := contractapi.NewChaincode(&smart.Contract{})
	if err != nil {
		return errors.Wrap(err, "failed to init chaincode")
	}

	err = cc.Start()
	return errors.Wrap(err, "failed to start chaincode")
}
