# TODO:
- use http/net instead: DONE
- construct project structure: modular monolith, utilize go workspace

### Stack
routing: http/net
db: bun, postgressql, goose
blockchain: geth
logging & mornitoring: zerolog(golang)
    stage 0 : gcp, aws native logging service
    stage 1 : promtail + logi + grafana (self-hosted)

## project structure : with go workspace
module:server
    /cmd
        /api
            api.go
        /database
            seed.go
            migrate.go
    /server
        server.go // construct server, 
    /internal
        /route
        /handler
        /service
        /repo

module:auth
    main.go // expose auth func
    /internal
        repo.go
        service.go
        service_test.go

module:stdx
    /httpx
    /errorx
    /ciphers
        /jwt.go
        /hash.go

module:mailing
        
module:dtacc
    dtacc.go   // package dataac
    init.go
    migrate.go
    testing.go // provide test container, functionality
    Makefile   // 
    /models
        user.go, session.go, org.go
    /internal
        /migrations
            migration_0001.go
            migration_0002.go

# for integration/smoke test
https://testcontainers.com/
TestContainer-go

# mailing
- maillite
- mailgun

# golang ?
- panic, recover
- graceful shutdown

# techinal requirement
https://www.youtube.com/watch?v=MhfH1H6fAIM&t=1600s

- fail over management
    - curcuitbreaker
    - retry, retry limit
