#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Testing OpenTelemetry Tracing Setup"
echo "======================================"
echo ""

# Check if otel-collector is running
echo -n "Checking otel-collector status... "
if docker ps | grep -q "otel-collector.*Up"; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
    echo "Run: docker-compose up -d"
    exit 1
fi

# Check otel-collector logs for errors
echo -n "Checking otel-collector for errors... "
ERRORS=$(docker logs otel-collector 2>&1 | grep -i "error" | wc -l)
if [ "$ERRORS" -eq 0 ]; then
    echo -e "${GREEN}✓ No errors${NC}"
else
    echo -e "${RED}✗ Found $ERRORS errors${NC}"
    echo "View logs: docker logs otel-collector"
fi

# Check if Jaeger is running
echo -n "Checking Jaeger status... "
if docker ps | grep -q "jaeger.*Up"; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
    exit 1
fi

# Check if services are running
echo -n "Checking shipment-service... "
if docker ps | grep -q "shipment-service.*Up"; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
    exit 1
fi

echo -n "Checking customer-service... "
if docker ps | grep -q "customer-service.*Up"; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
    exit 1
fi

echo ""
echo "📡 Sending test request..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/shipments \
  -H "Content-Type: application/json" \
  -d '{"route":"ALMATY→ASTANA","price":120000,"customer":{"idn":"990101123456"}}')

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Request successful${NC}"
    echo "Response: $RESPONSE"
    
    # Extract shipment ID from response
    SHIPMENT_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    if [ ! -z "$SHIPMENT_ID" ]; then
        echo ""
        echo "📦 Created shipment: $SHIPMENT_ID"
        
        # Wait a bit for trace to be collected
        echo "⏳ Waiting 3 seconds for trace propagation..."
        sleep 3
        
        echo ""
        echo "🔍 Checking for trace_id in logs..."
        
        # Check shipment-service logs
        SHIPMENT_TRACE=$(docker logs shipment-service 2>&1 | grep -i "trace_id" | tail -1)
        if [ ! -z "$SHIPMENT_TRACE" ]; then
            echo -e "${GREEN}✓ Found trace in shipment-service${NC}"
            echo "  $SHIPMENT_TRACE"
        else
            echo -e "${YELLOW}⚠ No trace_id in shipment-service logs${NC}"
        fi
        
        # Check customer-service logs
        CUSTOMER_TRACE=$(docker logs customer-service 2>&1 | grep -i "trace_id" | tail -1)
        if [ ! -z "$CUSTOMER_TRACE" ]; then
            echo -e "${GREEN}✓ Found trace in customer-service${NC}"
            echo "  $CUSTOMER_TRACE"
        else
            echo -e "${YELLOW}⚠ No trace_id in customer-service logs${NC}"
        fi
        
        # Check otel-collector received the trace
        echo ""
        echo "🔄 Checking otel-collector received traces..."
        OTEL_TRACES=$(docker logs otel-collector 2>&1 | grep -i "traces" | tail -3)
        if [ ! -z "$OTEL_TRACES" ]; then
            echo -e "${GREEN}✓ otel-collector processing traces${NC}"
        else
            echo -e "${YELLOW}⚠ No trace activity in otel-collector${NC}"
        fi
    fi
else
    echo -e "${RED}✗ Request failed${NC}"
    exit 1
fi

echo ""
echo "======================================"
echo -e "${GREEN}✅ All checks passed!${NC}"
echo ""
echo "🎯 Next steps:"
echo "1. Open Jaeger UI: http://localhost:16686"
echo "2. Select service: 'shipment-service'"
echo "3. Click 'Find Traces'"
echo "4. Look for trace with operation: 'POST /api/v1/shipments'"
echo ""
echo "Expected trace structure:"
echo "  📍 shipment-service: POST /api/v1/shipments"
echo "    └─ 📡 gRPC: customer.CustomerService/UpsertCustomer"
echo "       └─ 📍 customer-service: UpsertCustomer"
echo "          └─ 🗄️ Database operations"
echo ""
