#!/bin/bash

# Test script for CashCandy API

BASE_URL="http://localhost:1323"

echo "=== Testing CashCandy API ==="
echo

# Test 1: Hello World
echo "1. Testing Hello World endpoint..."
curl -s "$BASE_URL/" && echo
echo

# Test 2: Register a new user
echo "2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }')

echo "$REGISTER_RESPONSE"
echo

# Extract token from response
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Failed to get token, trying login..."
  
  # Test 3: Login
  echo "3. Testing user login..."
  LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "john@example.com",
      "password": "password123"
    }')
  
  echo "$LOGIN_RESPONSE"
  echo
  
  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ]; then
  echo "Failed to get authentication token. Exiting."
  exit 1
fi

echo "Token: $TOKEN"
echo

# Test 4: Create a goal
echo "4. Testing goal creation..."
GOAL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/goals" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Emergency Fund",
    "target_amount": 1000.00,
    "deadline": "2024-12-31T23:59:59Z"
  }')

echo "$GOAL_RESPONSE"
echo

# Extract goal ID
GOAL_ID=$(echo "$GOAL_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)

# Test 5: Get goals
echo "5. Testing get goals..."
curl -s -X GET "$BASE_URL/api/goals" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

# Test 6: Add transaction to goal
if [ ! -z "$GOAL_ID" ]; then
  echo "6. Testing add transaction to goal $GOAL_ID..."
  curl -s -X POST "$BASE_URL/api/transactions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"goal_id\": $GOAL_ID,
      \"amount\": 100.00,
      \"description\": \"Initial deposit\",
      \"type\": \"add\"
    }" | jq .
  echo
fi

# Test 7: Get dashboard stats
echo "7. Testing dashboard stats..."
curl -s -X GET "$BASE_URL/api/dashboard" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

echo "=== Test completed ==="
