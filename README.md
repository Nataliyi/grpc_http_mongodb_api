# grpc-api

This is a gRPC app build with Go, [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
and [ProtoBuff](https://developers.google.com/protocol-buffers) using HTTP as transport layer. It was build with the
purpose of having both gRPC and RESTful style integrated.

grpc-with-rest performs a [transcoding of HTTP calls to gRPC](https://cloud.google.com/endpoints/docs/grpc/transcoding)
using a proxy server.

The server and client side is all defined in _./proto/api.proto_ 

The following environment variables set in docker-compose.yml are used to configure the server:

- GRPC_HOST: host to connect to the server
- GRPC_PORT: port for gRPC requests
- GRPC_GATEWAY_PORT: port for HTTP requests

Description:

`grpcurl -plaintext localhost:8090 describe`


## How to Run it?

**In case you have docker-compose installed in your machine, just execute the following:**

~~~~
make up 
~~~~

Docker-compose will build all the dependencies and will add a MongoDB image in your container alongside
with the server so that we can interact with data.

*Once the docker-compose is finished, you should see an output in terminal:*

~~~
grpc-api       | 2022/10/25 14:09:57 Serving gRPC on localhost:8090
grpc-api       | 2022/10/25 14:09:57 Serving gRPC gateway on localhost:8080
~~~

*Send a POST request using cURL:*

`curl -X POST -k http://localhost:8080/api/v1/users/add -d '{"first_name": "FACEIT"}'`

You should have a response from server (example): 
~~~~
{"id":"bf7006db-4543-45c4-9be8-3b337b423524","status":200}
~~~~


## Available HTTP methods http://localhost:8080



| HTTP call |     Endpoints                                 |    Description                            | Required fields |
|:---------:|:---------------------------------------------:|:-----------------------------------------:|:---------------:|
|   POST    |     http://localhost:8080/api/v1/users/add    | Will create a user in database            | None            |
|   PUT     |     http://localhost:8080/api/v1/users/modify | Will update a user by id in database      | ID              |
|   DELETE  |     http://localhost:8080/api/v1/users/delete | Will delete a user by id in database      | ID              |
|   GET     |     http://localhost:8080/api/v1/get-all      | Will find all users in database           | None            |
|   POST    |     http://localhost:8080/api/v1/users/get    | Will find users by the filter in database | UsersFilter{Any}|


*Send a request using grpcurl:*

`echo '{"first_name": "FACEIT"}' | grpcurl -plaintext -d @ localhost:8090 UsersStore/AddUser`

You should have a response from server (example): 
~~~~
{
  "id": "ce555b8b-eeb3-47ac-9949-126e1d6b46a9",
  "status": 200
}
~~~~


## Available gRPC methods http://localhost:8090



| Method                      |                    Description                         |   Required fields     |
|:---------------------------:|:------------------------------------------------------:|:----------------------|
|    UsersStore/AddUser       |     Will create a user in database                     |      None             |
|    UsersStore/ModifyUser    |     Will update a user by id in databas                |      ID               |
|    UsersStore/DeleteUser    |     Will delete a user by id in database               |      ID               |
|    UsersStore/GetAllUsers   |     Will find all users in database                    |      None             |
|    UsersStore/GetUsers      |     Will find users by the filter in database          |      UsersFilter{Any} |


## UsersFilter

This filter helps to select users by certain fields. Possible selections include any combination
fields. For example, to select only those users whose first name is "Ally" and the last name is "Smit" or "Black" you need to write the following query:

`echo '{"first_name": ["Ally"], "last_name": ["Smit", "Black"]}' | grpcurl -plaintext -d @ localhost:8090 UsersStore/AddUser`

Filter by ID is also possible:

`echo '{"id": ["ad076657-bd10-4d66-97c5-7f228b521ae8", "53a14348-0cdc-485c-92c8-458018fe147c"]}' \
  | grpcurl -plaintext -d @ localhost:8090 UsersStore/AddUser`

Filter example using curl

`curl -X POST http://localhost:8080/api/v1/users/get \
    --data-binary 'id=ad076657-bd10-4d66-97c5-7f228b521ae8' \
    --data-binary 'id=53a14348-0cdc-485c-92c8-458018fe147c'`

## Logs

At the moment, a custom log collector is configured, which collects logs from actions in the api. The logs are saved to the current directory in the logs folder. You can change the settings using environment variables:

- LOGS_FREQUENCY_CREATING: frequency of creating new log files
- LOGS_PREFIX: prefix log file name
- LOGS_PATH: path to logs


## Data Base

The MongoDB database is used to store users. Environment variables are used to configure the database:

- DB_ADDR: host
- DB_PORT: port
- DB_DATABASE: db name
- DB_TABLE: collection name
- DB_LOGIN: login
- DB_PASSWORD: password


## Watcher

The service has a watcher stub. As one of the options, you can set up sending notifications to Google Pub Sub.
At the moment, the stub prints messages to the standard output of the service logs about the actions that led to changes in the database.

## HealthCheck

Available. The system uses https://github.com/grpc-ecosystem/grpc-health-probe


## How to Stop it?

~~~~
make up 
~~~~


## How to Test it?

~~~~
make test 
~~~~


## Request examples HTTP and gRPC

curl -X POST -k http://localhost:8080/api/v1/users/add -d '{"first_name": "FACEIT"}'

curl -X PUT -k http://localhost:8080/api/v1/users/modify -d '{"id": "bf7006db-4543-45c4-9be8-3b337b423524", "nickname": "face"}'

curl -X DELETE -k http://localhost:8080/api/v1/users/delete -d '{"id": "bf7006db-4543-45c4-9be8-3b337b423524"}'

curl -X GET -k http://localhost:8080/api/v1/users/get-all 

curl -X POST http://localhost:8080/api/v1/users/get \
    --data-binary 'id=ad076657-bd10-4d66-97c5-7f228b521ae8' \
    --data-binary 'first_name=Ally' --data-binary 'first_name=Andy'


echo '{"first_name": "FACEIT"}' \
    | grpcurl -plaintext -d @ localhost:8090 UsersStore/AddUser

echo '{"id": "53a14348-0cdc-485c-92c8-458018fe147c", "email": "ex@vv.com"}' \
    | grpcurl -plaintext -d @ localhost:8090 UsersStore/ModifyUser

echo '{"id": "0bdd3ac8-d745-4d26-929c-d18f1f5bf828"}' \
    | grpcurl -plaintext -d @ localhost:8090 UsersStore/DeleteUser

grpcurl -plaintext localhost:8090 UsersStore.GetAllUsers

echo '{"id": ["ad076657-bd10-4d66-97c5-7f228b521ae8"], "first_name": ["Ally","Andy"]}' \
    | grpcurl -plaintext -d @ localhost:8090 UsersStore.GetUsers


## Possible improvements

For faster interaction with the API, you can implement GRPC stream methods. This will increase the number
requests and responses per second with the correct implementation of the Golang code (using goroutines).

Example proto rpc for AddUser:

~~~~
rpc AddUser (stream User) returns (stream UserResponse)
~~~~

Example for GetAllUsers:

~~~~
rpc GetAllUsers (google.protobuf.Empty) returns (stream UsersList)
~~~~

It is also recommended to use Kubernetes to implement scaling and good control over replicas.

## An explanation of the choices taken and assumptions made during development

- gRPC is one of the most efficient means of data transfer with strong typing. The gRPC messaging system uses the binary Protobuf format instead of JSON, resulting in smaller messages and increased throughput - the transfer rate eventually increases by 7-10 times. This is indispensable in highly loaded systems.

- Prometheus Metrics allows you to export metrics to any monitoring service convenient for you.

- Pub Sub messaging system is used to effectively inform interested services. The ability to integrate with the Google Cloud Platform is especially convenient.

- Buf is used to replace the current REST/JSON based API development paradigm with a schema based one. Defining an API with an IDL has many advantages over REST/JSON, and Protobuf is by far the most stable and widely adopted IDL in the industry.

- MongoDB is a document database used to build highly available and scalable internet applications. With its flexible schema approach, it's popular with development teams using agile methodologies. A plus for me as a Golang developer is the opportunity to learn something new, since this is the first experience of developing services using MongoDB and it did not disappoint.


*During the development of the service, due to the limited amount of time, there were some difficulties:*

1. Format of MongoDB with datetime. I had to use string date format. because it was not possible to understand type conversion. In this regard, at the moment the service does not have the ability to use a filter by date.

2. Also, according to the task, the approximate number of requests to the API per second was not indicated. Therefore, methods based on stream were not created. But in my development experience, I can say that stream allows you to significantly increase the number of requests.

3. While testing the filters, it turned out that the standard curl -d does not convert the request to Protobuf. Therefore, to query by filters, I had to use the --data-binary parameter. But I believe that there is a more convenient way to transfer and is ready to accept your advice.

4. While testing the user ID, it became clear that the Object ID in MongoDB does not support the UUID format required by the task. Thanks to a small parser from github implemented in the servises/mongo-registry.go file, it became possible to store and retrieve the _id field in a string format that stores the UUID ID.

5. One of the development rules is that there are never too many tests. I understand that the implemented tests are absolutely not enough to fully verify the operation of the service. The minimum number was written only to check some points.

6. Docker-compose is a good tool for deploying a test job. But for the release of the product, I assume that it is not enough. I think that Kubernetes is more efficient and interesting. Especially in conjunction with GCP.


## Bonus

Additionally, metrics are configured in the API:

- *custom_api_heath_check* - when checking health-check sends data to the metric
- *custom_api_errors* - sends the number of errors received by the API

Environment variables:

- METRICS_SERVER_PORT: port where metrics are presented
- METRICS_SERVER_RUNTIME: metrics server running time after service shutdown
- METRICS_PATH: path to host metrics

You can view the metrics at http://localhost:9090/metrics. It also provides standard Golang metrics.