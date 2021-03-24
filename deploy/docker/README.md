
````shell script
clone https://github.com/adacta-ru/mattermost-server.git

clone https://github.com/adacta-ru/mattermost-webapp.git

cd mattermost-webapp

make build

cd ../docker

./run.sh
````