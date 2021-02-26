- BusinessProcesses

- Logging
    - improve logging with Fields and Context + 
    - setup logging configuration for each service
    - logging to external system (sentry, loki etc.)
    - mdw for tracing all gRPC/HTTP requests

- Translations
    - study approaches

- Users & Sessions
    - think of moving sessions from Api to a separate service
    - reconnect for chat ws connection  
    - check expiration of a token (depends on what identity provider is going to be used)
    - persistence - do we really need it?
    - cluster mode - broadcast a new session to all available nodes
    - roles and groups

- Socket hub 
    - gRPC API to send socket messages from app (done through NATS) + 
    - move to a separate service (now it's in Api)
    - gracefully close socket connections (currently we got an error)
    - keep alive - close WS server after configured keep alive period of not getting ping message +

- Configuration (config + .env) 
    - make configuration available before initialization

- panic recover
    - decide on which level recover makes sense
    - suggest a general approach

- Elastic Search
    - use in users service (and other services)
    - mapping shouldn't be hardcoded in src as it's implemented now (maybe to store mapping as files)

- Cache
    - suggest a common approach for caching (the current idea is to cache on the storage level and leave it up to storage implementation) 

- Error handling
    - the main question: do we need a custom type for errors? the main issue is to recognize errors type (code isn't available)
    
- API versioning
    - protobuf message versioning
    - http api versioning

- Task service
    - task config persistence
    - auto-assign and scheduler in cluster mode (we don't want to have multiple instances of some processes)

- Monitoring
    - study a way how to prepare metrics: http, gRpc, database, cache etc
    - monitoring of tools (zeebe, Postgre, Nats, etc)

- Validation 
    - study go-playground/validator
    - implement basic validations for business services
    
- Integrations with external systems
    - suggest approach
    
- Storage
    - Optimistic locking
    - Master-Slave 
    
- Cluster mode 
    - support simple leader election based on Redis
    - gRpc LB
    - coordinating with etcd
    
- Mattermost
    - reconnect to NATS after NATS goes down +
    - cluster mode - deep study 
    
- Prometheus
    - metrics support
    
- Processes
    - handle lost calls from zeebe which happens just after app starts

- Refactoring
    - queue topics must be declared as consts + 
    
- Webrtc
    - ION implementation

- Deployment
    - microservices split
    - prepare Makefile/Dockerfiles
    - prepare Kuber configs
    - run app under Skaffold
    - initial configuration (db, chat etc.)
    - put mattermost to Docker from source (there is an issue with Frontend)
     
- Documentation
    - code guide
    - architecture (main points)