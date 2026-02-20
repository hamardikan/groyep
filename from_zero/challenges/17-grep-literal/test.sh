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

echo "=== Challenge 17: Literal Match ==="
echo ""

# Test: single char match
output=$(./mygrep "J" testdata/rockbands.txt)
expected="Judas Priest
Bon Jovi
Junkyard"
assert_output "match 'J' in rockbands" "$expected" "$output"

# Test: multi-word
output=$(./mygrep "Bad" testdata/rockbands.txt)
expected="Bad English
Bad Company"
assert_output "match 'Bad' in rockbands" "$expected" "$output"

# Test: exact match
output=$(./mygrep "Nirvana" testdata/rockbands.txt)
assert_output "match 'Nirvana'" "Nirvana" "$output"

# Test: special chars
output=$(./mygrep "AC/DC" testdata/rockbands.txt)
assert_output "match 'AC/DC'" "AC/DC" "$output"

# Test: no match exit code
set +e
./mygrep "ZZZZZ" testdata/rockbands.txt > /dev/null 2>&1
assert_exit "no match exit code" "1" "$?"
set -e

# Test: match exit code
set +e
./mygrep "Nirvana" testdata/rockbands.txt > /dev/null 2>&1
assert_exit "match exit code" "0" "$?"
set -e

# Test: stdin pipe
output=$(echo -e "hello world\nfoo bar\nhello again" | ./mygrep "hello")
expected="hello world
hello again"
assert_output "stdin pipe match" "$expected" "$output"

# Test: stdin pipe no match
set +e
echo "hello world" | ./mygrep "zzz" > /dev/null 2>&1
assert_exit "stdin pipe no match" "1" "$?"
set -e

# Test: no args
set +e
./mygrep > /dev/null 2>&1
assert_exit "no args exit 2" "2" "$?"
set -e

# Test: file not found
set +e
./mygrep "pattern" nonexistent.txt > /dev/null 2>&1
assert_exit "file not found exit 2" "2" "$?"
set -e

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
