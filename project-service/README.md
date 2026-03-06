TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzI0Njk3OTMsInVzZXJfaWQiOjJ9.ZwD955GZE3ECwy1twvQWuWm4cIjrVhBkePqarPLaTrI

curl -X POST localhost:8082/projects \
  -d '{"name":"First Project"}' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"

curl localhost:8082/projects \
  -H "Authorization: Bearer $TOKEN"

curl -X POST localhost:8082/projects \
  -d '{"name":"Test"}' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"

curl -X DELETE localhost:8082/projects/5 \
  -H "Authorization: Bearer $TOKEN"

curl localhost:8082/projects/2 \
  -H "Authorization: Bearer $TOKEN"

curl -X PATCH localhost:8082/projects/1 \
  -d '{"name":"NewProject"}' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"