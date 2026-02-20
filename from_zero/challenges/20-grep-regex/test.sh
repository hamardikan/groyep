#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
PASS=0
FAIL=0

assert_contains() {
    local desc="$1" expected="$2" actual="$3"
    if echo "$actual" | grep -qF "$expected"; then
        echo -e "  ${GREEN}PASS${NC}: $desc"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc"
        echo "    Expected to contain: $expected"
        FAIL=$((FAIL + 1))
    fi
}

assert_exit() {
    local desc="$1" expected="$2" actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo -e "  ${GREEN}PASS${NC}: $desc (exit $actual)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc (expected exit $expected, got $actual)"
        FAIL=$((FAIL + 1))
    fi
}

assert_line_count() {
    local desc="$1" expected="$2" actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo -e "  ${GREEN}PASS${NC}: $desc ($actual lines)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc (expected $expected lines, got $actual)"
        FAIL=$((FAIL + 1))
    fi
}

echo "Building mygrep..."
go build -o mygrep . || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
echo ""

echo "=== Challenge 20: Regular Expressions ==="
echo ""

# Test: \d in BFS1985
output=$(./mygrep '\d' testdata/test-subdir/BFS1985.txt)
assert_contains "\\d finds 'turned 24'" "turned 24" "$output"
assert_contains "\\d finds 'U2'" "U2" "$output"
assert_contains "\\d finds '1985'" "1985" "$output"
count=$(echo "$output" | wc -l | tr -d ' ')
assert_line_count "\\d line count in BFS1985" "9" "$count"

# Test: \w in symbols
output=$(./mygrep '\w' testdata/symbols.txt)
assert_contains "\\w finds 'pound'" "pound" "$output"
assert_contains "\\w finds 'dollar'" "dollar" "$output"
count=$(echo "$output" | wc -l | tr -d ' ')
assert_line_count "\\w line count in symbols" "2" "$count"

# Test: character class
output=$(./mygrep '[YZ]' testdata/rockbands.txt)
assert_contains "[YZ] finds Y&T" "Y&T" "$output"
assert_contains "[YZ] finds ZZ Top" "ZZ Top" "$output"

# Test: dot matches any char
output=$(./mygrep 'K.x' testdata/rockbands.txt)
assert_contains "K.x matches Kix" "Kix" "$output"

# Test: invalid regex
set +e
./mygrep '[invalid' testdata/rockbands.txt > /dev/null 2>&1
assert_exit "invalid regex exit 2" "2" "$?"
set -e

# Test: empty pattern still works
output=$(./mygrep '' testdata/symbols.txt)
expected=$(cat testdata/symbols.txt)
if [ "$output" = "$expected" ]; then
    echo -e "  ${GREEN}PASS${NC}: empty pattern still matches all"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: empty pattern should match all lines"
    FAIL=$((FAIL + 1))
fi

# Test: recursive with regex
output=$(./mygrep -r '\d{4}' testdata/)
assert_contains "recursive \\d{4} finds 1985" "1985" "$output"

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
