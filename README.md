#### Development

- rename `./.env-default` to `./.env`
- setup these variables in your environment by either explicitly setting in IDE or by running the shell command `source .env`
- check and configure `./config.yml` properly (make sure your have all infrastructure components up and running)
- run ./cmd/main

##### Build

``make build``

##### Docker

- create a docker network `mm` or make sure it has been created before
- rename `./.env-default` to `./.env` for all the services
    - ./api
    - ./bp
    - ./chat
    - ./services
    - ./sessions
    - ./tasks
    - ./users
    - ./webrtc
- check and configure `./config/config.yml` properly (make sure your have all infrastructure components up and running)
- run ``docker-compose up --build --remove-orphans``