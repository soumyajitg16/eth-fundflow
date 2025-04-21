# eth-fundflow
This project aims to trace the flow of Ethereum funds for any given  address, helping to uncover both the final recipients (beneficiaries) of  outgoing funds and the original senders (payers) of incoming funds.

## Ethereum Fund Flow Analysis

# Provides two endpoints:

- **/beneficiary?address=<eth_address>**: final outflow recipients
- **/payer?address=<eth_address>**: original inflow senders

## Setup

cd eth-fundflow
1. Clone repo ("git clone https://github.com/soumyajitg16/eth-fundflow.git") & `cd eth-fundflow`
2. code. & `go mod download`
3. Add the API key: "ETHERSCAN_API_KEY=VEV3VC4CA3M2BS77TEIXC64PUHQAR77T7P"
   ** In the .env file in the root directory **
   ETHERSCAN_API_KEY=YourAPIKeyHere
4. `go run cmd/server/main.go`

## Example Addresses

**Beneficiary**:
- 0x3bac7400123bde8bbe813fbb781796d470164a1c
- 0x28c6c06298d514db089934071355e5743bf21d60

**Payer**:
- 0x3bac7400123bde8bbe813fbb781796d470164a1c
- 0x28c6c06298d514db089934071355e5743bf21d60

- # Test URL+Addresses-
# Beneficiary
// curl "http://localhost:8080/beneficiary?address=0x28c6c06298d514db089934071355e5743bf21d60"
// curl "http://localhost:8080/beneficiary?address=0x3bac7400123bde8bbe813fbb781796d470164a1c"
# Payer
// curl "http://localhost:8080/payer?address=0x28c6c06298d514db089934071355e5743bf21d60"
// curl "http://localhost:8080/payer?address=0x3bac7400123bde8bbe813fbb781796d470164a1c"
