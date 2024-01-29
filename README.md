# Parallel Specs

Simple api to grab spec files in order: new specs -> longest specs -> short specs and calculate estimated time for each instance to reduce total execution time  

# Use

- [Web UI](https://parallel-specs.shelex.dev) - for demo purposes you can use email: `test@test.com` password: `test`
- [Swagger](https://parallel-specs.shelex.dev/swagger) - documentation and api playground
- [JS client](https://github.com/Shelex/parallel-specs-client) - api wrapper and CLI for creating session or running cypress

# Flow

```mermaid
    graph LR;
    project-->spec_files;
    project-->session;
    session-->spec_executions;
    spec_files-->spec_executions;
```

- Register user or login with existing one (for js client library you can create api key).
- Create new session (it will be attached to existing project or will create new)
- Get nextSpec for your sessionID and machineID, every query will finish previous spec for this session + machine and return next. Final query will return message "session finished" and finish spec and session for specific machineID. in case machineID is not passed it will be "default"

# Try it locally
- clone this repository
- `cd parallel-specs`
- `make deps` - download dependencies
- `make keys` - generate private and public keys for jwt auth
- prepare postgres db, example: `docker run -p 5432:5432/tcp --name postgres-for-parallel-specs -e POSTGRES_PASSWORD=verysecret -e POSTGRES_DB="parallel-specs" -d postgres`
- rename `config.env.sample` -> `config.env`
- set `APP_DB_CONNECTION_URL` in `config.env` with your db instance connection string
- `make migration` - run migrations for postgres, first execution will setup db schema
- `make start` - start the service
- `make web-dev` - start ui
- open `http://localhost:3000/swagger` for Swagger documentation
- use `http://localhost:3000/api` for Altair/Postman/Insomnia api clients
- use `http://localhost:3000/` for ui interface
