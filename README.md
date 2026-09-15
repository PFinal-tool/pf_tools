# PF_Tools

基于 golang 终端小工具集

![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311081714504.png)

#### 查询天气

```shell
 pf_tools pf_wt  # 默认本地天气
 pf_tools pf_wt 成都 # 指定城市天气
```
![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311081716928.png)

#### 查询手机归属地

```shell
 pf_tools pf_m 17721123677 查询手机归属地
```

![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311091002314.png)

#### md5 小工具

```shell
 pf_tools pf_md5 pfinaclub #加密 pfinalclub
 pf_tools pf_md5 pfinaclub pf 1 #加密 pfinalclub 右边加盐 pf
```
> 按 q 退出

> 按 c 建复制到剪贴板上

![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311141146354.png)

![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311141148825.png)

#### base64 小工具

```shell
pf_tools pf_b64 pfinalclub # 加密
pf_tools pf_b64 cGZpbmFsY2x1Yg== -d # 解密
```
> 按 q 退出

> 按 c 建复制到剪贴板上

![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311141431206.png)

54ix55Sf5rS754ix5Luj56CBLOWWnOasouaKmOiFvu+8jOeXtOi/t+S6juaAp+iDveWSjOaVj+aNt+W8gOWPkSzmiJHllpzmrKLor7vnvJbnoIHjgIHor63kuYnlkozlj6/ph43nlKjmgKfjgII=


#### 查询网络词汇

```shell
pf_tools pf_s yyds
pf_tools pf_s qnmd
```
![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202311151407598.png)


#### 清除目录中的.DS_Store 文件
```shell
pf_tools pf_cd  # 清除当前目录
pf_tools pf_cd "" # 清除指定路径
```
![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202401250948067.png)


#### 微信URL检测
```shell
pf_tools pf_cwx [urlstring]
```


![](https://raw.githubusercontent.com/pfinal-nc/iGallery/master/blog/202403050921393.png)

#### JSON 工具
```shell
pf_tools pf_json '{"name":"test","age":18}'
pf_tools pf_json file.json -f
pf_tools pf_json '{"a":1}' -c
pf_tools pf_json '{"a":{"b":1}}' -p a.b -q
```

#### URL 编解码
```shell
pf_tools pf_url "https://example.com?name=张三"
pf_tools pf_url "https://a.com?x=%E5%BC%A0" -d
pf_tools pf_url "https://a.com/path?a=1&b=2" --parse -q
```

#### JWT 解码
```shell
pf_tools pf_jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
pf_tools pf_jwt $TOKEN --exp -q
```

#### UUID / 随机串 / 密码
```shell
pf_tools pf_uuid -c 5 -q
pf_tools pf_uuid -n -q
pf_tools pf_rand 32 -q
pf_tools pf_rand 16 -t hex -q
pf_tools pf_pass -l 20 -s -q
```

#### 端口排查
```shell
pf_tools pf_whichport 8080
pf_tools pf_killport 8080
pf_tools pf_killport 8080 -f
```

#### TLS 证书
```shell
pf_tools pf_cert example.com
pf_tools pf_cert example.com --days -q
```

#### Cron / 状态码 / Hex / 命名 / Unicode
```shell
pf_tools pf_cron "0 */5 * * *" -q
pf_tools pf_status 429 -q
pf_tools pf_hex "hello" -q
pf_tools pf_case user_id --to camel -q
pf_tools pf_unicode "你好" -q
```

#### Cookie / Header / 剪贴板 / 短链 / Mock / IP
```shell
pf_tools pf_cookie "a=1; b=2" -q
pf_tools pf_mock user -c 2 -q
pf_tools pf_mock phone -c 5 -q
pf_tools pf_ip -q
pf_tools pf_ip -l
pf_tools pf_expand https://example.com -q
pf_tools pf_clip jwt -q
```

#### DNS / HTTP / CIDR / SemVer
```shell
pf_tools pf_dns example.com -q
pf_tools pf_dns example.com -t MX -q
pf_tools pf_http https://example.com --time -q
pf_tools pf_http https://httpbin.org/post -X POST -d '{"a":1}' -j
pf_tools pf_cidr 10.0.0.8 10.0.0.0/24 -q
pf_tools pf_cidr 10.0.0.0/24 --list -n 5 -q
pf_tools pf_semver 1.2.3 --minor -q
pf_tools pf_semver v1.2.3-beta.1 --inc-prerelease -q
```

#### 管道友好（-q）
多数编解码类命令支持 `-q`，只打印结果，方便脚本：
```shell
pf_tools pf_md5 hello -q
pf_tools pf_b64 hello -q | pbcopy
pf_tools pf_t -m -q
pf_tools pf_b64 hello -u -q   # Base64URL
pf_tools pf_md5 -f ./go.mod -q
```

时间戳增强：
```shell
pf_tools pf_t -q                 # 当前秒级时间戳
pf_tools pf_t -m -q              # 毫秒时间戳
pf_tools pf_t 1699123456 -q      # 戳 → 时间
pf_tools pf_t "2023-11-05 10:30:00" -q
```

## Stargazers over time
[![Stargazers over time](https://starchart.cc/PFinal-tool/pf_tools.svg?variant=adaptive)](https://starchart.cc/PFinal-tool/pf_tools)
