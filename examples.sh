#!/bin/bash

# Create an Article
echo "Create an Article"
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"title":"abc","body":"xyz", "date": "2016-09-24", "tags": ["red", "white"}' \
  http://localhost:8000/articles
echo ""

# Get an article
echo "Get Article 1"
curl http://localhost:8000/articles/1
echo ""

# Get tags
echo "Get tags"
curl http://localhost:8000/tag/white/20160924
