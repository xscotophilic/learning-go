package main

import (
	"crypto/tls"
	"net/http"

	vault "github.com/hashicorp/vault/api"
)

type VaultModel struct {
	address            string
	token              string
	insecureSkipVerify bool
}

func (vaultModel *VaultModel) NewClient() (*vault.Client, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = vaultModel.address

	if vaultModel.insecureSkipVerify {
		vaultConfig.HttpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	vaultClient.SetToken(vaultModel.token)

	return vaultClient, nil
}
