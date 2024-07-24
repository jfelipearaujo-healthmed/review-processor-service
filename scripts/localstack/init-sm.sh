#!/bin/sh

echo "Initializing Secrets Manager..."

awslocal secretsmanager create-secret \
    --name db-secret-url \
    --description "DB URL" \
    --secret-string "postgres://appointment:appointment@localhost:5432/appointment_db?sslmode=disable"

echo "Secrets Manager initialized!"