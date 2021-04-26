# Lana Flagship Store

Project to resolve Lana store challenge. Made with Golang :clap :clap (my first code in Golang actually)


## Installation

*Golang disclaimer*: Project should be at _$GOPATH/src/lana/_

To install the project run:

    make install

To build the project run:

    make build

To run the container run:

    make run

then the API will be ready at `http://localhost:3080/`


## Project folders

    ./flagship-store
    |-- models
    |-- persistence
    |-- services
    |   |-- commands
    |   |-- errors
    |   └-- responses
    |-- utils
        └-- mocks

_models_: Domain objects classes.

_persistence_: Repository classes and interfaces to deal with our persistence system(local array, database or whatever)

_services_: Services for each action that can be executed.

_services/commands_: Objects used as a parameters of services.

_services/errors_: Services response errors.

_services/responses_: Application response based on service response. Used at controller layer.

_utils/mocks_: Services mocks used for testing.

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
