#!/bin/bash
# Create dummy input
seq 10 | awk '{print "cmd" $1}' > input.txt

# Run ch with d3 output
cat input.txt | ./ch --output d3 > d3_output.html 2>&1

# Extract the HTML file path
HTML_FILE=$(grep "Opening D3 chart at" d3_output.html | head -n 1 | awk '{print $5}')

if [ -z "$HTML_FILE" ]; then
  echo "FAIL: Could not find HTML file path in output"
  cat d3_output.html
  exit 1
fi

echo "Checking HTML file: $HTML_FILE"

# Check for NaN in the file (as seen in user screenshot)
if grep -q "NaN" "$HTML_FILE"; then
  echo "FAIL: NaN detected in D3 output"
  exit 1
fi

echo "PASS: No NaN found"
