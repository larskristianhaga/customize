#!/bin/bash

# Test CORS functionality
# This script makes a preflight OPTIONS request to the API and checks if CORS headers are present in the response

# Replace with your actual API endpoint
API_ENDPOINT="http://localhost:8080/api/v1/custom/test-user-id"

echo "Testing CORS functionality..."
echo "Making preflight OPTIONS request to $API_ENDPOINT"

# Make a preflight OPTIONS request
response=$(curl -s -i -X OPTIONS \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  $API_ENDPOINT)

echo "Response headers:"
echo "$response"

# Check if CORS headers are present in the response
if echo "$response" | grep -q "Access-Control-Allow-Origin"; then
  echo "✅ CORS headers are present in the response"
else
  echo "❌ CORS headers are not present in the response"
fi

echo "Test completed."