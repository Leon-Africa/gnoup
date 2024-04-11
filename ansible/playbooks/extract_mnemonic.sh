#!/bin/bash

# Read the contents of the mnemonic.txt file
text=$(cat /home/keys/mnemonic.txt)

# Extract the mnemonic phrase using sed
mnemonic_phrase=$(echo "$text" | sed -n '/\*\*IMPORTANT\*\*/,$p' | sed '1d;2d')

echo "$mnemonic_phrase"
