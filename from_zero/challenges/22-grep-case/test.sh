#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
PASS=0
FAIL=0

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

assert_contains() {
    local desc="$1" expected="$2" actual="$3"
    if echo "$actual" | grep -qF "$expected"; then
        echo -e "  ${GREEN}PASS${NC}: $desc"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc (expected to contain '$expected')"
        FAIL=$((FAIL + 1))
    fi
}

echo "Building mygrep..."
go build -o mygrep . || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
echo ""

echo "=== Challenge 22: Case Insensitive ==="
echo ""

# THE SPEC TEST: case-sensitive A = 8 lines
count=$(./mygrep A testdata/rockbands.txt | wc -l | tr -d ' ')
assert_line_count "case-sensitive 'A' count" "8" "$count"

# THE SPEC TEST: case-insensitive A = 58 lines
count=$(./mygrep -i A testdata/rockbands.txt | wc -l | tr -d ' ')
assert_line_count "case-insensitive 'A' count" "58" "$count"

# Test: lowercase pattern finds uppercase
output=$(./mygrep -i "nirvana" testdata/rockbands.txt)
assert_contains "-i 'nirvana' finds Nirvana" "Nirvana" "$output"

# Test: uppercase pattern finds mixed case
output=$(./mygrep -i "NIRVANA" testdata/rockbands.txt)
assert_contains "-i 'NIRVANA' finds Nirvana" "Nirvana" "$output"

# Test: without -i, lowercase doesn't match uppercase
set +e
./mygrep "nirvana" testdata/rockbands.txt > /dev/null 2>&1
assert_exit "'nirvana' without -i = no match" "1" "$?"
set -e

# Test: -i with anchor
output=$(./mygrep -i '^a' testdata/rockbands.txt)
assert_contains "-i ^a finds AC/DC" "AC/DC" "$output"
assert_contains "-i ^a finds Aerosmith" "Aerosmith" "$output"

# Test: -ri combined
output=$(./mygrep -ri "nirvana" testdata/)
assert_contains "-ri finds in rockbands" "rockbands.txt" "$output"
assert_contains "-ri finds in BFS1985" "BFS1985.txt" "$output"

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
