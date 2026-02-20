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
        echo "    expected: '$(echo "$expected" | head -2)'"
        echo "    got:      '$(echo "$actual" | head -2)'"
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

echo "Building upper..."
go build -o upper .

echo ""
echo "=== Uppercase Filter Tests ==="

echo ""
echo "--- piped input ---"

ACTUAL=$(echo "hello world" | ./upper)
assert_output "simple lowercase pipe" "HELLO WORLD" "$ACTUAL"

ACTUAL=$(echo "Hello World" | ./upper)
assert_output "mixed case pipe" "HELLO WORLD" "$ACTUAL"

ACTUAL=$(printf "hello\nworld\n" | ./upper)
EXPECTED=$(printf "HELLO\nWORLD")
assert_output "multiline pipe" "$EXPECTED" "$ACTUAL"

ACTUAL=$(echo "abc 123 def" | ./upper)
assert_output "numbers pass through" "ABC 123 DEF" "$ACTUAL"

ACTUAL=$(echo "café au lait" | ./upper)
assert_output "unicode pipe" "CAFÉ AU LAIT" "$ACTUAL"

echo ""
echo "--- file input ---"

ACTUAL=$(./upper testdata/mixed.txt)
EXPECTED=$(printf 'HELLO WORLD\nTHIS IS LOWERCASE\nTHIS IS UPPERCASE\nMIXED CASE TEXT\nGO IS A GREAT LANGUAGE\n123 NUMBERS STAY THE SAME\nCAFÉ AU LAIT')
assert_output "read from file" "$EXPECTED" "$ACTUAL"

echo ""
echo "--- pipe from cat ---"

ACTUAL=$(cat testdata/mixed.txt | ./upper)
assert_output "cat file | upper" "$EXPECTED" "$ACTUAL"

echo ""
echo "--- file takes priority ---"

ACTUAL=$(echo "ignored" | ./upper testdata/mixed.txt)
if echo "$ACTUAL" | grep -q "HELLO WORLD"; then
    echo -e "  ${GREEN}PASS${NC}: file arg takes priority over stdin"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: file arg takes priority over stdin"
    echo "    got: '$ACTUAL'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- error handling ---"

assert_exit_code "missing file exits 1" 1 ./upper nonexistent.txt

STDERR_OUTPUT=$(./upper nonexistent.txt 2>&1 1>/dev/null) || true
if echo "$STDERR_OUTPUT" | grep -q "no such file or directory"; then
    echo -e "  ${GREEN}PASS${NC}: error message for missing file"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: error message for missing file"
    echo "    stderr: '$STDERR_OUTPUT'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- special characters preserved ---"

ACTUAL=$(echo "hello! @#\$ world?" | ./upper)
assert_output "special chars preserved" 'HELLO! @#$ WORLD?' "$ACTUAL"

echo ""
echo "================================"
echo -e "Results: ${GREEN}${PASS} passed${NC}, ${RED}${FAIL} failed${NC}"

rm -f upper

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
