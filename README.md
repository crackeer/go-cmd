# 简介

使用Golang写的命令行工具，集合多种常用功能，持续增加中

# 功能列表

- base64-decode base64解码

```sh
got base64-decode 'aGVsbG8='
```

- base64-encode base64编码

```sh
got base64-encode 'hello'
```

- csv2json 解析csv文件到json

```sh
got csv2json test.csv
```

- decrypt-file 文件内容解密

```sh
got decrypt-file /your/input/file.txt.enc --output=/your/output/file.txt
```

- encrypt-file 文件内容加密

```sh
got encrypt-file /your/input/file.txt --output=/your/output/file.txt.enc
```

- go-md5 计算md5

```sh
got go-md5 'hello'
```

- http-post HttpPost请求

```sh
got http-post https://www.baidu.com body.json header.json
```

- http-server 简单的httpServer

```sh
got http-server --port=8888 /host/directory 
```

- ip 获取内网ip

```sh
got ip
```

- qrcode 生成二维码

```sh
got qrcode 'hello'
```

- regex-extract 正则提取

```sh
got regex-extract /input/file.txt 'hello'
```

- ssh-exec 远程执行

```sh
got ssh-exec test-server-name commands_alias
```

- ts 转换时间戳

```sh
got ts 3528532763523
```

- unzip 解压zip文件

```sh
got unzip test.zip
```

- zip 压缩文件

```sh
got zip test.zip
```

- download 下载文件 

```sh
got download file/url1 file/url2
```

- json-extract json提取

```sh
got json-extract file.json
```

- mysql-insert MySQL插入

```sh
got mysql-insert table table/data.json
```

```text
ctime         获取当前 time
download      批量下载文件
json-extract  json提取
mysql-insert  MySQL插入
mysql-query   MySQL查询
```
