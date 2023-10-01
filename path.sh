
# Get the current directory
CURRENT_DIRECTORY=$(pwd)

# Path to the directory you want to add to PATH
DIRECTORY_PATH="$CURRENT_DIRECTORY/build"

# Check if the directory exists
if [ -d "$DIRECTORY_PATH" ]; then
    # Add directory to PATH if not already added
    if [[ ":$PATH:" != *":$DIRECTORY_PATH:"* ]]; then
        PATH="$DIRECTORY_PATH:$PATH"
        echo "Directory added to PATH. You can now run executables from the 'build' directory."
    else
        echo "Directory is already in PATH."
    fi
else
    echo "Error: Directory not found."
fi