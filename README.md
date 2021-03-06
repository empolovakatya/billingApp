# billingApp

Using tools: Golang, PostgreSQL, RabbitMQ, Docker 

Microservice can increase a user's balance or decrease, send money to another user, freeze the user's balance and unfreeze. Accepts requests in json ([]byte).  

## Run the app

Running building

    make build
    
Running app

    make run
    
After first run you have to exit to make DB (CTRL+C)

    make createdb
    
Init DB

    make migrate
    
DB starts with 4 rows in table `balances` having amount = 0. 

Restart app

    make run
    
Microservice accepts requests to Queue

    "rpc_queue"

To start app on localhost

    go run ./cmd/worker.go
    
To send requests on localhost

    go run ./cmd/sender.go
    
In sender.go you can write request body to var data.

    Example: data := billing.WorkerSender{Method: "approve", FreezeId: 1, IsApproved: false}

## Increase balance
using params: `method`, `balance_id`, `amount` 

### Request example

    {"method":"increase","balance_id":1,"amount":4,"receiver":0,"freeze_id":0,"freezed_amount":0,"is_approved":false} 
    
### Response example

    {"data":{"balance_id":1,"amount":4,"msg":"balance-changed"}} 

## Decrease balance
using params: `method`, `balance_id`, `amount`

### Request example 

    {"method":"decrease","balance_id":1,"amount":4,"receiver":0,"freeze_id":0,"freezed_amount":0,"is_approved":false} 
    
### Response example

    {"data":{"balance_id":1,"amount":0,"msg":"balance-changed"}} 

## Send to user
using params: `method`, `balance_id`, `amount`, `receiver`

### Request example 

    {"method":"send","balance_id":1,"amount":4,"receiver":2,"freeze_id":0,"freezed_amount":0,"is_approved":false}

### Response example

    {"data":{"sender_id":1,"sender_balance":0,"receiver_id":2,"receiver_balance":4,"msg":"money-transfered"}} 

## Freeze balance
using params: `method`, `balance_id`, `freezed_amount`

### Request example

    {"method":"freeze","balance_id":2,"amount":0,"receiver":0,"freeze_id":0,"freezed_amount":1,"is_approved":false} 

### Response example

    {"data":{"freeze_id":1,"freezed_amount":1,"msg":"balance-freezed"}} 

## Approve freezed amount
using params: `method`, `freeze_id`, `is_approved`

### Request example 

    {"method":"approve","balance_id":0,"amount":0,"receiver":0,"freeze_id":1,"freezed_amount":0,"is_approved":false} 

### Response example
    
    {"data":{"balance_id":2,"amount":4,"msg":"balance-unfreezed"}} 
