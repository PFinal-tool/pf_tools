# 开发者实战工具命令建议

基于现有项目代码结构（`cmd/` + `pak/` + cobra + termui + 剪贴板），从开发者日常实战痛点出发给出建议。

## 现有能力盘点

| 命令 | 能力 | 风格特征 |
|------|------|----------|
| `pf_wt` | 天气查询 | 外部 HTTP |
| `pf_m` | 手机归属地 | 本地数据 `phone.dat` |
| `pf_md5` | MD5（可加盐） | termui + `c` 复制 |
| `pf_b64` | Base64 编解码 | termui + `c` 复制 |
| `pf_s` | 网络词汇释义 | 外部查询 |
| `pf_t` | 时间戳 / 时间→戳 | termui + `c` 复制 |
| `pf_cd` | 清 `.DS_Store` | 本地文件操作 |
| `pf_cwx` | 微信域名拦截检测 | 外部 HTTP + 表格 |
| `pf_mk` | 按日期建目录 | 本地文件操作 |

**项目 DNA**：单命令、秒开、结果可复制、偏「调试瞬间刚需」，不是替代 `jq`/`curl`/`git` 全家桶。

---

## 📋 目录

- [先补现有工具的缺口](#先补现有工具的缺口)
- [编码/加密工具](#编码加密工具)
- [网络工具](#网络工具)
- [文件操作工具](#文件操作工具)
- [开发辅助工具](#开发辅助工具)
- [系统工具](#系统工具)
- [数据转换工具](#数据转换工具)
- [第二批：更贴实战的建议（重点）](#第二批更贴实战的建议重点)
- [优先级建议](#优先级建议)
- [实现建议](#实现建议)

---

## 先补现有工具的缺口

比新增命令更划算：先把已经在用的工具做「管线友好」。

### A. 全命令支持 Quiet / Pipe 模式

现状：几乎都强制进 termui，无法 `pf_tools pf_md5 xxx | pbcopy`，脚本里也用不了。

```bash
pf_tools pf_md5 hello -q          # 只打印结果到 stdout
pf_tools pf_b64 hello -q | pbcopy
echo '{"a":1}' | pf_tools pf_json -q
```

**实战价值**：终端工具一旦能进管道，使用频率会立刻上去。

### B. 增强 `pf_t`（不要另起 `pf_unixtime`）

```bash
pf_tools pf_t                     # 秒级时间戳（现有）
pf_tools pf_t -m                  # 毫秒时间戳（日志/JS 常用）
pf_tools pf_t 1699123456          # 戳 → 本地时间
pf_tools pf_t 1699123456000 -m    # 毫秒戳 → 时间
pf_tools pf_t "2023-11-05 10:30:00" -z Asia/Shanghai
pf_tools pf_t -f "2006-01-02"
```

### C. 增强 `pf_b64` / `pf_md5`

```bash
pf_tools pf_b64 hello -u          # Base64URL（JWT / URL 安全）
pf_tools pf_md5 -f ./file.bin     # 文件 MD5（发版校验）
pf_tools pf_md5 hello -a sha256   # 顺手支持 SHA（或拆 pf_sha）
```

### D. 增强 `pf_cwx`

```bash
pf_tools pf_cwx url1 url2 url3    # 批量检测投放链接
pf_tools pf_cwx -f urls.txt
```

---

## 🔐 编码/加密工具

### 1. **pf_json** - JSON 格式化 / 校验
```bash
pf_tools pf_json '{"name":"test","age":18}'
pf_tools pf_json file.json -f
pf_tools pf_json '{"a":1}' -c     # 压缩
pf_tools pf_json '{"a":1}' -v     # 仅校验
pf_tools pf_json resp.json -p data.items.0.id   # 取路径（轻量 jq）
```
**实战场景**: 抓包结果整理、配置校验、从响应里抠字段

### 2. **pf_url** - URL 编解码 / 拆解
```bash
pf_tools pf_url "https://example.com?name=张三"
pf_tools pf_url "https://a.com?x=%E5%BC%A0" -d
pf_tools pf_url "https://a.com/path?a=1&b=2" --parse   # scheme/host/path/query 拆开
```
**实战场景**: 回调参数乱码、Query 拼接排查

### 3. **pf_uuid** - UUID 生成器
```bash
pf_tools pf_uuid
pf_tools pf_uuid -c 5
pf_tools pf_uuid -n               # 无横杠（部分业务主键习惯）
```
**实战场景**: 测试数据、幂等 key、trace id

### 4. **pf_sha** - SHA / 文件哈希
```bash
pf_tools pf_sha "hello"
pf_tools pf_sha "hello" -a sha1
pf_tools pf_sha file.txt -f
```
**实战场景**: 发版包校验、缓存 key、签名对比

### 5. **pf_pass** - 随机密码 / Token
```bash
pf_tools pf_pass
pf_tools pf_pass -l 32 -s
pf_tools pf_pass -t hex -l 64     # API secret / webhook token
```
**实战场景**: 本地 `.env`、临时密钥、压测账号

---

## 🌐 网络工具

### 6. **pf_ip** - IP 地址查询
```bash
pf_tools pf_ip
pf_tools pf_ip 8.8.8.8
pf_tools pf_ip -l
```
**实战场景**: 部署核对出口 IP、CDN/回源排查

### 7. **pf_port** - 端口占用 / 连通
```bash
pf_tools pf_port 8080
pf_tools pf_port 80 -h example.com
pf_tools pf_port 8000-8010
```
**实战场景**: 本地服务起不来、端口冲突

### 8. **pf_http**（建议别叫 pf_curl）- 轻量请求
> 不要做成完整 curl。只做「看状态码 + 耗时 + 响应头 + JSON 美化」。

```bash
pf_tools pf_http https://api.example.com
pf_tools pf_http https://api.example.com -X POST -d '{"a":1}' -j
pf_tools pf_http https://api.example.com -H "Authorization: Bearer xxx" --time
```
**实战场景**: 快速探活、看延迟、看响应头

### 9. **pf_dns** - DNS 查询
```bash
pf_tools pf_dns example.com
pf_tools pf_dns example.com -t MX
```
**实战场景**: 域名刚解析完验证、邮件 MX 排查

---

## 📁 文件操作工具

### 10. **pf_tree** - 目录树（比 find 更贴 CLI 习惯）
```bash
pf_tools pf_tree .
pf_tools pf_tree . -d 2 -e go,md
```
**实战场景**: 快速摸生疏仓库结构

### 11. **pf_size** - 体积排序
```bash
pf_tools pf_size .
pf_tools pf_size . -s -n 20       # 最大 20 项
```
**实战场景**: `node_modules` / 日志撑爆磁盘

### 12. **pf_lines** - 代码行数（原 pf_count）
```bash
pf_tools pf_lines .
pf_tools pf_lines . -e go,ts
```
**实战场景**: 周报、重构前后对比

---

## 🛠️ 开发辅助工具

### 13. **pf_env** - `.env` / 环境变量
```bash
pf_tools pf_env PATH
pf_tools pf_env -f .env           # 校验重复 key、空值、导出为 export 语句
```
**实战场景**: 本地配置排查、把 `.env` 转 shell

### 14. **pf_case** - 命名风格转换
```bash
pf_tools pf_case "user_id" --camel
pf_tools pf_case "UserId" --snake
pf_tools pf_case "user-id" --export
```
**实战场景**: Go/JSON/DB 字段命名来回转

### 15. **pf_regex** - 正则试跑
```bash
pf_tools pf_regex "test@example.com" `^[\w.-]+@[\w.-]+\.\w+$`
pf_tools pf_regex -f access.log "status=(\d+)" --group 1
```
**实战场景**: 写规则前先本地验证

### 16. **pf_qr** - 终端二维码
```bash
pf_tools pf_qr "https://example.com"
pf_tools pf_qr "https://example.com" -o qr.png
```
**实战场景**: 手机扫调试页、分享内网地址

---

## 💻 系统工具

### 17. **pf_killport** - 按端口杀进程（比通用 pf_proc 更刚需）
```bash
pf_tools pf_killport 8080
pf_tools pf_killport 8080 -f      # SIGKILL
```
**实战场景**: 旧服务占端口，本地重启最烦的一步

### 18. **pf_whichport** - 谁占用了端口
```bash
pf_tools pf_whichport 3306
```
**实战场景**: 查清是 Docker 还是本地 mysqld

---

## 🔄 数据转换工具

### 19. **pf_hex** - 十六进制
```bash
pf_tools pf_hex "hello"
pf_tools pf_hex "68656c6c6f" -d
```

### 20. **pf_html** - HTML 实体编解码
```bash
pf_tools pf_html "<script>" 
pf_tools pf_html "&lt;script&gt;" -d
```

---

## 第二批：更贴实战的建议（重点）

下面这批比「再封装一个 find/git」更符合本仓库气质：**调试瞬间高频、实现轻、和现有 MD5/B64/时间戳同一赛道**。

### 21. **pf_jwt** - JWT 解码（不校验签名）⭐
```bash
pf_tools pf_jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
pf_tools pf_jwt $TOKEN --exp      # 是否过期、剩余秒数
```
**为什么刚需**: 联调登录态时，每天都要看 payload / exp，不必打开 jwt.io。  
**实现**: 标准库 `encoding/base64` + `encoding/json`，零依赖。

### 22. **pf_cert** - 域名证书信息 ⭐
```bash
pf_tools pf_cert example.com
pf_tools pf_cert example.com:443 --days   # 还有多少天过期
```
**为什么刚需**: HTTPS 挂了、证书快过期、链不完整，运维/开发都常查。  
**实现**: `crypto/tls` 拨号读 `PeerCertificates`。

### 23. **pf_cron** - Cron 表达式解读
```bash
pf_tools pf_cron "0 */5 * * *"
pf_tools pf_cron "0 9 * * 1-5" -n 5     # 接下来 5 次触发时间
```
**为什么刚需**: 写定时任务时最容易写错；看别人 crontab 也费劲。

### 24. **pf_mock** - 快速造测试数据
```bash
pf_tools pf_mock user -c 5
pf_tools pf_mock phone -c 10          # 可复用 phone.dat 思路
pf_tools pf_mock json --fields id:uuid,name:zh,email,created_at:ts
```
**为什么刚需**: 写接口 / 前端联调时，临时 JSON 比手敲快。  
**贴合项目**: 已有 `phone.dat` + 手机号能力，可自然扩展。

### 25. **pf_rand** - 随机串（区别于可读密码）
```bash
pf_tools pf_rand 32
pf_tools pf_rand 16 -t digit
pf_tools pf_rand 24 -t base58
```
**为什么刚需**: order_no、nonce、state、CSRF token，和「人类密码」不是一类需求。

### 26. **pf_header** - 解析裸 HTTP 头 / curl -v 片段
```bash
pf_tools pf_header $'Content-Type: application/json\nAuthorization: Bearer xxx'
pbpaste | pf_tools pf_header
```
**为什么刚需**: 从抓包/日志里粘一段头出来，快速看 Host、Cookie、鉴权字段。

### 27. **pf_cookie** - Cookie 串拆解
```bash
pf_tools pf_cookie "a=1; b=2; session=xxx"
```
**为什么刚需**: 浏览器 Application 面板复制出来一长串，拆开比肉眼扫快。

### 28. **pf_status** - HTTP 状态码速查
```bash
pf_tools pf_status 418
pf_tools pf_status 502
```
**为什么刚需**: `451`/`422`/`429` 不常用时容易忘语义；实现就是一张表。

### 29. **pf_cidr** - IP / CIDR 判断
```bash
pf_tools pf_cidr 10.0.0.8 10.0.0.0/24
pf_tools pf_cidr 10.0.0.0/24 --list -n 5
```
**为什么刚需**: 网关白名单、内网判断、安全组规则核对。

### 30. **pf_semver** - 版本号 bump
```bash
pf_tools pf_semver 1.2.3 --patch
pf_tools pf_semver 1.2.3 --minor
pf_tools pf_semver 1.2.3-beta.1 --inc-prerelease
```
**为什么刚需**: 发版改 tag / CHANGELOG 前快速算下一版。

### 31. **pf_yaml** / **pf_toml** - 配置互转
```bash
pf_tools pf_yaml config.yaml --json
pf_tools pf_json config.json --yaml
```
**为什么刚需**: K8s / CI / 应用配置格式来回切，调试时高频。

### 32. **pf_sqlstr** - SQL 字符串转义 / 反转义
```bash
pf_tools pf_sqlstr "O'Reilly"
pf_tools pf_sqlstr "'O\'Reilly'" -d
```
**为什么刚需**: 拼测试 SQL、看日志里的转义串（不做 SQL 注入工具）。

### 33. **pf_unicode** - Unicode / 中文码点
```bash
pf_tools pf_unicode "你好"
pf_tools pf_unicode "U+4F60 U+597D" -d
pf_tools pf_unicode "\\u4f60\\u597d" -d
```
**为什么刚需**: 日志里一堆 `\uXXXX`、特殊空白字符排查。

### 34. **pf_imgb64** - 图片 ↔ Base64 Data URI
```bash
pf_tools pf_imgb64 ./logo.png
pf_tools pf_imgb64 "data:image/png;base64,..." -o out.png
```
**为什么刚需**: 写邮件模板、Markdown、小图标内嵌时常用。

### 35. **pf_expand** - 短链 / 重定向追踪
```bash
pf_tools pf_expand https://t.cn/xxxxx
pf_tools pf_expand https://bit.ly/xxx --max 10
```
**为什么刚需**: 运营短链、短信链接最终落到哪，和现有 `pf_cwx` 同一业务域。

### 36. **pf_whois**（可选）/ **pf_ttl**
```bash
pf_tools pf_ttl example.com         # DNS TTL + 当前解析
```
**为什么刚需**: 改完解析「怎么还不生效」时看 TTL 最直接。

### 37. **pf_diffj** - JSON 结构 diff
```bash
pf_tools pf_diffj old.json new.json
```
**为什么刚需**: 接口字段增减对比，比文本 diff 更适合 JSON。

### 38. **pf_luhn** / **pf_idcard** - 校验位（国内场景）
```bash
pf_tools pf_idcard 110101199003078 agx   # 校验格式/校验位（测试号）
pf_tools pf_luhn 622202xxxxxxxxxxxx
```
**为什么刚需**: 和 `pf_m` 一样偏国内业务联调；**仅校验算法，不采集真实隐私数据**。

### 39. **pf_mime** - 扩展名 / Content-Type
```bash
pf_tools pf_mime .wasm
pf_tools pf_mime file.unknown --sniff
```
**为什么刚需**: 上传接口、静态资源 Content-Type 设错时秒查。

### 40. **pf_clip** - 剪贴板编解码流水线
```bash
pf_tools pf_clip b64              # 读剪贴板 → base64 → 写回
pf_tools pf_clip jwt              # 读剪贴板 JWT 并解码展示
pf_tools pf_clip md5
```
**为什么刚需**: 本项目已经重度依赖 clipboard；做成「剪贴板进、结果出」会非常顺手。

---

## 🎯 优先级建议

### P0（最贴本仓库、实现便宜、日用频率高）
1. **全命令 `-q` / pipe 模式**（基础设施，先做）
2. **增强 `pf_t`**（毫秒、双向转换、时区）
3. **`pf_json`**
4. **`pf_url`（含 --parse）**
5. **`pf_jwt`**
6. **`pf_uuid` / `pf_rand`**

### P1（联调 / 排障高频）
7. **`pf_killport` + `pf_whichport`**
8. **`pf_cert`**
9. **`pf_http`（轻量，不替代 curl）**
10. **`pf_ip` / `pf_dns` / `pf_ttl`**
11. **`pf_sha` 或扩展 `pf_md5 -f/-a`**
12. **`pf_pass`**

### P2（提升质感 / 垂直场景）
13. **`pf_cron`**
14. **`pf_mock`（承接 phone.dat）**
15. **`pf_cookie` / `pf_header`**
16. **`pf_case` / `pf_unicode`**
17. **`pf_expand`（搭配 pf_cwx）**
18. **`pf_clip`**

### 不建议优先做（容易做成劣质轮子）
- 完整 `pf_git` / `pf_find` / `pf_curl` 全家桶
- 系统监控大盘（`pf_sys`/`pf_disk` 不如用系统自带）
- 重型扫描器、爬虫框架

---

## 💡 实现建议

### 代码结构（保持现状）
1. 命令入口: `cmd/pfJson.go`、`cmd/pfJwt.go` …
2. 逻辑实现: `pak/json_format.go`、`pak/jwt_decode.go` …
3. UI: 继续 termui；**同时**支持 `-q` 直接 `fmt.Println`
4. 复制: 保留 `c` 键写剪贴板；`-q` 模式交给管道

### UX 统一约定
```text
成功: 绿色边框 / 绿色字
失败: 红色提示
按 q 退出
按 c 复制「主结果」（hash、token、时间戳等）
-q / --quiet: 无 UI，stdout 纯文本，便于脚本
```

### 技术栈（尽量标准库）
| 能力 | 建议 |
|------|------|
| JSON | `encoding/json` |
| URL | `net/url` |
| UUID | `github.com/google/uuid` 或自研 v4 |
| JWT 解码 | `encoding/base64` + `encoding/json` |
| 证书 | `crypto/tls` |
| HTTP | `net/http` |
| 哈希 | `crypto/sha256` 等 |
| YAML | `gopkg.in/yaml.v3`（若做互转再引入） |

### 和现有工具的关系（避免重复造名）

| 需求 | 建议 |
|------|------|
| 时间戳 | **扩展 `pf_t`**，不要新建 `pf_unixtime` |
| Base64URL | **扩展 `pf_b64 -u`** |
| 文件哈希 | **扩展 `pf_md5 -f` 或加 `pf_sha`** |
| 微信链路 | **`pf_cwx` 批量 + 新 `pf_expand`** |
| 手机测试号 | **扩展数据能力到 `pf_mock`** |

---

## 📝 现有 vs 建议增强

| 能力域 | 现有 | 建议 |
|--------|------|------|
| 摘要 | MD5 | + 文件哈希、SHA、`-q` |
| 编码 | Base64 | + URL / Base64URL / Hex / HTML 实体 |
| 时间 | 时间→戳 | + 毫秒、戳→时间、时区、格式 |
| 网络 | 微信 URL 检测 | + IP/DNS/证书/短链展开/轻量 HTTP |
| 身份调试 | 无 | + JWT 解码、Cookie/Header 解析 |
| 本地开发 | 清 DS_Store、建日期目录 | + killport、`.env` 检查、mock 数据 |
| 交互 | 仅 TUI | + 管道模式，保留一键复制 |

---

## 🚀 两周落地路线（务实版）

**Week 1 — 基础设施 + 三板斧**
1. 抽出公共 `pak/ui_result.go`：`Show(title, text)` + `ShowTable` + `QuietPrintln`
2. 给 `pf_md5` / `pf_b64` / `pf_t` 加上 `-q`
3. 落地 `pf_json`、`pf_url`、`pf_jwt`

**Week 2 — 联调排障**
4. 增强 `pf_t`（毫秒/双向）
5. `pf_uuid` + `pf_rand`/`pf_pass`
6. `pf_killport` + `pf_cert`

这 8 项做完，工具集会从「小玩具集合」变成「开发者口袋刀」，且不偏离现有架构与交互习惯。
