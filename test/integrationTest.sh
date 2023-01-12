#!/bin/bash

response=$(curl --location --request GET 'localhost:8083/message' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request PUT 'localhost:8083/state' \
--header 'Content-Type: text/plain' \
--data-raw 'PAUSED' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request GET 'localhost:8083/state' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request PUT 'localhost:8083/state' \
--header 'Content-Type: text/plain' \
--data-raw 'RUNNING' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request GET 'localhost:8083/state' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request GET 'localhost:8083/node-statistic' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request GET 'localhost:8083/queue-statistic' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

response=$(curl --location --request GET 'localhost:8083/run-log' -s -o /dev/null -w "%{http_code}")
if [ "$response" -ne 200 ]; then
  exit 1
fi

exit 0
