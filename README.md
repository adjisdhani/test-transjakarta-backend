# TEST BACKEND ENGINEER TRANSJAKARTA Sistem Manajemen Armada

## Author: Adjis Ramadhani Utomo

## Project Setup Guide

This is a guide to install this code.

### 1. Clone the Repository

```sh
git clone https://github.com/adjisdhani/test-transjakarta-backend.git
```

### 2. Navigate into the Project Directory

```sh
cd test-transjakarta-backend
```

### 2. Requirement

```sh
install docker first
```

### 3. Install Dependencies

```sh
docker compose up -d --build
```

## API Testing using Postman

To test the API, use Postman or any other API testing tool.

# API Spec

## Last Location Vehicle API

Endpoint : GET vehicles/:vehicle_id/location

Response Body Success :

```json
{
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": 1715003456
}
```

Response Body Error :

```json
{
  "message": "vehicle not found"
}
```

## History Location Vehicle API

Endpoint : GET vehicles/:vehicle_id/history?start=11111&end=22222

Response Body Success :

```json
[
  {
    "vehicle_id": "B1234XYZ",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "timestamp": 1715003456
  }
]
```

Response Body Error :

```json
[]
```
