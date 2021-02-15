- Logging (request id) 
    TODO:
    - improve logging with Fields
    - setup logging configuration for each service
    - logging to external system (sentry, loki etc.)
    - mdw for tracing all gRPC/HTTP requests

- Translations
    TODO:
    - study approaches

- Sessions
    TODO:
    - think of moving sessions from Api to a separate service
    - reconnect for chat ws connection  
    - check expiration of a token (depends on what identity provider is going to be used)
    - persistence - do we really need it?
    - cluster mode - broadcast a new session to all available nodes

- Socket hub 
    TODO:
    - gRPC API to send socket messages from app
    - move to a separate service (now it's in Api)
    - gracefully close socket connections (currently we got an error)
    - keep alive - close WS server after configured keep alive period of not getting ping message

- Configuration (config + .env) 
    TODO:
    - make configuration available before initialization

- panic recover
    TODO:
    - decide on which level recover makes sense
    - suggest a general approach

- Elastic Search
    TODO:
    - use in users service (and other services)
    - mapping shouldn't be hardcoded in src as it's implemented now (maybe to store mapping as files)

- Cache
    TODO:
    - suggest a common approach for caching (the current idea is to cache on the storage level and leave it up to storage implementation) 

- Error handling
    TODO:
    - the main question: do we need a custom type for errors? the main issue is to recognize errors type (code isn't available)
    
- API versioning
    TODO:
    - protobuf message versioning
    - http api versioning

- Task service
    TODO:
    - task config persistence
    - auto-assign and scheduler in cluster mode (we don't want to have multiple instances of some processes)

- Monitoring
    TODO:
    - study a way how to prepare metrics: http, gRpc, database, cache etc
    - monitoring of tools (zeebe, Postgre, Nats, etc)

- Validation 
    TODO:
    - study go-playground/validator
    - implement basic validations for business services
    
- Integrations with external systems
    TODO:
    - suggest approach
    
- Storage
    TODO:
    - Optimistic locking
    
- Cluster mode 
    TODO:
    - support simple leader election based on Redis
    
- Mattermost
    TODO:
    - reconnect to NATS after NATS goes down
    - cluster mode - deep study 

- Refactoring
    TODO:
    - queue topics must be declared as consts
    
- Deployment
    TODO:
    - microservices split
    - prepare Makefile/Dockerfiles
    - prepare Kuber configs
    - run app under Skaffold
    - initial configuration (db, chat etc.)
     
- Documentation
    TODO:
    - code guide
    - architecture (main points)