#!/bin/bash

if ! type "make" > /dev/null; then
  echo "Make is needed to build systembolaget-api-fetch."
  echo "Check out if there's a release that suits your needs on GitHub!"
  exit 1
fi

if ! type "go" > /dev/null; then
  echo "Go is needed to build systembolaget-api-fetch."
  echo "Check out if there's a release that suits your needs on GitHub!"
  exit 1
fi

make
