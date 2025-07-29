## Go Wallet API

go-wallet is a mini example Go API project that provides a simple CRUD.  
It allows users to register, log in, create and manage multiple wallets, and perform transactions such as deposit, withdraw, and transfer between wallets.  

## Getting Started Locally

1. Make sure you have [Docker](https://www.docker.com/) installed.

2. Clone this repository to your local machine.

3. In the project directory, run:
   ```sh
   make run-dev
   # or
   docker compose -f docker-compose.dev.yml up
   ```

4. The API will be available at [http://localhost:8080/healthz](http://localhost:8080/healthz).

## Run Unit Tests

To run all unit tests, use:

```sh
make test
# or
go test -v ./...
```

## API Spec

You can get the OpenAPI spec at `docs/server.yml` to import into your Postman or API client.

## Flow to Test the API

1. **Register**  
   Send a POST request to `/public/register` with your email, password, and displayName.

2. **Login**  
   Send a POST request to `/public/login` with your email and password.  
   Copy the `accessToken` from the response.

3. **Authenticate**  
   For all `/secure` endpoints, set the `Authorization` header:  
   ```
   Authorization: Bearer <accessToken>
   ```

4. **Wallet Operations**
   - **Create Wallet:**  
     POST `/secure/wallet` create new wallet to user
   - **List Wallets:**  
     GET `/secure/wallets` to see all your wallets.
   - **Update Wallet:**  
     PUT `/secure/wallet/{walletId}` to update wallet info.
   - **Delete Wallet:**  
     DELETE `/secure/wallet/{walletId}` to remove a wallet.

5. **Transaction Operations**
   - **Deposit:**  
     POST `/secure/deposit` to add initial points to your wallet.
   - **Transfer:**  
     POST `/secure/transfer` to move balance between wallets.
   - **Withdraw:**  
     POST `/secure/withdraw` to directly remove points from a wallet.
   - **List Transactions:**  
     GET `/secure/wallet/{walletId}/transactions` (supports `page` and `limit` query params).

