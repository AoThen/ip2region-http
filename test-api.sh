#!/bin/bash

# Test script for ip2region-http API

BASE_URL="${BASE_URL:-http://localhost:8999}"

echo "Testing ip2region-http API..."
echo "================================"

# Test 1: Health check
echo -e "\n1. Health Check"
curl -s "${BASE_URL}/health" | jq .

# Test 2: Search with full format (default)
echo -e "\n2. Search IP (full format)"
curl -s "${BASE_URL}/search?ip=1.2.3.4" | jq .

# Test 3: Search with text format
echo -e "\n3. Search IP (text format)"
curl -s "${BASE_URL}/search?ip=120.229.45.92&format=text"
echo ""

# Test 4: Search with fields format
echo -e "\n4. Search IP (fields format)"
curl -s "${BASE_URL}/search?ip=8.8.8.8&format=fields" | jq .

# Test 5: Invalid IP
echo -e "\n5. Invalid IP (should return error)"
curl -s "${BASE_URL}/search?ip=999.999.999.999" | jq .

# Test 6: Missing IP parameter
echo -e "\n6. Missing IP parameter (should return error)"
curl -s "${BASE_URL}/search" | jq .

echo -e "\n================================"
echo "All tests completed!"
