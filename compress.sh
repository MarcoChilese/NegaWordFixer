#!/bin/sh
echo About to compress.
tar -czvf $1 $2/html
echo Done.