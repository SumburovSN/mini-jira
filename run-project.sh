#!/usr/bin/env bash
cd project-service
export DB_URL="postgres://postgres:postgres@localhost:5432/mini_jira"
export JWT_SECRET="supersecret"
export PORT="8082"
go run ./cmd