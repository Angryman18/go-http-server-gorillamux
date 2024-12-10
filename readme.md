### How to Run?


* Point a postgress db and create two table with the following command


##### Users Table
```
create table users (
    id uuid primary key not null,
    name varchar(100) not null,
    email varchar(100) unique not null,
    password varchar(255) not null
);
```
##### Todo Table
```
create table todo (
    id uuid primary key not null,
    user_id uuid not null,
    todo text not null,
    is_completed bool,
    created_at timestamp,
    updated_at timestamp,
    foreign key (user_id) references users(id)
);
```

* Get a postgress connection string and paste in the `config.env` file.

* `config.env` should look like this all these values.
```
PORT = 8000
JWT_SECRET = ""
POSTGRESS_CONN = ""
```

## Swagger Endpoint

* Swagger docs are available on `http://localhost:8000/swagger/index.html`

### Generate Swagger Docs For Apis (Optional)
* Need to install `swag`
* run `go install github.com/swaggo/swag/cmd/swag@latest` and this saves at `go/bin/swag`

* if you are using windows then need to add `go/bin` as path variable
* everytime you make changes in the api docs run `swag init` to generate docs

### Run Server
* `go get` and then `go run .` make sure you are hitting this command at root where you `main.go` file is located.

## Thats it. Ready 🚀🚀🚀🚀🚀




---

# Project Routes

```
HOME         = "/"
LOGIN        = "/login"
SIGNUP       = "/signup"
HEALTH       = "/health"
CREATE_TODO  = "/create_todo"
GET_ALL_TODO = "/get-all-todo"
UPDATE_TODO  = "/update-todo"
DELETE_TODO  = "/delete-todo"
```