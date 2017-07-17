#!/bin/bash
set -x
# Package building script
PROJECT_ROOT=/app

rm -rf platform-visibility-all* || true
fpm -s dir -t deb \
  -n platform-visibility-all \
  -v $VERSION \
  --after-install /app/tools/build/debian/postinstall \
  $PROJECT_ROOT/etc/platformvisibility/=/etc/platformvisibility/ \
  $PROJECT_ROOT/build/linux-amd64/visualizationapi=/usr/bin/ \
  $PROJECT_ROOT/build/linux-amd64/sql-migrate=/usr/bin/ \
  $PROJECT_ROOT/tools/database-migrations/=/var/lib/platformvisibility/database-migrations/

mkdir -p $PROJECT_ROOT/build/deb
mv platform-visibility-all* $PROJECT_ROOT/build/deb/
