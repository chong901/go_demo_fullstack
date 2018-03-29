#!/bin/bash
BASENAME=${PWD##*/}
EXPORT="export"
filename=$BASENAME

case $GOOS in
 "windows")     filename+=".exe";;
esac

go build -o $filename

rm -rf $EXPORT/$BASENAME

mkdir -p $EXPORT/$BASENAME
mkdir -p $EXPORT/$BASENAME/configs
cp configs/*.toml $EXPORT/$BASENAME/configs
cp -r views $EXPORT/$BASENAME/
cp -r i18n $EXPORT/$BASENAME/
cp -r static $EXPORT/$BASENAME/
mv $filename $EXPORT/$BASENAME/
