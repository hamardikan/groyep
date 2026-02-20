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
        echo -e "  ${RED}FAIL${NC}: $desc (expected to contain '$expected')"
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
# Try cmd/mygrep first, then current dir
if [ -f "cmd/mygrep/main.go" ]; then
    go build -o mygrep ./cmd/mygrep || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
else
    go build -o mygrep . || { echo -e "${RED}BUILD FAILED${NC}"; exit 1; }
fi
echo ""

echo "========================================="
echo "  Challenge 25: THE FINAL TEST SUITE"
echo "========================================="
echo ""

# ===== STEP 1: Empty Match =====
echo "--- Step 1: Empty Match ---"
./mygrep "" testdata/test.txt | diff testdata/test.txt - > /dev/null 2>&1
assert_exit "empty pattern = entire file" "0" "$?"

# ===== STEP 2: Literal Match =====
echo "--- Step 2: Literal Match ---"
output=$(./mygrep "J" testdata/rockbands.txt)
expected="Judas Priest
Bon Jovi
Junkyard"
assert_output "J matches 3 bands" "$expected" "$output"

set +e
./mygrep "ZZZZZ" testdata/rockbands.txt > /dev/null 2>&1
assert_exit "no match exit 1" "1" "$?"
set -e

# ===== STEP 3: Recursive =====
echo "--- Step 3: Recursive ---"
output=$(./mygrep -r Nirvana testdata/)
assert_contains "recursive: rockbands" "testdata/rockbands.txt:Nirvana" "$output"
assert_contains "recursive: BFS1985" "testdata/test-subdir/BFS1985.txt:" "$output"
count=$(echo "$output" | grep -c "Nirvana")
assert_line_count "recursive: 5 Nirvana matches" "5" "$count"

# ===== STEP 4: Invert =====
echo "--- Step 4: Invert ---"
output=$(./mygrep -r Nirvana testdata/ | ./mygrep -v Madonna)
assert_output "Nirvana without Madonna" "testdata/rockbands.txt:Nirvana" "$output"

# ===== STEP 5: Regex =====
echo "--- Step 5: Regex ---"
output=$(./mygrep '\d' testdata/test-subdir/BFS1985.txt)
count=$(echo "$output" | wc -l | tr -d ' ')
assert_line_count "\\d matches 9 lines in BFS1985" "9" "$count"

output=$(./mygrep '\w' testdata/symbols.txt)
expected="pound
dollar"
assert_output "\\w matches pound and dollar" "$expected" "$output"

# ===== STEP 6: Anchors =====
echo "--- Step 6: Anchors ---"
output=$(./mygrep '^A' testdata/rockbands.txt)
expected="AC/DC
Aerosmith
Accept
April Wine
Autograph"
assert_output "^A matches 5 bands" "$expected" "$output"

output=$(./mygrep 'na$' testdata/rockbands.txt)
assert_output "na\$ matches Nirvana" "Nirvana" "$output"

# ===== STEP 7: Case Insensitive =====
echo "--- Step 7: Case Insensitive ---"
count=$(./mygrep A testdata/rockbands.txt | wc -l | tr -d ' ')
assert_line_count "case-sensitive A" "8" "$count"

count=$(./mygrep -i A testdata/rockbands.txt | wc -l | tr -d ' ')
assert_line_count "case-insensitive A" "58" "$count"

# ===== COMBINED FLAGS =====
echo "--- Combined Flags ---"
set +e
./mygrep -ri "nirvana" testdata/ > /dev/null 2>&1
assert_exit "-ri works" "0" "$?"

./mygrep -rv "ZZZZZ" testdata/ > /dev/null 2>&1
assert_exit "-rv works" "0" "$?"
set -e

# ===== STDIN =====
echo "--- Stdin ---"
output=$(echo -e "hello world\nfoo bar\nhello again" | ./mygrep "hello")
expected="hello world
hello again"
assert_output "stdin pipe" "$expected" "$output"

# ===== EXIT CODES =====
echo "--- Exit Codes ---"
set +e
./mygrep > /dev/null 2>&1
assert_exit "no args = exit 2" "2" "$?"

./mygrep "pattern" nonexistent.txt > /dev/null 2>&1
assert_exit "file not found = exit 2" "2" "$?"

./mygrep '[bad' testdata/rockbands.txt > /dev/null 2>&1
assert_exit "invalid regex = exit 2" "2" "$?"
set -e

# ===== LINE NUMBERS =====
echo "--- Line Numbers ---"
output=$(./mygrep -n "Nirvana" testdata/rockbands.txt)
if echo "$output" | grep -qP '^\d+:Nirvana$'; then
    echo -e "  ${GREEN}PASS${NC}: -n shows line numbers"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC}: -n should show NUMBER:Nirvana, got: $output"
    FAIL=$((FAIL + 1))
fi

# ===== COUNT =====
echo "--- Count ---"
output=$(./mygrep -c "Bad" testdata/rockbands.txt)
assert_output "-c counts 2 Bad matches" "2" "$(echo "$output" | tr -d ' ')"

# Cleanup
rm -f mygrep

echo ""
echo "========================================="
echo "  RESULTS"
echo "========================================="
echo -e "  ${GREEN}Passed: $PASS${NC}"
echo -e "  ${RED}Failed: $FAIL${NC}"
echo ""

if [ $FAIL -gt 0 ]; then
    echo "Some tests failed. Keep going!"
    exit 1
else
    echo "ALL TESTS PASSED!"
    echo ""
    echo "You built grep from scratch. Well done."
fi
