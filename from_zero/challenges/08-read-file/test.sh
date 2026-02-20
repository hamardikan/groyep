#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

PASS=0
FAIL=0

assert_stdout() {
    local description="$1"
    local expected="$2"
    shift 2
    local actual
    actual=$("$@" 2>/dev/null) || true

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

assert_stderr_contains() {
    local description="$1"
    local expected_substr="$2"
    shift 2
    local actual
    actual=$("$@" 2>&1 1>/dev/null) || true

    if echo "$actual" | grep -q "$expected_substr"; then
        echo -e "  ${GREEN}PASS${NC}: $description"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $description"
        echo "    expected stderr to contain: '$expected_substr'"
        echo "    got stderr: '$actual'"
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

echo "Building mycat..."
go build -o mycat .

echo ""
echo "=== Read File (mycat) Tests ==="

echo ""
echo "--- single file ---"
assert_stdout "print hello.txt" "Hello, World!" ./mycat testdata/hello.txt
assert_stdout "print numbers.txt" "$(printf '1\n2\n3')" ./mycat testdata/numbers.txt
assert_stdout "print empty.txt" "" ./mycat testdata/empty.txt

echo ""
echo "--- multiple files ---"
assert_stdout "concat hello + numbers" "$(printf 'Hello, World!\n1\n2\n3')" ./mycat testdata/hello.txt testdata/numbers.txt
assert_stdout "concat with empty in middle" "$(printf 'Hello, World!\n1\n2\n3')" ./mycat testdata/hello.txt testdata/empty.txt testdata/numbers.txt

echo ""
echo "--- error handling ---"
assert_stderr_contains "missing file error" "no such file or directory" ./mycat testdata/nonexistent.txt
assert_exit_code "missing file exits 1" 1 ./mycat testdata/nonexistent.txt
assert_exit_code "valid file exits 0" 0 ./mycat testdata/hello.txt
assert_exit_code "mixed valid+invalid exits 1" 1 ./mycat testdata/hello.txt testdata/nonexistent.txt

echo ""
echo "--- no args ---"
assert_stderr_contains "no args shows usage" "usage:" ./mycat
assert_exit_code "no args exits 1" 1 ./mycat

echo ""
echo "--- mixed valid and invalid continues processing ---"
MIXED_OUTPUT=$(./mycat testdata/hello.txt testdata/nonexistent.txt testdata/numbers.txt 2>/dev/null) || true
EXPECTED_MIXED="$(printf 'Hello, World!\n1\n2\n3')"
if [ "$MIXED_OUTPUT" = "$EXPECTED_MIXED" ]; then
    echo -e "  ${GREEN}PASS${NC}: continues after error"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: continues after error"
    echo "    expected stdout: '$EXPECTED_MIXED'"
    echo "    got stdout:      '$MIXED_OUTPUT'"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "================================"
echo -e "Results: ${GREEN}${PASS} passed${NC}, ${RED}${FAIL} failed${NC}"

rm -f mycat

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
