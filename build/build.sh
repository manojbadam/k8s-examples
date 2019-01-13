#! /bin/bash

docker build -t sample-app app/.
docker tag sample-app manojbadam/liveness
docker push manojbadam/liveness