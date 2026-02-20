#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

PASS=0
FAIL=0

assert_output() {
    local description="$1"
    local expected="$2"
    shift 2
    local actual
    actual=$("$@" 2>&1) || true
    actual=$(echo "$actual" | tr -d '\r')
    expected=$(echo "$expected" | tr -d '\r')

    if [ "$actual" = "$expected" ]; then
        echo -e "  ${GREEN}PASS${NC}: $description"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $description"
        echo "    expected: '$expected'"
        echo "    got:      '$actual'"
        FAIL=$((FAIL + 1))
    fi
}

assert_exit_code() {
    local description="$1"
    local expected_code="$2"
    shift 2
    set +e
    "$@" > /dev/null 2>&1
    local actual_code=$?
    set -e

    if [ "$actual_code" -eq "$expected_code" ]; then
        echo -e "  ${GREEN}PASS${NC}: $description"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $description"
        echo "    expected exit code: $expected_code"
        echo "    got exit code:      $actual_code"
        FAIL=$((FAIL + 1))
    fi
}

echo "Building strtool..."
go build -o strtool .

echo ""
echo "=== String Checker Tests ==="

echo ""
echo "--- length command ---"
assert_output "length of 'Hello'" "Length: 5" ./strtool length "Hello"
assert_output "length of empty string" "Length: 0" ./strtool length ""
assert_output "length of 'hello world'" "Length: 11" ./strtool length "hello world"
assert_output "length of unicode 'café'" "Length: 4" ./strtool length "café"
assert_output "length of unicode '日本語'" "Length: 3" ./strtool length "日本語"

echo ""
echo "--- upper command ---"
assert_output "upper of 'hello'" "Upper: HELLO" ./strtool upper "hello"
assert_output "upper of 'Hello World'" "Upper: HELLO WORLD" ./strtool upper "Hello World"
assert_output "upper of 'café'" "Upper: CAFÉ" ./strtool upper "café"

echo ""
echo "--- lower command ---"
assert_output "lower of 'HELLO'" "Lower: hello" ./strtool lower "HELLO"
assert_output "lower of 'Go Is FUN'" "Lower: go is fun" ./strtool lower "Go Is FUN"

echo ""
echo "--- reverse command ---"
assert_output "reverse of 'hello'" "Reverse: olleh" ./strtool reverse "hello"
assert_output "reverse of 'racecar'" "Reverse: racecar" ./strtool reverse "racecar"
assert_output "reverse of 'café'" "Reverse: éfac" ./strtool reverse "café"
assert_output "reverse of empty string" "Reverse: " ./strtool reverse ""

echo ""
echo "--- contains command ---"
assert_output "contains 'hello world' 'world'" "Contains: true" ./strtool contains "hello world" "world"
assert_output "contains 'hello world' 'xyz'" "Contains: false" ./strtool contains "hello world" "xyz"
assert_output "contains case sensitive" "Contains: false" ./strtool contains "Hello" "hello"

echo ""
echo "--- error handling ---"
assert_output "unknown command" "error: unknown command 'explode'" ./strtool explode "test"
assert_output "no arguments" "usage: strtool <command> <string> [args...]" ./strtool
assert_exit_code "unknown command exits 1" 1 ./strtool explode "test"
assert_exit_code "no args exits 1" 1 ./strtool

echo ""
echo "================================"
echo -e "Results: ${GREEN}${PASS} passed${NC}, ${RED}${FAIL} failed${NC}"

# Clean up
rm -f strtool

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
