if ! type "go" > /dev/null; then
  echo "Go is needed to run systembolaget-api-fetch"
  exit 1
fi

go run source/*.go
