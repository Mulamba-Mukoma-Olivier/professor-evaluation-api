#!/usr/bin/env bash

set -euo pipefail

BASE="http://localhost:8080"

# ============================================================
# Vérification des dépendances
# ============================================================

command -v curl >/dev/null 2>&1 || {
    echo "curl is required" >&2
    exit 1
}

command -v jq >/dev/null 2>&1 || {
    echo "jq is required" >&2
    exit 1
}

echo "========================================"
echo " Professor Evaluation API - E2E TEST"
echo "========================================"

# ============================================================
# 1. Health check
# ============================================================

echo
echo "== 1. Health check =="

HEALTH_RESP=$(curl -sS -w "\n%{http_code}" \
    "$BASE/health")

HEALTH_STATUS=$(echo "$HEALTH_RESP" | tail -n1)
HEALTH_BODY=$(echo "$HEALTH_RESP" | sed '$d')

echo "status: $HEALTH_STATUS"
echo "$HEALTH_BODY" | jq .

if [[ "$HEALTH_STATUS" != "200" ]]; then
    echo "Health check failed" >&2
    exit 1
fi

# ============================================================
# 2. Register
# ============================================================

echo
echo "== 2. Register user =="

REG_PAYLOAD=$(jq -n \
    --arg matricule "20259998" \
    --arg name "E2E Tester" \
    --arg email "e2e+test@example.com" \
    --arg password "secret123" \
    --arg role "STUDENT" \
    '{
        matricule: $matricule,
        name: $name,
        email: $email,
        password: $password,
        role: $role
    }'
)

REG_RESP=$(curl -sS -w "\n%{http_code}" \
    -X POST \
    "$BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REG_PAYLOAD")

REG_STATUS=$(echo "$REG_RESP" | tail -n1)
REG_BODY=$(echo "$REG_RESP" | sed '$d')

echo "status: $REG_STATUS"
echo "$REG_BODY" | jq .

# 201 = nouvel utilisateur créé
# 400 = utilisateur probablement déjà existant
if [[ "$REG_STATUS" != "201" && "$REG_STATUS" != "400" ]]; then
    echo "Register failed with status $REG_STATUS" >&2
    exit 1
fi

# ============================================================
# 3. Login
# ============================================================

echo
echo "== 3. Login =="

LOGIN_PAYLOAD=$(jq -n \
    --arg email "e2e+test@example.com" \
    --arg password "secret123" \
    '{
        email: $email,
        password: $password
    }'
)

LOGIN_RESP=$(curl -sS -w "\n%{http_code}" \
    -X POST \
    "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "$LOGIN_PAYLOAD")

LOGIN_STATUS=$(echo "$LOGIN_RESP" | tail -n1)
LOGIN_BODY=$(echo "$LOGIN_RESP" | sed '$d')

echo "status: $LOGIN_STATUS"
echo "$LOGIN_BODY" | jq .

if [[ "$LOGIN_STATUS" != "200" ]]; then
    echo "Login failed" >&2
    exit 1
fi

# ============================================================
# 4. Extract JWT
# ============================================================

echo
echo "== 4. Extract JWT =="

TOKEN=$(echo "$LOGIN_BODY" | jq -r '.data.access_token // .access_token // empty')

if [[ -z "$TOKEN" ]]; then
    echo "No access token returned" >&2
    exit 1
fi

echo "JWT successfully received"

# ============================================================
# 5. Protected endpoint
# ============================================================

echo
echo "== 5. Protected endpoint =="

PROT_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE/professors")

PSTATUS=$(echo "$PROT_RESP" | tail -n1)
PBODY=$(echo "$PROT_RESP" | sed '$d')

echo "status: $PSTATUS"
echo "$PBODY" | jq .

if [[ "$PSTATUS" != "200" ]]; then
    echo "Protected endpoint failed" >&2
    exit 1
fi

# ============================================================
# 6. Test invalid token
# ============================================================

echo
echo "== 6. Invalid token test =="

INVALID_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "Authorization: Bearer invalid-token" \
    "$BASE/professors")

INVALID_STATUS=$(echo "$INVALID_RESP" | tail -n1)
INVALID_BODY=$(echo "$INVALID_RESP" | sed '$d')

echo "status: $INVALID_STATUS"
echo "$INVALID_BODY" | jq .

if [[ "$INVALID_STATUS" != "401" ]]; then
    echo "Invalid token test failed" >&2
    exit 1
fi

# ============================================================
# 7. Test missing token
# ============================================================

echo
echo "== 7. Missing token test =="

NO_TOKEN_RESP=$(curl -sS -w "\n%{http_code}" \
    "$BASE/professors")

NO_TOKEN_STATUS=$(echo "$NO_TOKEN_RESP" | tail -n1)
NO_TOKEN_BODY=$(echo "$NO_TOKEN_RESP" | sed '$d')

echo "status: $NO_TOKEN_STATUS"
echo "$NO_TOKEN_BODY" | jq .

if [[ "$NO_TOKEN_STATUS" != "401" ]]; then
    echo "Missing token test failed" >&2
    exit 1
fi

# ============================================================
# SUCCESS
# ============================================================

echo
echo "========================================"
echo " E2E TESTS PASSED"
echo "========================================"