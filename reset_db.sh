#!/bin/bash

docker exec -e DISABLE_DATABASE_ENVIRONMENT_CHECK=1 acceptance-web-1 bundle exec rails db:schema:load
docker exec acceptance-web-1 bundle exec rails runner "Ops::Bots::Create.new(:booking).call(username: 'test_beauty36_bot', token: '6750404821:AAEuOxy3Z9k09i5q4tceHoI8ZHTbarwJS1Y')"
