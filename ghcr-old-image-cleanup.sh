#!/bin/bash

ORG="vetchium"
IMAGES=$(gh api "orgs/$ORG/packages?package_type=container" --jq '.[].name')

for IMAGE in $IMAGES; do
  echo "Processing image: $IMAGE"

  # Get all versions of this image, sorted by creation date (latest first)
  VERSIONS=$(gh api "orgs/$ORG/packages/container/$IMAGE/versions" \
              --paginate \
              --jq 'sort_by(.created_at) | reverse | .[].id')

  COUNT=0
  for VERSION_ID in $VERSIONS; do
    ((COUNT++))
    if [ "$COUNT" -le 3 ]; then
      echo "  Keeping version ID: $VERSION_ID"
    else
      echo "  Deleting old version ID: $VERSION_ID"
      gh api --method DELETE "orgs/$ORG/packages/container/$IMAGE/versions/$VERSION_ID"
    fi
  done
done
