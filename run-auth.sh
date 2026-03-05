#!/usr/bin/env bash
cd auth-service
export DB_URL="postgres://postgres:postgres@localhost:5432/mini_jira"
export JWT_SECRET="supersecret"
export PORT="8081"
go run ./cmd
