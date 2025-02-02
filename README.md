# ICE Global

![go-version](https://img.shields.io/badge/Golang-1.23-66ADD8?style=for-the-badge&logo=go)
![app-version](https://img.shields.io/github/v/tag/mohammadne/ice-global?sort=semver&style=for-the-badge&logo=github)
![coverage](https://img.shields.io/codecov/c/github/mohammadne/ice-global?logo=codecov&style=for-the-badge)
![repo-size](https://img.shields.io/github/repo-size/mohammadne/ice-global?logo=github&style=for-the-badge)

## Introduction

The `ICE Global` is a simple shopping cart manager for handling user cart iems.

This Project was part of an [interview refactor](https://git.ice.global/packages/golang-interview-refactor) for the ICE Global.

![Shopping cart manager](assets/application.png)

### Architecture

```tree
> tree
.
├── cmd
│   ├── migration
│   │   └── schemas
│   ├── version.go
│   └── web-api
├── deployments
│   ├── docker
│   └── k8s
├── internal
│   ├── api
│   │   └── http
│   ├── config
│   ├── entities
│   ├── services
│   └── storage
└── pkg
    ├── mysql
    └── redis
```

#### CI

I have used `Github Action` to build the image which uploads it into `ghcr.io` registry.

#### Deployments

I have deployed the application via `docker-compose`, the docker directory contains 2 compose file which we can use `.local` for resolving depencies for local development and the `.prod` for deploying the whole application into a server via docker-compose file.

Also I have developed a `k8s` directory which can deploy the apploication via `helmsman` and `helm` charts. for deploying the ice-global I have setup a kubernetes cluster via `Kind`.

#### API

for the api layer we at this stage we only have `http` handler and I have used `embed` for embedding template files,the pervious design used an absolute-path which causes problems in the deployment.

As pervious I have used the `Gin` engine and I don't see any benefit for changing the http engine, but in high trouphot systems like `Snapp` which I have exprience in it, they using frameworks based on `fast http` like `fiber` but here the Gin was enough.

Also I have developed a `middleware` to decouple `cookie` checking for users and extract theses codes into some middlewares which makes our code much more cleaner and maintanable.

**NOTE:** There was some confiusions about term `session` which is was used incorrectly instead of cookie, in code I have changed this term but for the database I put the term `session_id` unchanged.

#### Config

For reading configurations I have used a very simple package named `envconfig` to avoid any hard coding configuration into the code. for much more advanded usecases we can use other tools like `koanf` but for our case I see the envconfig enough.

#### Services

The services diectory acts as `usecases` in the clean-architecture which contains all the business-logic stuffs and the outer layers like http(gin) is unaware of any business-logic of the application. the http layer only works with queries, arguments, user input validation and ... .

#### Storage

As pervious I have used mysql but with some changes:

1. I have changed the `gorm`
2. develop migrations logic

In my attitude if we can use raw sql queries, it give us lot's of benefits:

1. more readable and much more cleaner
2. more performance
3. having more optimized queries like you want different lists of a table with some atributes, if you have raw sql queries you can optimize your sql queries based on given attributes.
4. readable by everyone specially by DBAs.

#### Cache

I have used Redis cache for caching items and prohibit access to items of the primary database each time.

As you know the Redis is an in an in-memory storage which read-writes are much more faster rather than databases which stores the data on the storage like mysql in this case.

#### Packages

In my codes I like to put my general package creation like mysql and redis (in this case) into the `/pkg` directory and putting the business logic only related to this project into `/internal` directoy which is something recommanded by Go communities.

### Data Migration

```sh
go run cmd/migration/main.go --direction up
```

The command above will migrate the current data schemas to a newer and better data schemas.

At first we have the `01_initial.up.sql` sql file, as you can see there is no any index and ..., at the given steps we migrate the existing data, add `items` table and ... to migrate to the new design. you can see the exact migration steps by reading migration files.

#### Why migration was necessary?

1. In the `cart_items` the cart_id has no relation cfonstraint to `cart_entities` table.
2. Add `items` table to make `cart_items` more concise and cleaner and remove hard codes in the `pkg/calculator/add_item.go`.
3. The columnd `session_id` should be unique accross table.
4. The `total` column in `cart_entities` makes the table un-normalize because this field can be calculated from `cart_items` and the items associated with an `cart_entities`.
5. for tracking changes it's better to use `TIMESTAMP` rather than `DATETIME`.
