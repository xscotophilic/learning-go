package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
)

const secretsStoragePath string = "userdata/users/%s/secrets"

func (app *application) GenerateEncryptionKeysHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "user_id is required")
		return
	}

	privateKey, publicKey, err := app.generateEncryptionKeyPair()
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error generating key pair: %s", err))
		return
	}

	err = app.uploadEncryptionKeyPair(userID, privateKey, publicKey)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error storing keys: %s", err))
		return
	}

	err = app.writeStatusCode(w, http.StatusCreated)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetPublicEncryptionKeyHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	app.retrieveEncryptionKey(w, r, "public_key")
}

func (app *application) GetPrivateEncryptionKeyHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	app.retrieveEncryptionKey(w, r, "private_key")
}

func (app *application) generateEncryptionKeyPair() (
	string, string, error,
) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key pair: %v", err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	privateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}))
	publicKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}))

	return privateKeyPEM, publicKeyPEM, nil
}

func (app *application) uploadEncryptionKeyPair(
	userID, privateKey, publicKey string,
) (err error) {
	secretPath := fmt.Sprintf(secretsStoragePath, userID)
	secretData := map[string]interface{}{
		"data": map[string]interface{}{
			"private_key": privateKey,
			"public_key":  publicKey,
		},
	}

	_, err = app.vaultClient.Logical().Write(secretPath, secretData)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) retrieveEncryptionKey(
	w http.ResponseWriter,
	r *http.Request,
	keyType string,
) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "user_id is required")
		return
	}
	if keyType != "public_key" && keyType != "private_key" {
		app.serverErrorResponse(w, r, errors.New("invalid key type"))
		return
	}

	secretsPath := fmt.Sprintf(secretsStoragePath, userID)
	secret, err := app.vaultClient.Logical().Read(secretsPath)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get the keys: %s", err))
		return
	}

	if secret == nil || secret.Data["data"] == nil {
		app.errorResponse(w, r, http.StatusNotFound, "no keys found")
		return
	}

	data := secret.Data["data"].(map[string]interface{})
	key, exists := data[keyType].(string)
	if !exists {
		app.errorResponse(w, r, http.StatusNotFound, "no keys data found")
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{keyType: key}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
