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
    local actual="$3"

    if [ "$actual" = "$expected" ]; then
        echo -e "  ${GREEN}PASS${NC}: $description"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $description"
        echo "    expected: '$(echo "$expected" | head -3)'"
        echo "    got:      '$(echo "$actual" | head -3)'"
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

assert_exit_code_piped() {
    local description="$1"
    local expected_code="$2"
    local input="$3"
    local pattern="$4"
    set +e
    echo "$input" | ./finder "$pattern" > /dev/null 2>&1
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

echo "Building finder..."
go build -o finder .

echo ""
echo "=== Substring Finder (Mini Grep) Tests ==="

echo ""
echo "--- file search: matches found ---"

ACTUAL=$(./finder "ap" testdata/fruits.txt)
EXPECTED=$(printf 'apple\napricot')
assert_output "find 'ap' in fruits" "$EXPECTED" "$ACTUAL"

ACTUAL=$(./finder "berry" testdata/fruits.txt)
assert_output "find 'berry' in fruits" "blueberry" "$ACTUAL"

ACTUAL=$(./finder "banana" testdata/fruits.txt)
assert_output "find exact 'banana'" "banana" "$ACTUAL"

echo ""
echo "--- file search: no matches ---"

ACTUAL=$(./finder "xyz" testdata/fruits.txt 2>/dev/null) || true
assert_output "no match 'xyz' produces no output" "" "$ACTUAL"

ACTUAL=$(./finder "Apple" testdata/fruits.txt 2>/dev/null) || true
assert_output "case sensitive 'Apple' no match" "" "$ACTUAL"

echo ""
echo "--- empty pattern matches all ---"

ACTUAL=$(./finder "" testdata/fruits.txt)
EXPECTED=$(printf 'apple\nbanana\ncherry\napricot\nblueberry\navocado')
assert_output "empty pattern matches all lines" "$EXPECTED" "$ACTUAL"

echo ""
echo "--- stdin input ---"

ACTUAL=$(echo "hello world" | ./finder "world")
assert_output "stdin match 'world'" "hello world" "$ACTUAL"

ACTUAL=$(printf "hello\nworld\nfoo\n" | ./finder "o")
EXPECTED=$(printf 'hello\nworld\nfoo')
assert_output "stdin multiple matches" "$EXPECTED" "$ACTUAL"

ACTUAL=$(echo "hello world" | ./finder "xyz" 2>/dev/null) || true
assert_output "stdin no match" "" "$ACTUAL"

echo ""
echo "--- exit codes ---"

assert_exit_code "match found exits 0" 0 ./finder "apple" testdata/fruits.txt
assert_exit_code "no match exits 1" 1 ./finder "xyz" testdata/fruits.txt
assert_exit_code "no args exits 2" 2 ./finder
assert_exit_code "file not found exits 2" 2 ./finder "test" nonexistent.txt

echo ""
echo "--- exit codes with piped input ---"

assert_exit_code_piped "stdin match exits 0" 0 "hello world" "hello"
assert_exit_code_piped "stdin no match exits 1" 1 "hello world" "xyz"

echo ""
echo "--- error messages ---"

STDERR_OUTPUT=$(./finder 2>&1 1>/dev/null) || true
if echo "$STDERR_OUTPUT" | grep -q "usage: finder"; then
    echo -e "  ${GREEN}PASS${NC}: no args shows usage message"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: no args shows usage message"
    echo "    stderr: '$STDERR_OUTPUT'"
    FAIL=$((FAIL + 1))
fi

STDERR_OUTPUT=$(./finder "test" nonexistent.txt 2>&1 1>/dev/null) || true
if echo "$STDERR_OUTPUT" | grep -q "no such file or directory"; then
    echo -e "  ${GREEN}PASS${NC}: file not found error message"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: file not found error message"
    echo "    stderr: '$STDERR_OUTPUT'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- multiple matches on pattern 'a' ---"

ACTUAL=$(./finder "a" testdata/fruits.txt)
EXPECTED=$(printf 'banana\napricot\navocado')
assert_output "pattern 'a' matches banana, apricot, avocado" "$EXPECTED" "$ACTUAL"

echo ""
echo "================================"
echo -e "Results: ${GREEN}${PASS} passed${NC}, ${RED}${FAIL} failed${NC}"

rm -f finder

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
