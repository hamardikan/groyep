#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
PASS=0
FAIL=0

assert_exit() {
    local desc="$1" expected="$2" actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo -e "  ${GREEN}PASS${NC}: $desc"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC}: $desc (expected exit $expected, got $actual)"
        FAIL=$((FAIL + 1))
    fi
}

echo "Building mygrep..."
go build -o mygrep . || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
echo ""

echo "=== Challenge 23: Color ==="
echo ""

# Test: --color=always contains ANSI codes
output=$(./mygrep -color=always "Nirvana" testdata/rockbands.txt)
if echo "$output" | grep -qP '\033\[1;31m'; then
    echo -e "  ${GREEN}PASS${NC}: -color=always produces ANSI codes"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: -color=always should produce ANSI codes"
    FAIL=$((FAIL + 1))
fi

# Test: --color=never has no ANSI codes
output=$(./mygrep -color=never "Nirvana" testdata/rockbands.txt)
if echo "$output" | grep -qP '\033\['; then
    echo -e "  ${RED}FAIL${NC}: -color=never should NOT produce ANSI codes"
    FAIL=$((FAIL + 1))
else
    echo -e "  ${GREEN}PASS${NC}: -color=never has no ANSI codes"
    PASS=$((PASS + 1))
fi

# Test: no color flag = no ANSI codes
output=$(./mygrep "Nirvana" testdata/rockbands.txt)
if echo "$output" | grep -qP '\033\['; then
    echo -e "  ${RED}FAIL${NC}: no flag should NOT produce ANSI codes"
    FAIL=$((FAIL + 1))
else
    echo -e "  ${GREEN}PASS${NC}: no color flag = no ANSI codes"
    PASS=$((PASS + 1))
fi

# Test: color stripped = plain
color_out=$(./mygrep -color=always "Nirvana" testdata/rockbands.txt | sed 's/\x1b\[[0-9;]*m//g')
plain_out=$(./mygrep "Nirvana" testdata/rockbands.txt)
if [ "$color_out" = "$plain_out" ]; then
    echo -e "  ${GREEN}PASS${NC}: color stripped matches plain output"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: color stripped should match plain output"
    FAIL=$((FAIL + 1))
fi

# Test: previous features still work
set +e
./mygrep -r "Nirvana" testdata/ > /dev/null 2>&1
assert_exit "-r still works" "0" "$?"
./mygrep -i "nirvana" testdata/rockbands.txt > /dev/null 2>&1
assert_exit "-i still works" "0" "$?"
./mygrep -v "ZZZZ" testdata/rockbands.txt > /dev/null 2>&1
assert_exit "-v still works" "0" "$?"
set -e

rm -f mygrep
echo ""
echo "=== Results ==="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
[ $FAIL -gt 0 ] && exit 1 || true
