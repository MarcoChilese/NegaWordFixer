#!/bin/sh
echo About to compress.
mkdir html && mv $2/html/* html
tar -czvf $1 html
echo Done.