# Lana Flagship Store

Project to resolve Lana store challenge.


## Installation

To install the project run:

    make install

To build the project run:

    make build

To run the container run:

    make env-start

then the API will be ready at `http://localhost:3080/`


## Testing

To execute test run:

    make test


## Operations

### Create a new checkout basket

To create a new checkout basket, in terminal execute:

    curl -w "%{http_code}" --location --request POST 'http://localhost:3080/checkouts' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "product-code": "PEN"
    }'

Possible responses:
- Success: Code 200 with body

            {"id":"eefc5ac5-8f90-4f87-91e2-1f425781d8fb","products":["PEN"]}

- Failed: Code 404 with body

            {"message":"Product FAKE not found"}

.

### Add a product to a basket

To add a product to a basket, in terminal execute:

    curl -w "%{http_code}" --location --request PATCH 'http://localhost:3080/checkouts/45120489-458f-4567-9d7a-c0d83b55128e' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "product": "TSHIRT"
    }'

Possible responses:
- Success: Code 204

- Failed:

  - Code 404 with body

            {"message":"Checkout a_fake_checkout not found"}

  - Code 422 with body

            {"message":"Product FAKE not found"}

.


### Get the total amount in a basket

To get the total amount in a basket, in terminal execute:

    curl -w "%{http_code}" --location --request GET 'http://localhost:3080/checkouts/45120489-458f-4567-9d7a-c0d83b55128e/amount'

Possible responses:
- Success: Code 200 with body

    {"amount":"27.50€"}

- Failed:

  - Code 404 with body

            {"message":"Checkout a_fake_checkout not found"}

.

### Remove the basket

To remove the basket, in terminal execute:

    curl -w "%{http_code}" --location --request DELETE 'http://localhost:3080/checkouts/45120489-458f-4567-9d7a-c0d83b55128e'

Possible responses:
- Success: Code 204

- Failed:

  - Code 404 with body

            {"message":"Checkout a_fake_checkout not found"}