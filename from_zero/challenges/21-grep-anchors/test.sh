#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
PASS=0
FAIL=0

assert_output() {
    local desc="$1" expected="$2" actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo -e "  ${GREEN}PASS${NC}: $desc"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc"
        echo "    Expected: $(echo "$expected" | head -2)"
        echo "    Actual:   $(echo "$actual" | head -2)"
        FAIL=$((FAIL + 1))
    fi
}

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

echo "=== Challenge 21: Anchors ==="
echo ""

# Test: ^A
output=$(./mygrep '^A' testdata/rockbands.txt)
expected="AC/DC
Aerosmith
Accept
April Wine
Autograph"
assert_output "^A matches bands starting with A" "$expected" "$output"

# Test: na$
output=$(./mygrep 'na$' testdata/rockbands.txt)
assert_output "na\$ matches Nirvana" "Nirvana" "$output"

# Test: ^$ (empty lines)
output=$(./mygrep '^$' testdata/test-subdir/BFS1985.txt)
count=$(echo "$output" | wc -l | tr -d ' ')
if [ "$count" -ge 4 ]; then
    echo -e "  ${GREEN}PASS${NC}: ^$ finds empty lines ($count)"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: ^$ expected >=4 empty lines, got $count"
    FAIL=$((FAIL + 1))
fi

# Test: ^.{3}$ (exactly 3 chars)
output=$(./mygrep '^.{3}$' testdata/rockbands.txt)
assert_contains "^.{3}$ finds UFO" "UFO" "$output"
assert_contains "^.{3}$ finds Kix" "Kix" "$output"
assert_contains "^.{3}$ finds TNT" "TNT" "$output"

# Test: anchored vs unanchored count
anchored_count=$(./mygrep '^A' testdata/rockbands.txt | wc -l | tr -d ' ')
unanchored_count=$(./mygrep 'A' testdata/rockbands.txt | wc -l | tr -d ' ')
if [ "$anchored_count" -lt "$unanchored_count" ]; then
    echo -e "  ${GREEN}PASS${NC}: ^A ($anchored_count) < A ($unanchored_count)"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: ^A ($anchored_count) should be less than A ($unanchored_count)"
    FAIL=$((FAIL + 1))
fi

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
