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
echo -e "${YELLOW}=== Challenge 12: CLI Args (Line Search Tool) ===${NC}"
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

# --- Test: single match ---
ACTUAL=$(./"${BINARY}" "Priest" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:Judas Priest"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Single match: 'Priest'"
else
    fail "Single match: 'Priest'" "$EXPECTED" "$ACTUAL"
fi

# --- Test: multiple matches ---
ACTUAL=$(./"${BINARY}" "Bad" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Multiple matches: 'Bad'"
else
    fail "Multiple matches: 'Bad'" "$EXPECTED" "$ACTUAL"
fi

# --- Test: White matches ---
ACTUAL=$(./"${BINARY}" "White" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:Whitesnake\ntestdata/rockbands.txt:Great White\ntestdata/rockbands.txt:White Lion\ntestdata/rockbands.txt:Whitecross'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Multiple matches: 'White'"
else
    fail "Multiple matches: 'White'" "$EXPECTED" "$ACTUAL"
fi

# --- Test: match in fruits file ---
ACTUAL=$(./"${BINARY}" "berry" testdata/fruits.txt 2>&1)
EXPECTED=$'testdata/fruits.txt:blueberry'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Match in fruits: 'berry'"
else
    fail "Match in fruits: 'berry'" "$EXPECTED" "$ACTUAL"
fi

# --- Test: multiple matches in fruits ---
ACTUAL=$(./"${BINARY}" "a" testdata/fruits.txt 2>&1)
EXPECTED=$'testdata/fruits.txt:apple\ntestdata/fruits.txt:banana\ntestdata/fruits.txt:apricot\ntestdata/fruits.txt:avocado'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Multiple matches in fruits: 'a'"
else
    fail "Multiple matches in fruits: 'a'" "$EXPECTED" "$ACTUAL"
fi

# --- Test: no match ---
ACTUAL=$(./"${BINARY}" "zzz" testdata/rockbands.txt 2>&1 || true)
EXIT_CODE=$(./"${BINARY}" "zzz" testdata/rockbands.txt 2>/dev/null; echo $?)
if [ -z "$ACTUAL" ] && [ "$EXIT_CODE" = "1" ]; then
    pass "No match: exit code 1, no output"
else
    fail "No match: 'zzz'" "exit code 1, no output" "exit=$EXIT_CODE, output='$ACTUAL'"
fi

# --- Test: case sensitive ---
ACTUAL=$(./"${BINARY}" "priest" testdata/rockbands.txt 2>&1 || true)
EXIT_CODE=$(./"${BINARY}" "priest" testdata/rockbands.txt 2>/dev/null; echo $?)
if [ -z "$ACTUAL" ] && [ "$EXIT_CODE" = "1" ]; then
    pass "Case sensitive: 'priest' (lowercase) no match"
else
    fail "Case sensitive" "exit code 1, no output" "exit=$EXIT_CODE, output='$ACTUAL'"
fi

# --- Test: no arguments ---
ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: search <pattern> <filename>"
EXIT_CODE=$(./"${BINARY}" 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "No arguments: usage message, exit code 2"
else
    fail "No arguments" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

# --- Test: only pattern, missing filename ---
ACTUAL=$(./"${BINARY}" "hello" 2>&1 || true)
EXPECTED="usage: search <pattern> <filename>"
EXIT_CODE=$(./"${BINARY}" "hello" 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "Missing filename: usage message, exit code 2"
else
    fail "Missing filename" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

# --- Test: file not found ---
ACTUAL=$(./"${BINARY}" "hello" nonexistent.txt 2>&1 || true)
EXPECTED="error: cannot open 'nonexistent.txt'"
EXIT_CODE=$(./"${BINARY}" "hello" nonexistent.txt 2>/dev/null; echo $?)
if [ "$ACTUAL" = "$EXPECTED" ] && [ "$EXIT_CODE" = "2" ]; then
    pass "File not found: error message, exit code 2"
else
    fail "File not found" "$EXPECTED (exit 2)" "$ACTUAL (exit $EXIT_CODE)"
fi

# --- Test: special characters in pattern ---
ACTUAL=$(./"${BINARY}" "/" testdata/rockbands.txt 2>&1)
EXPECTED=$'testdata/rockbands.txt:AC/DC\ntestdata/rockbands.txt:Love/Hate'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Special chars in pattern: '/'"
else
    fail "Special chars: '/'" "$EXPECTED" "$ACTUAL"
fi

# --- Test: match at end of file ---
ACTUAL=$(./"${BINARY}" "Nirvana" testdata/rockbands.txt 2>&1)
EXPECTED="testdata/rockbands.txt:Nirvana"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Match at end of file: 'Nirvana'"
else
    fail "Match at end of file" "$EXPECTED" "$ACTUAL"
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
