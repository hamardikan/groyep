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
        echo -e "  ${RED}FAIL${NC}: $desc (expected to contain '$expected')"
        FAIL=$((FAIL + 1))
    fi
}

assert_min_lines() {
    local desc="$1" min="$2" actual="$3"
    local count
    count=$(echo "$actual" | grep -v '^--$' | wc -l | tr -d ' ')
    if [ "$count" -ge "$min" ]; then
        echo -e "  ${GREEN}PASS${NC}: $desc ($count lines >= $min)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc ($count lines < $min minimum)"
        FAIL=$((FAIL + 1))
    fi
}

echo "Building mygrep..."
go build -o mygrep . || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
echo ""

echo "=== Challenge 24: Context Lines ==="
echo ""

# Test: -A (after context)
output=$(./mygrep -A 1 "Nirvana" testdata/rockbands.txt)
assert_contains "-A 1 shows Nirvana" "Nirvana" "$output"

# Test: -B (before context)
output=$(./mygrep -B 2 "Accept" testdata/rockbands.txt)
assert_contains "-B 2 shows Accept" "Accept" "$output"
assert_min_lines "-B 2 shows >= 3 lines" "3" "$output"

# Test: -C (both context)
output=$(./mygrep -C 1 "Accept" testdata/rockbands.txt)
assert_contains "-C 1 shows Accept" "Accept" "$output"
assert_min_lines "-C 1 shows >= 3 lines" "3" "$output"

# Test: separator exists for non-adjacent groups
output=$(./mygrep -B 1 "1985" testdata/test-subdir/BFS1985.txt)
assert_contains "-B 1 '1985' has content" "1985" "$output"
if echo "$output" | grep -q '^--$'; then
    echo -e "  ${GREEN}PASS${NC}: group separator found"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: expected -- separator between groups"
    FAIL=$((FAIL + 1))
fi

# Test: -C 0 behaves like normal grep
context_out=$(./mygrep -C 0 "Nirvana" testdata/rockbands.txt | grep -v '^--$')
normal_out=$(./mygrep "Nirvana" testdata/rockbands.txt)
if [ "$context_out" = "$normal_out" ]; then
    echo -e "  ${GREEN}PASS${NC}: -C 0 matches normal grep"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: -C 0 should match normal grep"
    FAIL=$((FAIL + 1))
fi

# Test: -i with context
output=$(./mygrep -i -A 1 "nirvana" testdata/rockbands.txt)
assert_contains "-i with -A works" "Nirvana" "$output"

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
