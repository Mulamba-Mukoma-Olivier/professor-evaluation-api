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

TOKEN=$(echo "$LOGIN_BODY" | jq -r \
    '.data.access_token // .access_token // empty')

if [[ -z "$TOKEN" ]]; then
    echo "No access token returned" >&2
    exit 1
fi

echo "JWT successfully received"

AUTH_HEADER="Authorization: Bearer $TOKEN"

# ============================================================
# 5. Protected endpoint
# ============================================================

echo
echo "== 5. Protected endpoint =="

PROT_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "$AUTH_HEADER" \
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
# 8. Get professors
# ============================================================

echo
echo "== 8. Get professors =="

PROF_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "$AUTH_HEADER" \
    "$BASE/professors")

PROF_STATUS=$(echo "$PROF_RESP" | tail -n1)
PROF_BODY=$(echo "$PROF_RESP" | sed '$d')

echo "status: $PROF_STATUS"
echo "$PROF_BODY" | jq .

if [[ "$PROF_STATUS" != "200" ]]; then
    echo "Get professors failed" >&2
    exit 1
fi

PROFESSOR_ID=$(echo "$PROF_BODY" | jq -r '
    .professors[0].id //
    .data[0].id //
    empty
')

if [[ -z "$PROFESSOR_ID" ]]; then
    echo "No professor available for evaluation test" >&2
    exit 1
fi

echo "Professor ID: $PROFESSOR_ID"

# ============================================================
# 9. Get courses
# ============================================================

echo
echo "== 9. Get courses =="

COURSE_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "$AUTH_HEADER" \
    "$BASE/courses")

COURSE_STATUS=$(echo "$COURSE_RESP" | tail -n1)
COURSE_BODY=$(echo "$COURSE_RESP" | sed '$d')

echo "status: $COURSE_STATUS"
echo "$COURSE_BODY" | jq .

if [[ "$COURSE_STATUS" != "200" ]]; then
    echo "Get courses failed" >&2
    exit 1
fi

COURSE_ID=$(echo "$COURSE_BODY" | jq -r '
    .courses[0].id //
    .data[0].id //
    empty
')

if [[ -z "$COURSE_ID" ]]; then
    echo "No course available for evaluation test" >&2
    exit 1
fi

echo "Course ID: $COURSE_ID"

# ============================================================
# 10. Get active criteria
# ============================================================

echo
echo "== 10. Get active criteria =="

CRITERIA_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "$AUTH_HEADER" \
    "$BASE/criteria/active")

CRITERIA_STATUS=$(echo "$CRITERIA_RESP" | tail -n1)
CRITERIA_BODY=$(echo "$CRITERIA_RESP" | sed '$d')

echo "status: $CRITERIA_STATUS"
echo "$CRITERIA_BODY" | jq .

if [[ "$CRITERIA_STATUS" != "200" ]]; then
    echo "Get active criteria failed" >&2
    exit 1
fi

CRITERION_ID=$(echo "$CRITERIA_BODY" | jq -r '
    .criteria[0].id //
    .data[0].id //
    empty
')

if [[ -z "$CRITERION_ID" ]]; then
    echo "No active criterion available for evaluation test" >&2
    exit 1
fi

echo "Criterion ID: $CRITERION_ID"

# ============================================================
# 11. Create evaluation
# ============================================================

echo
echo "== 11. Create evaluation =="

EVALUATION_PAYLOAD=$(jq -n \
    --argjson professor_id "$PROFESSOR_ID" \
    --argjson course_id "$COURSE_ID" \
    --argjson criterion_id "$CRITERION_ID" \
    '{
        professor_id: $professor_id,
        course_id: $course_id,
        academic_year: "2025-2026",
        period: "E2E",
        answers: [
            {
                criterion_id: $criterion_id,
                score: 5
            }
        ]
    }'
)

echo "$EVALUATION_PAYLOAD" | jq .

EVALUATION_RESP=$(curl -sS -w "\n%{http_code}" \
    -X POST \
    "$BASE/evaluations" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d "$EVALUATION_PAYLOAD")

EVALUATION_STATUS=$(echo "$EVALUATION_RESP" | tail -n1)
EVALUATION_BODY=$(echo "$EVALUATION_RESP" | sed '$d')

echo "status: $EVALUATION_STATUS"
echo "$EVALUATION_BODY" | jq .

if [[ "$EVALUATION_STATUS" != "201" ]]; then
    echo "Create evaluation failed" >&2
    exit 1
fi

EVALUATION_ID=$(echo "$EVALUATION_BODY" | jq -r '.id // empty')

if [[ -z "$EVALUATION_ID" ]]; then
    echo "No evaluation ID returned" >&2
    exit 1
fi

echo "Evaluation ID: $EVALUATION_ID"

# ============================================================
# 12. Get evaluation by ID
# ============================================================

echo
echo "== 12. Get evaluation by ID =="

GET_EVALUATION_RESP=$(curl -sS -w "\n%{http_code}" \
    -H "$AUTH_HEADER" \
    "$BASE/evaluations/$EVALUATION_ID")

GET_EVALUATION_STATUS=$(echo "$GET_EVALUATION_RESP" | tail -n1)
GET_EVALUATION_BODY=$(echo "$GET_EVALUATION_RESP" | sed '$d')

echo "status: $GET_EVALUATION_STATUS"
echo "$GET_EVALUATION_BODY" | jq .

if [[ "$GET_EVALUATION_STATUS" != "200" ]]; then
    echo "Get evaluation failed" >&2
    exit 1
fi

# ============================================================
# 13. Duplicate evaluation
# ============================================================

echo
echo "== 13. Duplicate evaluation test =="

DUPLICATE_RESP=$(curl -sS -w "\n%{http_code}" \
    -X POST \
    "$BASE/evaluations" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d "$EVALUATION_PAYLOAD")

DUPLICATE_STATUS=$(echo "$DUPLICATE_RESP" | tail -n1)
DUPLICATE_BODY=$(echo "$DUPLICATE_RESP" | sed '$d')

echo "status: $DUPLICATE_STATUS"
echo "$DUPLICATE_BODY" | jq .

if [[ "$DUPLICATE_STATUS" != "409" ]]; then
    echo "Duplicate evaluation test failed" >&2
    exit 1
fi

# ============================================================
# SUCCESS
# ============================================================

echo
echo "========================================"
echo " ALL E2E TESTS PASSED"
echo "========================================"