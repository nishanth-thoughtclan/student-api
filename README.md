# student-api
A student-api service is an app built using Go demonstrating the auth and the CRUD api's for student entity

### Installed swagger for api-docs using the below commands

go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

### Swagger Docs Url
http://localhost:8080/swagger/index.html

### Postman Collection Dump
provided the postman api dump with the examples [Link to file](./student-api.postman_collection.json)

### Provided MySql Schema & Dump for Testing
[Link to file](./mysql/create_schema.sql)
[Link to file](./mysql/student_dump.sql)

### Run the app using Docker
Hit 'docker-compose up --build' to run the app  