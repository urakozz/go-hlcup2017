#!/usr/bin/env bash
set -e

curl localhost:3000/users/new -XPOST -d '{"id":10,"email":"ff@ff.ru","first_name":"Елизавета","last_name":"Хопленлова","gender":"f","birth_date":406080000}'
curl localhost:3000/locations/new -XPOST -d '{"id":100,"distance":9,"city":"омск","place":"Рест","country":"элла"}'
curl localhost:3000/locations/new -XPOST -d '{"id":101,"distance":9,"city":"о1","place":"Рест11","country":"элла1"}'
curl localhost:3000/visits/new -XPOST -d '{"id": 1000, "user": 10, "visited_at": 1385572234, "location": 100, "mark": 2}'
curl localhost:3000/locations/100 -XPOST -d '{"place": "upd"}'
curl localhost:3000/users/10/visits
