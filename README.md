
# JustLend Energy

## Overview

This project implements energy rental functionality on the JustLend platform using Golang, providing energy cost calculation, rental, and refund features for a fast and convenient experience.

## Quick Start

You can [download the precompiled version here](https://github.com/19byte/justlend_energy/releases) or compile and run it yourself with the following steps:

```shell
git clone https://github.com/19byte/justlend_energy.git
cd justlend_energy
go run cmd/justlend.go
```

## HTTP Interface (Port: 8085)

For detailed definitions of the `type` field, please refer to the [tronprotocol/protocol GitHub repository](https://github.com/tronprotocol/protocol/blob/2a678934da3992b1a67f975769bbb2d31989451f/core/contract/common.proto#L9).

### Query Fee

- **Description**: Retrieve the current rental fee information
- **Method**: GET
- **Endpoint**: `/fee`

**Query Parameters**

| Name       | Method | Type   | Required | Remark                     |
|------------|--------|--------|----------|----------------------------|
| energy     | query  | string | Yes      | Amount of energy to rent   |
| privateKey | query  | string | Yes      | User's private key         |
| type       | query  | string | Yes      | Rental type                |

**Response Example**

```json
{
  "code": 1000,
  "data": {
    "rentAmount": 1000,
    "stakePerTrx": 82,
    "liquidateThreshold": "0",
    "rentalRate": "0.0000000084933752",
    "feeRatio": "40",
    "minFee": "40",
    "curFeeRatio": "0.041",
    "rentFee": "0.12034772923392",
    "prePayFee": 40.12034772923392
  }
}
```

---

### Rent Energy

- **Description**: Rent a specified amount of energy
- **Method**: POST
- **Endpoint**: `/rent`

**Body Parameters**

```json
{
  "receive": "receiver address",
  "type": 1,
  "amount": 100000,
  "privateKey": "user's private key"
}
```

**Response Example**

```json
{
  "code": 1000,
  "data": {
    "txId": "transaction ID",
    "stakePerTrx": 8174000000
  }
}
```

---

### Return Energy

- **Description**: Return previously rented energy
- **Method**: POST
- **Endpoint**: `/return`

**Body Parameters**

```json
{
  "receive": "receiver address",
  "type": 1,
  "stakePerTrx": 8174000000,
  "privateKey": "user's private key"
}
```

**Response Example**

```json
{
  "code": 1000,
  "data": {
    "txId": "transaction ID"
  }
}
```

---

## Contributing

We welcome contributions to improve this project! You can submit your code via Pull Requests or leave your feedback in the Issues section.
