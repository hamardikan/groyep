#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="search"

pass() {
    echo -e "  ${GREEN}PASS${NC}: $1"
    PASS=$((PASS + 1))
}

fail() {
    echo -e "  ${RED}FAIL${NC}: $1"
    echo -e "        Expected: '$2'"
    echo -e "        Got:      '$3'"
    FAIL=$((FAIL + 1))
}

echo ""
echo -e "${YELLOW}=== Challenge 15: Packages (Project Structure) ===${NC}"
echo ""

# --- Structural checks ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "run: go mod init search" ""
fi

if [ -f "cmd/search/main.go" ]; then
    pass "cmd/search/main.go exists"
else
    fail "cmd/search/main.go not found" "cmd/search/main.go" ""
fi

if [ -f "internal/matcher/matcher.go" ]; then
    pass "internal/matcher/matcher.go exists"
else
    fail "internal/matcher/matcher.go not found" "internal/matcher/matcher.go" ""
fi

if [ -f "internal/matcher/matcher_test.go" ]; then
    pass "internal/matcher/matcher_test.go exists"
else
    fail "internal/matcher/matcher_test.go not found" "copy matcher_test.go from challenge dir" ""
fi

# --- Run package unit tests ---
echo ""
echo "Running matcher package unit tests..."
if go test ./internal/matcher/ -v 2>&1 | tail -1 | grep -q "^ok"; then
    pass "Matcher package tests pass"
else
    fail "Matcher package tests" "all pass" "some tests failed"
    echo ""
    echo "  Details:"
    go test ./internal/matcher/ -v 2>&1 | grep -E "FAIL|PASS" | head -20
fi

# --- Build binary ---
echo ""
echo "Building binary..."
if go build -o "${BINARY}" ./cmd/search 2>/dev/null; then
    pass "Build succeeded (go build ./cmd/search)"
else
    fail "Build failed" "compiles" "build errors"
    rm -f "${BINARY}"
    exit 1
fi

echo ""
echo "Running binary tests..."

# ============================================================
# SINGLE FILE
# ============================================================
echo ""
echo "--- Single file ---"

ACTUAL=$(./"${BINARY}" "Priest" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:Judas Priest"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Default: 'Priest'"
else
    fail "Default: 'Priest'" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" "Bad" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Default: 'Bad'"
else
    fail "Default: 'Bad'" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -n "Bad" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:19:Bad English\ntestdata/rockbands.txt:25:Bad Company'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n: 'Bad' at lines 19, 25"
else
    fail "-n: 'Bad' at lines 19, 25" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -c "Bad" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:2"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-c: 'Bad' count 2"
else
    fail "-c: 'Bad' count 2" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -n -c "Bad" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:2"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n -c: count wins"
else
    fail "-n -c: count wins" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" "zzz" testdata/rockbands.txt 2>&1 || true)
EXIT_CODE=$(./"${BINARY}" "zzz" testdata/rockbands.txt 2>/dev/null; echo $?)
if [ -z "$ACTUAL" ] && [ "$EXIT_CODE" = "1" ]; then
    pass "No match: exit 1"
else
    fail "No match" "exit 1, no output" "exit=$EXIT_CODE, output='$ACTUAL'"
fi

# ============================================================
# MULTI-FILE
# ============================================================
echo ""
echo "--- Multi-file ---"

ACTUAL=$(./"${BINARY}" "Bad" testdata/rockbands.txt testdata/fruits.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Multi: 'Bad' in first file only"
else
    fail "Multi: 'Bad' in first file only" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" "cherry" testdata/rockbands.txt testdata/fruits.txt 2>&1)
EXPECTED="testdata/fruits.txt:cherry"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Multi: 'cherry' in second file only"
else
    fail "Multi: 'cherry' in second file only" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -c "a" testdata/rockbands.txt testdata/fruits.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:52\ntestdata/fruits.txt:4'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Multi -c: 'a' count per file"
else
    fail "Multi -c: 'a' count per file" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" "zzz" testdata/rockbands.txt testdata/fruits.txt 2>&1 || true)
EXIT_CODE=$(./"${BINARY}" "zzz" testdata/rockbands.txt testdata/fruits.txt 2>/dev/null; echo $?)
if [ -z "$ACTUAL" ] && [ "$EXIT_CODE" = "1" ]; then
    pass "Multi: no match anywhere, exit 1"
else
    fail "Multi: no match anywhere" "exit 1, no output" "exit=$EXIT_CODE, output='$ACTUAL'"
fi

# ============================================================
# ERROR HANDLING
# ============================================================
echo ""
echo "--- Error handling ---"

ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: search <pattern> <filename>"
EXIT_CODE=$(./"${BINARY}" 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "No args: usage, exit 2"
else
    fail "No args" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

ACTUAL=$(./"${BINARY}" "hello" nonexistent.txt 2>&1 || true)
EXPECTED="error: cannot open 'nonexistent.txt'"
EXIT_CODE=$(./"${BINARY}" "hello" nonexistent.txt 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "File not found: error, exit 2"
else
    fail "File not found" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

# Multi-file with one missing
ACTUAL_STDOUT=$(./"${BINARY}" "Bad" testdata/rockbands.txt nonexistent.txt 2>/dev/null || true)
ACTUAL_STDERR=$(./"${BINARY}" "Bad" testdata/rockbands.txt nonexistent.txt 2>&1 1>/dev/null || true)
EXIT_CODE=$(./"${BINARY}" "Bad" testdata/rockbands.txt nonexistent.txt >/dev/null 2>/dev/null; echo $?)
CONTAINS_MATCH=false
CONTAINS_ERROR=false
if echo "$ACTUAL_STDOUT" | grep -q "Bad English"; then
    CONTAINS_MATCH=true
fi
if echo "$ACTUAL_STDERR" | grep -q "cannot open 'nonexistent.txt'"; then
    CONTAINS_ERROR=true
fi
if [ "$CONTAINS_MATCH" = "true" ] && [ "$CONTAINS_ERROR" = "true" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "Multi: one missing file — matches + error, exit 2"
else
    fail "Multi: one missing file" "matches + error, exit 2" "match=$CONTAINS_MATCH, error=$CONTAINS_ERROR, exit=$EXIT_CODE"
fi

# --- Cleanup ---
rm -f "${BINARY}"

# --- Summary ---
echo ""
TOTAL=$((PASS + FAIL))
if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}All ${TOTAL} tests passed!${NC}"
else
    echo -e "${RED}${FAIL} of ${TOTAL} tests failed.${NC}"
    exit 1
fi
