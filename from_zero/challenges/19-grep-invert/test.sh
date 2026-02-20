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

assert_not_contains() {
    local desc="$1" unexpected="$2" actual="$3"
    if echo "$actual" | grep -qF "$unexpected"; then
        echo -e "  ${RED}FAIL${NC}: $desc (found '$unexpected' in output)"
        FAIL=$((FAIL + 1))
    else
        echo -e "  ${GREEN}PASS${NC}: $desc"
        PASS=$((PASS + 1))
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

echo "Building mygrep..."
go build -o mygrep . || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
echo ""

echo "=== Challenge 19: Invert Match ==="
echo ""

# Test: -v excludes matching lines
output=$(./mygrep -v "o" testdata/symbols.txt)
assert_not_contains "-v excludes 'pound'" "pound" "$output"
assert_not_contains "-v excludes 'dollar'" "dollar" "$output"
assert_contains "-v keeps '!'" "!" "$output"

# Test: pipe composition (the spec test)
output=$(./mygrep -r Nirvana testdata/ | ./mygrep -v Madonna)
assert_contains "pipe: Nirvana without Madonna" "Nirvana" "$output"
assert_not_contains "pipe: no Madonna lines" "Madonna" "$output"

# Test: -v with no matches means all lines pass
output=$(./mygrep -v "ZZZNOTFOUND" testdata/rockbands.txt)
count=$(echo "$output" | wc -l | tr -d ' ')
if [ "$count" -ge 100 ]; then
    echo -e "  ${GREEN}PASS${NC}: -v with no-match pattern outputs all lines ($count)"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: -v with no-match pattern (expected >=100 lines, got $count)"
    FAIL=$((FAIL + 1))
fi

# Test: -v with empty pattern (all match) = no output, exit 1
set +e
output=$(./mygrep -v "" testdata/symbols.txt 2>/dev/null)
exitcode=$?
set -e
assert_exit "-v '' (all match) exit 1" "1" "$exitcode"

# Test: -rv combined
set +e
output=$(./mygrep -rv "e" testdata/test-subdir/BFS1985.txt 2>/dev/null)
exitcode=$?
set -e
assert_exit "-rv combined works" "0" "$exitcode"

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
