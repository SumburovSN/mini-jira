curl -X POST localhost:8081/register \
  -d '{"email":"admin@example.com","password":"123"}' \
  -H "Content-Type: application/json" 

curl -X POST localhost:8081/login \
  -d '{"email":"admin@example.com","password":"123"}' \
  -H "Content-Type: application/json" 

curl -X POST localhost:8081/login -d '{"email":"admin@example.com","password":"123"}' -H "Content-Type: application/json" 

TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzI0Njk3OTMsInVzZXJfaWQiOjJ9.ZwD955GZE3ECwy1twvQWuWm4cIjrVhBkePqarPLaTrI
curl localhost:8081/me -H "Authorization: Bearer $TOKEN"