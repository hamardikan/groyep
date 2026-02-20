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
        echo "    expected: $(echo "$expected" | head -3)"
        echo "    got:      $(echo "$actual" | head -3)"
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

echo "Building mynl..."
go build -o mynl .

echo ""
echo "=== Number Lines (mynl) Tests ==="

echo ""
echo "--- file input ---"

# Test sample.txt output
ACTUAL=$(./mynl testdata/sample.txt)
EXPECTED=$(printf '     1\tHello, World!\n     2\tThis is line two.\n     3\t\n     4\tThis is line four.\n     5\tThe end.')
assert_output "number lines sample.txt" "$EXPECTED" "$ACTUAL"

echo ""
echo "--- stdin pipe input ---"

# Test piped input
ACTUAL=$(printf 'first\nsecond\nthird\n' | ./mynl)
EXPECTED=$(printf '     1\tfirst\n     2\tsecond\n     3\tthird')
assert_output "pipe three lines" "$EXPECTED" "$ACTUAL"

# Test single line pipe
ACTUAL=$(echo "hello" | ./mynl)
EXPECTED=$(printf '     1\thello')
assert_output "pipe single line" "$EXPECTED" "$ACTUAL"

echo ""
echo "--- empty lines get numbers ---"

ACTUAL=$(printf 'a\n\nb\n' | ./mynl)
EXPECTED=$(printf '     1\ta\n     2\t\n     3\tb')
assert_output "empty line gets number" "$EXPECTED" "$ACTUAL"

echo ""
echo "--- tab separator ---"

# Verify output contains tab character
ACTUAL=$(echo "test" | ./mynl)
if echo "$ACTUAL" | grep -qP '\t'; then
    echo -e "  ${GREEN}PASS${NC}: output contains tab separator"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: output contains tab separator"
    echo "    output was: '$ACTUAL'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- error handling ---"

assert_exit_code "missing file exits 1" 1 ./mynl nonexistent.txt

# Verify error message goes to stderr
STDERR_OUTPUT=$(./mynl nonexistent.txt 2>&1 1>/dev/null) || true
if echo "$STDERR_OUTPUT" | grep -q "no such file or directory"; then
    echo -e "  ${GREEN}PASS${NC}: error message on stderr"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: error message on stderr"
    echo "    stderr was: '$STDERR_OUTPUT'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- right-aligned numbers ---"

# Line numbers should be right-aligned in 6-char field
ACTUAL=$(echo "test" | ./mynl)
if echo "$ACTUAL" | grep -q "^     1"; then
    echo -e "  ${GREEN}PASS${NC}: number is right-aligned (6 chars)"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: number is right-aligned (6 chars)"
    echo "    got: '$ACTUAL'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "================================"
echo -e "Results: ${GREEN}${PASS} passed${NC}, ${RED}${FAIL} failed${NC}"

rm -f mynl

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
