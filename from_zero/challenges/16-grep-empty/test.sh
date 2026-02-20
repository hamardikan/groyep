#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
PASS=0
FAIL=0

assert_output() {
    local desc="$1"
    local expected="$2"
    local actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo -e "  ${GREEN}PASS${NC}: $desc"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc"
        echo "    Expected: $(echo "$expected" | head -3)..."
        echo "    Actual:   $(echo "$actual" | head -3)..."
        FAIL=$((FAIL + 1))
    fi
}

assert_exit_code() {
    local desc="$1"
    local expected="$2"
    local actual="$3"
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

echo "=== Challenge 16: Empty Match ==="
echo ""

# Test 1: Empty pattern outputs entire file
echo "Test: Empty pattern matches all lines"
output=$(./mygrep "" testdata/test.txt)
expected=$(cat testdata/test.txt)
assert_output "empty pattern outputs all lines of test.txt" "$expected" "$output"

# Test 2: Empty pattern with diff
echo "Test: Empty pattern diff check"
./mygrep "" testdata/test.txt | diff testdata/test.txt - > /dev/null 2>&1
assert_exit_code "mygrep '' test.txt | diff test.txt - (no diff)" "0" "$?"

# Test 3: Empty pattern on rockbands.txt
echo "Test: Empty pattern on rockbands.txt"
./mygrep "" testdata/rockbands.txt | diff testdata/rockbands.txt - > /dev/null 2>&1
assert_exit_code "mygrep '' rockbands.txt | diff (no diff)" "0" "$?"

# Test 4: Literal match
echo "Test: Literal pattern match"
output=$(./mygrep "Nirvana" testdata/rockbands.txt)
assert_output "literal match 'Nirvana'" "Nirvana" "$output"

# Test 5: No match returns exit 1
echo "Test: No match exit code"
./mygrep "ZZZZNOTFOUND" testdata/rockbands.txt > /dev/null 2>&1 || true
exitcode=$?
# Re-run to capture actual exit code
set +e
./mygrep "ZZZZNOTFOUND" testdata/rockbands.txt > /dev/null 2>&1
exitcode=$?
set -e
assert_exit_code "no match returns exit 1" "1" "$exitcode"

# Test 6: No args returns exit 2
echo "Test: No args"
set +e
./mygrep > /dev/null 2>&1
exitcode=$?
set -e
assert_exit_code "no args returns exit 2" "2" "$exitcode"

# Test 7: File not found returns exit 2
echo "Test: File not found"
set +e
./mygrep "pattern" nonexistent.txt > /dev/null 2>&1
exitcode=$?
set -e
assert_exit_code "file not found returns exit 2" "2" "$exitcode"

# Cleanup
rm -f mygrep

echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"

if [ $FAIL -gt 0 ]; then
    exit 1
fi
