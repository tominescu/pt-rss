#!/bin/bash
MD5='md5sum'
unamestr=`uname`
if [[ "$unamestr" == 'Darwin' ]]; then
	MD5='md5'
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.VERSION=$VERSION -s -w"
GCFLAGS="-B"

OSES=(linux darwin windows freebsd)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]
		then
			suffix=".exe"
		fi
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o pt-rss_${os}_${arch}${suffix} github.com/tominescu/pt-rss/pt-rss
		if $UPX; then upx -9 pt-rss_${os}_${arch}${suffix};fi
		tar -zcf pt-rss-${os}-${arch}-$VERSION.tar.gz pt-rss_${os}_${arch}${suffix} assets/sample-config.json
		$MD5 pt-rss-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o pt-rss_linux_arm$v  github.com/tominescu/pt-rss/pt-rss
done
if $UPX; then upx -9 pt-rss_linux_arm*;fi
tar -zcf pt-rss-linux-arm-$VERSION.tar.gz pt-rss_linux_arm* assets/sample-config.json
$MD5 pt-rss-linux-arm-$VERSION.tar.gz
