package hyperledger

import (
	"io/ioutil"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
)

const walletLabel = "app_user"

type Config struct {
	Channel, ConfigPath, MSPID, CertPath, PrivateKeyPath string
	ChaincodeID                                          string
}

func Connect(cfg Config) (*gateway.Network, error) {
	wallet, err := generateWallet(cfg.MSPID, cfg.CertPath, cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(cfg.ConfigPath)),
		gateway.WithIdentity(wallet, walletLabel),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to gateway")
	}

	network, err := gw.GetNetwork(cfg.Channel)
	return network, errors.Wrap(err, "failed to get network")
}

func generateWallet(mspID, certPath, privateKeyPath string) (*gateway.Wallet, error) {
	// read the certificate pem
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read cert path %s", certPath)
	}

	// read private key
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", privateKeyPath)
	}

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))

	wallet := gateway.NewInMemoryWallet()
	if err := wallet.Put(walletLabel, identity); err != nil {
		return nil, errors.Wrapf(err, "failed to put identity to wallet with label %s", walletLabel)
	}

	return wallet, nil
}
