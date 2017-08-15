#!/usr/bin/env bash
set -e

curl localhost:3000/users/new -XPOST -d '{"id":10,"email":"ff@ff.ru","first_name":"Елизавета","last_name":"Хопленлова","gender":"f","birth_date":406080000}'
curl localhost:3000/locations/new -XPOST -d '{"id":100,"distance":9,"city":"омск","place":"Рест","country":"элла"}'
curl localhost:3000/locations/new -XPOST -d '{"id":101,"distance":9,"city":"о1","place":"Рест11","country":"элла1"}'
curl localhost:3000/visits/new -XPOST -d '{"id": 1000, "user": 10, "visited_at": 1385572234, "location": 100, "mark": 2}'
curl localhost:3000/users/10/visits
curl localhost:3000/users/new -XPOST -d '{"id":11,"email":"ff1@ff.ru","first_name":"Елизавета","last_name":"Хопленлова","gender":"f","birth_date":406080001}'
curl localhost:3000/visits/1000 -XPOST -d '{"user": 11}'
curl localhost:3000/users/10/visits
curl localhost:3000/users/11/visits
curl localhost:3000/locations/100/avg
curl localhost:3000/visits/new -XPOST -d '{"id": 1001, "user": 11, "visited_at": 1385572235, "location": 100, "mark": 4}'
curl localhost:3000/locations/100/avg
curl localhost:3000/visits/new -XPOST -d '{"id": 1002, "user": 11, "visited_at": 1385572236, "location": 100, "mark": 5}'
curl localhost:3000/locations/100/avg
curl localhost:3000/locations/101/avg
curl localhost:3000/visits/1000 -XPOST -d '{"location": 101}'
curl localhost:3000/locations/100/avg
curl localhost:3000/locations/101/avg
#curl localhost:3000/visits/888997 -XPOST -d '{"location": 1}'
#curl localhost:3000/locations/777777/avg
