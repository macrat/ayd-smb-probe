#!/bin/bash

mkdir -p /share{/public,/private}

echo hello > /share/public/test.txt
touch -m -t 202301021504.05 /share/public/test.txt

touch -m -t 202301010101.01 /share/private/my-secret.txt
touch -m -t 202301020202.02 /share/private/hidden-file.txt

touch -m -t 202301010000.00 /share{,/public,/private}

exec /usr/bin/samba.sh \
  -p \
  -u "foo;bar" \
  -s "public;/share/public" \
  -s "private;/share/private;no;no;no;foo"
