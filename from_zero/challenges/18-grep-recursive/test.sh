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
        echo "    Actual: $(echo "$actual" | head -3)"
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

echo "=== Challenge 18: Recursive Search ==="
echo ""

# Test: recursive Nirvana
output=$(./mygrep -r Nirvana testdata/)
assert_contains "recursive Nirvana finds rockbands" "testdata/rockbands.txt:Nirvana" "$output"
assert_contains "recursive Nirvana finds BFS1985" "testdata/test-subdir/BFS1985.txt:" "$output"

# Count total Nirvana matches (should be 5: 1 in rockbands + 4 in BFS1985)
count=$(echo "$output" | wc -l | tr -d ' ')
assert_line_count "recursive Nirvana total matches" "5" "$count"

# Test: recursive no match
set +e
./mygrep -r "ZZZZNOTFOUND" testdata/ > /dev/null 2>&1
assert_exit "recursive no match" "1" "$?"
set -e

# Test: recursive on single file
output=$(./mygrep -r Nirvana testdata/rockbands.txt)
assert_contains "recursive single file" "Nirvana" "$output"

# Test: recursive MTV
output=$(./mygrep -r MTV testdata/)
count=$(echo "$output" | wc -l | tr -d ' ')
assert_line_count "recursive MTV matches (>=4)" "4" "$([ "$count" -ge 4 ] && echo 4 || echo "$count")"

# Test: non-recursive still works
output=$(./mygrep "Nirvana" testdata/rockbands.txt)
assert_contains "non-recursive still works" "Nirvana" "$output"

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
