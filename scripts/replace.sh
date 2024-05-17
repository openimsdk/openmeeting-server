#!/bin/bash

# Use the parent directory of the current directory as the folder path
folder_path=$(dirname "$(pwd)")

# Define the strings to be replaced
old_string="openmeeting-server"
new_string="openmeeting-server"

# Debugging: Print the folder path
echo "Folder path: $folder_path"

# Traverse all files in the folder and replace the string
find "$folder_path" -type f | while read -r file; do
    if [[ -f "$file" ]]; then
        # Debugging: Print the current file being processed
        echo "Processing file: $file"

        # Use sed to replace the string in the file
        if grep -q "$old_string" "$file"; then
            sed -i "s/$old_string/$new_string/g" "$file"
            echo "Replaced in $file"
        else
            echo "No match found in $file"
        fi
    fi
done

echo "String replacement completed."
