# Blue Jay - Blueprint

[![Go Report Card](https://goreportcard.com/badge/github.com/blue-jay/blueprint)](https://goreportcard.com/report/github.com/blue-jay/blueprint)
[![GoDoc](https://godoc.org/github.com/blue-jay/blueprint?status.svg)](https://godoc.org/github.com/blue-jay/blueprint)

Blueprint for your next web application in Go.

###Differences to upstream project
This is a fork with the following changes: 

1. Glide added for dependency management
1. Ported from mysql to postgres using dbr & github.com/jackc/pgx 
1. Email verification (using gomail) of user account before login
1. gohealth added but not yet intergrated. 
1. migrations using jay are now broken as jay does not support postgres. github.com/mattes/migrate added instead. 


## Quick Start Website

(installation steps are a bit different to the upstream project)

1. To download Blueprint, clone this repo into $GOPATH/src/github.com/blue-jay/blueprint
1. To download Jay, run the following command: go get github.com/blue-jay/jay
1. Install glide: https://github.com/Masterminds/glide
1. In your terminal, CD to the blueprint folder.
1. Run `glide install`.
1. Run `jay env make` to create the env.json file from env.json.example.
1. Start PostgreSQL, create a database and run migrations from migration/postgres the way you like.
1. Edit the PostgreSQL section of env.json to match your database login information.
1. Run the application: go run blueprint.go
1. Open your web browser to http://localhost and you should see the welcome page.
1. Navigate to the register page at http://localhost/register and create a new user.
1. You can now login at http://localhost/login.

## Running Tests

Before running tests ensure that your postgres user can create and drop other databases.
You can get all the correct settings with a single command, which is also convenient for development:

    docker run -p 5432:5432 -e POSTGRES_DB=blueprint -e POSTGRES_USER=blueprint postgres:9.6

Tests themselves are run with this command:

    go test $(glide novendor)

Two batch commands are provided, one for windows (tests.bat) and one for linux (tests.sh), which run all tests and exclude the vendor folder. 

## Original Information

Documentation available here: https://blue-jay.github.io/

Blue Jay is a web toolkit for [Go](https://golang.org/). It's a collection of
command-line tools and a web blueprint that allows you to easily structure
your web application. There is no rigid framework to which you have to
conform.

There are a few components:

- [**Blueprint**](https://github.com/blue-jay/blueprint) is a
model-view-controller (MVC) style web skeleton.

- [**Jay**](https://github.com/blue-jay/jay) is a command-line tool with
modules for find/replace, database migrations, code generation, and env.json.

- [**Core**](https://github.com/blue-jay/core) is a collection of packages
available to Blueprint and Jay as well as other applications.

Check the [milestones](https://github.com/blue-jay/blueprint/milestones) for
project status.