# hdc-rss

这是一个将HDC下载筐中的种子下载到本地指定目录的小工具。兼容ttg的小货车。

## 安装方法

1. 安装golang并设置好$GOPATH
2. 执行 go get github.com/tominescu/hdc-rss/hdc-rss
3. 修改 $GOPATH/src/github.com/tominescu/hdc-rss/assets/sample-config.json, 填写你自己的rss地址.

## 运行

执行 $GOPATH/bin/hdc-rss -c $GOPATH/github.com/tominescu/hdc-rss/assets/sample-config.json
