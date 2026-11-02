# app team-app add

向团队应用库添加应用版本

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli app team-app add` |
| HTTP | `POST /open/api/app/teamApp/create` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--apk-icon-url` | string | 是 |  | apk图标url |
| `--apk-name` | string | 是 |  | apk名称 |
| `--apk-pkg` | string | 是 |  | apk包名 |
| `--apk-version-code` | integer | 是 |  | apk版本号 |
| `--apk-version-name` | string | 是 |  | apk版本名称 |
| `--md5` | string | 是 |  | 文件 MD5 |
| `--name` | string | 是 |  | 文件名称（原始文件名称md5值+文件后缀） |
| `--origin-name` | string | 是 |  | 原始文件名称 |
| `--size` | integer | 是 |  | 文件大小（单位：B） |
| `--url` | string | 是 |  | 文件url |
| `--cate-code` | integer |  | `-1` | 所属分类编码：-1代表未分类 |
| `--shelf-status` | integer |  | `1` | 应用状态：1-上架，0=下架 |
| `--weight` | integer |  | `0` | 权重，越大越靠前，默认：0 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli app team-app add --apk-icon-url <string> --apk-name <string> --apk-pkg <string> --apk-version-code <int> --apk-version-name <string> --md5 <string> --name <string> --origin-name <string> --size <int> --url <string>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RLong`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | integer | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-app](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)
