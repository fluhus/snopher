# Tests a single library.

set -e

# Build library.
go build -buildmode=c-shared -o $1.so ./$1

# Run code.
python $1/$1.py > $1.txt

# Check output.
if [ "$(cat $1.txt)" != "$(cat outputs/$1.txt)" ]; then
  echo "Output is incorrect:"
  diff outputs/$1.txt $1.txt
  exit 1
fi

# Clean up.
rm $1.so $1.h $1.txt
