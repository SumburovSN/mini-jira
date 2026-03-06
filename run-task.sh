#!/usr/bin/env bash
cd task-service
export DB_URL="postgres://postgres:postgres@localhost:5432/mini_jira"
export JWT_SECRET="supersecret"
export PORT="8083"
go run ./cmd