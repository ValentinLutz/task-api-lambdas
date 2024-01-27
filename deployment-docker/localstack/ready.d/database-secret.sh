#!/bin/bash

awslocal secretsmanager create-secret \
  --name database-secret \
  --region eu-central-1 \
  --secret-string '{"username": "test", "password": "test"}'
