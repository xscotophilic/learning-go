# Ban7er API Documentation

## Overview

This API provides a robust solution for integrating live chat functionality into your application, supporting both one-to-one and group chats. It also includes encryption endpoints, allowing clients to generate their own encryption keys and manage encryption/decryption on their end using RSA (public and private keys in PEM format).

## Features

- **Live Chat**: Supports one-to-one and group chats via WebSocket.
- **Encryption**: Clients can generate and manage their own RSA keys for secure communication.
- **Secure Authentication**: Only users with valid `X-API-Key` and `X-API-Secret` headers can access the API.
- **No Database Required**: Uses Hashicorp Vault to store user keys, eliminating the need for additional database tables.
- **Channel-Based Communication**: Users can communicate in specific channels by passing a channel_id.

## Authentication

- All requests must include the following headers:
  - `X-API-Key`: API Key provided for authentication.
  - `X-API-Secret`: API Secret to validate the request.

## API Endpoints

### Key Management Routes

These endpoints help users generate and retrieve encryption keys.

#### Generate Encryption Keys

- **Endpoint:** `POST /api/v1/keys/generate`
- **Description:** Generates a new RSA key pair and stores it in HashiCorp Vault.
- **Query Parameters:**
  - `user_id`: ID of the user for whom the keys are being generated.
- **Response:**
  - `Status code`: 201

#### Get Public Encryption Key

- **Endpoint:** `GET /api/v1/keys/public`
- **Description:** Retrieves the public key of a user.
- **Query Parameters:**
  - `user_id`: ID of the user whose public key is needed.
- **Response:**
  ```json
  {
    "public_key": "-----BEGIN PUBLIC KEY-----..."
  }
  ```

#### Get Private Encryption Key

- **Endpoint:** `GET /api/v1/keys/private`
- **Description:** Retrieves the private key of a user.
- **Query Parameters:**
  - `user_id`: ID of the user whose private key is needed.
- **Response:**
  ```json
  {
    "private_key": "-----BEGIN PRIVATE KEY-----..."
  }
  ```

### Chat Routes

#### WebSocket Connection

- **Endpoint:** `GET /api/v1/ws`
- **Description:** Establishes a WebSocket connection for real-time messaging.
- **Query Parameters:**
  - `channel_id`: Unique identifier for the chat session.
- **Usage:**
  - One-to-One Chat: `channel_id` should be a combination of both user IDs (e.g., `user_a-user_b`).
  - Group Chat: `channel_id` should be a unique group identifier (e.g., group name or UUID).

## Integration Guide

### How It Works

- **Key Generation**:
  - When a user signs up, call the `/api/v1/keys/generate` endpoint to generate their RSA key pair.
  - Store the keys securely on the client side.
- **Encryption/Decryption**:
  - Use the public key to encrypt messages before sending them to the server.
  - Use the private key to decrypt messages received from the server.
- **Channel-Based Communication**:
  - For one-to-one chats, use a unique `channel_id` like `user_a-user_b`.
  - For group chats, use a `unique group identifier` or `group name` as the `channel_id`.
- **WebSocket Communication**:
  - Establish a WebSocket connection using the `/api/v1/ws` endpoint.
  - Send and receive messages in real-time within the specified channel.

### Example Workflow

- `User A signs up`:
  - Call /api/v1/keys/generate with user_id: user_a to generate keys.
  - Store the keys securely.
- `User B signs up`:
  - Call /api/v1/keys/generate with user_id: user_b to generate keys.
  - Store the keys securely.
- `User A wants to chat with User B`:
  - Establish a WebSocket connection to /api/v1/ws with channel_id: user_a-user_b.
  - Encrypt messages using User B's public key before sending.
- `User B receives messages`:
  - Decrypt messages using their private key.

### Error Handling

- **401 Unauthorized**: Invalid or missing X-API-Key or X-API-Secret.
- **400 Bad Request**: Missing or invalid parameters (e.g., user_id, channel_id).
- **500 Internal Server Error**: Server-side issues

### Best Practices

- Always encrypt sensitive data before sending it to the server.
- Use unique and consistent `channel_id` values for communication.
- Securely store keys on the client side.
- Rotate API keys and secrets periodically for enhanced security.

### Request examples

#### Generate Encryption Keys

```sh
curl -X POST "http://localhost:4000/api/v1/keys/generate?user_id=user_a" \
     -H "X-API-Key: test" \
     -H "X-API-Secret: test"
```

#### Retrieve Public Encryption Key

```sh
curl -X GET "http://localhost:4000/api/v1/keys/public?user_id=user_a" \
     -H "X-API-Key: test" \
     -H "X-API-Secret: test"
```

#### Retrieve Private Encryption Key

```sh
curl -X GET "http://localhost:4000/api/v1/keys/private?user_id=user_a" \
     -H "X-API-Key: test" \
     -H "X-API-Secret: test"
```
