#!/bin/bash
# Test the full pipeline: create deals with different deck types
# Usage: ./scripts/test_deal.sh

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "=== Creating Offering Memorandum ==="
RESPONSE=$(curl -s -X POST "$BASE_URL/api/deals" \
  -F "property_name=Bayshore Commerce Center" \
  -F "street=1250 Bayshore Blvd" \
  -F "city=San Francisco" \
  -F "state=CA" \
  -F "zip=94124" \
  -F "asset_class=office" \
  -F "units=10" \
  -F "sq_ft=11300" \
  -F "year_built=1985" \
  -F "deck_type=offering_memorandum" \
  -F "thesis=Strong value-add opportunity in emerging SoMa submarket. Below-market rents with 20% vacancy present clear lease-up upside. Recent infrastructure investment and tech tenant demand support aggressive rent growth assumptions." \
  -F "rent_roll=@sample_data/rent_roll.csv" \
  -F "t12=@sample_data/t12.csv")

DEAL_ID=$(echo "$RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)

if [ -z "$DEAL_ID" ]; then
  echo "Failed to create deal. Response:"
  echo "$RESPONSE"
  exit 1
fi

echo "OM Deal created: $DEAL_ID"
echo "View: $BASE_URL/api/deals/$DEAL_ID/deck"
echo ""

echo "=== Creating BOV ==="
BOV_RESPONSE=$(curl -s -X POST "$BASE_URL/api/deals" \
  -F "property_name=Bayshore Commerce Center" \
  -F "street=1250 Bayshore Blvd" \
  -F "city=San Francisco" \
  -F "state=CA" \
  -F "zip=94124" \
  -F "asset_class=office" \
  -F "units=10" \
  -F "sq_ft=11300" \
  -F "year_built=1985" \
  -F "deck_type=broker_opinion_of_value" \
  -F "thesis=Strong value-add opportunity" \
  -F "rent_roll=@sample_data/rent_roll.csv" \
  -F "t12=@sample_data/t12.csv")

BOV_ID=$(echo "$BOV_RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
echo "BOV Deal created: $BOV_ID"
echo "View: $BASE_URL/api/deals/$BOV_ID/deck"
echo ""

echo "=== Creating Investment Teaser ==="
TEASER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/deals" \
  -F "property_name=Bayshore Commerce Center" \
  -F "street=1250 Bayshore Blvd" \
  -F "city=San Francisco" \
  -F "state=CA" \
  -F "zip=94124" \
  -F "asset_class=office" \
  -F "deck_type=investment_teaser" \
  -F "thesis=Prime office asset in high-growth submarket" \
  -F "rent_roll=@sample_data/rent_roll.csv" \
  -F "t12=@sample_data/t12.csv")

TEASER_ID=$(echo "$TEASER_RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
echo "Teaser Deal created: $TEASER_ID"
echo "View: $BASE_URL/api/deals/$TEASER_ID/deck"
echo ""

echo "=== Listing all deals ==="
curl -s "$BASE_URL/api/deals" | python3 -m json.tool 2>/dev/null

echo ""
echo "Opening OM in browser..."
open "$BASE_URL/api/deals/$DEAL_ID/deck" 2>/dev/null || echo "(open manually)"
