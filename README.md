# billingApp

Using tools: Golang, PostgreSQL, RabbitMQ, Docker 

Microservice can increase a user's balance or decrease, send money to another user, freeze the user's balance and unfreeze. Accept requests in json ([]byte).  

## Run the app

Running building

    make build
    
Running app

    make run
    
After first run you have to exit. (CTRL+C) To make DB

    make createdb
    
Init DB

    make migrate
    
DB start with 4 rows having amount = 0. 

Restart app

    make run
    
Microservice take requests to Queue

    "rpc_queue"

To start app on localhost

    go run ./cmd/worker.go
    
To send requests on localhost

    go run ./cmd/sender.go
    
In sender.go you can write request body to var data.

## Increase balance

### Request  

    {"method":"increase","balance_id":1,"amount":4,"receiver":0,"freeze_id":0,"freezed_amount":0,"is_approved":false} 
    
### Response

    {"data":{"balance_id":1,"amount":4,"msg":"balance-changed"}} 

## Decrease balance

### Request

    {"method":"decrease","balance_id":1,"amount":4,"receiver":0,"freeze_id":0,"freezed_amount":0,"is_approved":false} 
    
### Response

    {"data":{"balance_id":1,"amount":0,"msg":"balance-changed"}} 

## Send to user

### Request

    {"method":"send","balance_id":1,"amount":4,"receiver":2,"freeze_id":0,"freezed_amount":0,"is_approved":false}

### Response

    {"data":{"sender_id":1,"sender_balance":0,"receiver_id":2,"receiver_balance":4,"msg":"money-transfered"}} 

## Freeze balance

### Request

    {"method":"freeze","balance_id":2,"amount":0,"receiver":0,"freeze_id":0,"freezed_amount":1,"is_approved":false} 

### Response

    {"data":{"freeze_id":1,"freezed_amount":1,"msg":"balance-freezed"}} 

## Approve freezed amount

### Request

    {"method":"approve","balance_id":0,"amount":0,"receiver":0,"freeze_id":1,"freezed_amount":0,"is_approved":false} 

### Response
    
    {"data":{"balance_id":2,"amount":4,"msg":"balance-unfreezed"}} 
