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
echo -e "${YELLOW}=== Challenge 13: Flag Parsing ===${NC}"
echo ""

# --- Structural checks ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "run: go mod init search" ""
fi

# --- Build ---
echo "Building..."
if go build -o "${BINARY}" . 2>/dev/null; then
    pass "Build succeeded"
else
    fail "Build failed" "compiles" "build errors"
    rm -f "${BINARY}"
    exit 1
fi

echo "Running tests..."

# ============================================================
# DEFAULT BEHAVIOR (no flags) — backward compatible with ch 12
# ============================================================

ACTUAL=$(./"${BINARY}" "Priest" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:Judas Priest"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Default: single match 'Priest'"
else
    fail "Default: single match 'Priest'" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" "Bad" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Default: multiple matches 'Bad'"
else
    fail "Default: multiple matches 'Bad'" "$EXPECTED" "$ACTUAL"
fi

# ============================================================
# -n FLAG (line numbers)
# ============================================================

ACTUAL=$(./"${BINARY}" -n "Priest" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:1:Judas Priest"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n: Priest is line 1"
else
    fail "-n: Priest is line 1" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -n "Bad" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:19:Bad English\ntestdata/rockbands.txt:25:Bad Company'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n: Bad at lines 19, 25"
else
    fail "-n: Bad at lines 19, 25" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -n "Nirvana" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:101:Nirvana"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n: Nirvana is line 101"
else
    fail "-n: Nirvana is line 101" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -n "cherry" testdata/fruits.txt 2>&1)
EXPECTED="testdata/fruits.txt:3:cherry"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n: cherry is line 3 in fruits"
else
    fail "-n: cherry is line 3 in fruits" "$EXPECTED" "$ACTUAL"
fi

# ============================================================
# -c FLAG (count only)
# ============================================================

ACTUAL=$(./"${BINARY}" -c "Priest" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:1"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-c: Priest count is 1"
else
    fail "-c: Priest count is 1" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -c "Bad" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:2"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-c: Bad count is 2"
else
    fail "-c: Bad count is 2" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -c "White" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:4"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-c: White count is 4"
else
    fail "-c: White count is 4" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -c "a" testdata/fruits.txt 2>&1)
EXPECTED="testdata/fruits.txt:4"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-c: 'a' in fruits count is 4"
else
    fail "-c: 'a' in fruits count is 4" "$EXPECTED" "$ACTUAL"
fi

# ============================================================
# -n -c COMBINED (count wins)
# ============================================================

ACTUAL=$(./"${BINARY}" -n -c "Bad" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:2"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-n -c: count wins (Bad = 2)"
else
    fail "-n -c: count wins" "$EXPECTED" "$ACTUAL"
fi

ACTUAL=$(./"${BINARY}" -c -n "Bad" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:2"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "-c -n: count wins (reversed order)"
else
    fail "-c -n: count wins (reversed)" "$EXPECTED" "$ACTUAL"
fi

# ============================================================
# NO MATCH with flags
# ============================================================

ACTUAL=$(./"${BINARY}" -n "zzz" testdata/rockbands.txt 2>&1 || true)
EXIT_CODE=$(./"${BINARY}" -n "zzz" testdata/rockbands.txt 2>/dev/null; echo $?)
if [ -z "$ACTUAL" ] && [ "$EXIT_CODE" = "1" ]; then
    pass "-n: no match exits 1"
else
    fail "-n: no match exits 1" "exit 1, no output" "exit=$EXIT_CODE, output='$ACTUAL'"
fi

ACTUAL=$(./"${BINARY}" -c "zzz" testdata/rockbands.txt 2>&1 || true)
EXIT_CODE=$(./"${BINARY}" -c "zzz" testdata/rockbands.txt 2>/dev/null; echo $?)
if [ -z "$ACTUAL" ] && [ "$EXIT_CODE" = "1" ]; then
    pass "-c: no match exits 1"
else
    fail "-c: no match exits 1" "exit 1, no output" "exit=$EXIT_CODE, output='$ACTUAL'"
fi

# ============================================================
# ERROR HANDLING
# ============================================================

ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: search <pattern> <filename>"
EXIT_CODE=$(./"${BINARY}" 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "No args: usage message, exit 2"
else
    fail "No args" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

ACTUAL=$(./"${BINARY}" -n 2>&1 || true)
EXPECTED="usage: search <pattern> <filename>"
EXIT_CODE=$(./"${BINARY}" -n 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "Flag but no args: usage message, exit 2"
else
    fail "Flag but no args" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

ACTUAL=$(./"${BINARY}" "hello" nonexistent.txt 2>&1 || true)
EXPECTED="error: cannot open 'nonexistent.txt'"
EXIT_CODE=$(./"${BINARY}" "hello" nonexistent.txt 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "File not found: error, exit 2"
else
    fail "File not found" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

ACTUAL=$(./"${BINARY}" -n "hello" nonexistent.txt 2>&1 || true)
EXPECTED="error: cannot open 'nonexistent.txt'"
EXIT_CODE=$(./"${BINARY}" -n "hello" nonexistent.txt 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "File not found with -n: error, exit 2"
else
    fail "File not found with -n" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
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
