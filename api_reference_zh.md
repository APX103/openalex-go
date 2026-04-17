# openalex-go API 参考

`go get github.com/APX103/openalex-go`

## 初始化

```go
c := openalex.New(opts...)
```

| Option | 类型 | 说明 |
|--------|------|------|
| `WithAPIKey` | `string` | API Key，优先级最高。有 Key 时速率 100 次/秒，无 Key 约 10 次/秒 |
| `WithMailto` | `string` | 邮箱，无 API Key 时进入礼貌池。有 Key 时此字段不生效 |
| `WithTimeout` | `time.Duration` | 请求超时，默认 15s |
| `WithBaseURL` | `string` | 自定义 API 地址，默认 `https://api.openalex.org` |
| `WithHTTPClient` | `*http.Client` | 完全替换默认 HTTP 客户端 |

认证优先级：`apiKey` > `mailto` > 无认证（共享池 100 credits/天）。

---

## 通用类型

### PageParams — 分页

```go
type PageParams struct {
    Page    int // 页码，默认 1
    PerPage int // 每页条数，默认 25，API 硬上限 200
}
```

零值 `PageParams{}` 等于第 1 页、每页 25 条。`PerPage` 超过 200 会被自动钳制。

### SortOption — 排序

```go
type SortOption struct {
    Field string // 字段名
    Order string // "desc"（默认）或 "asc"
}
```

常用 Field：`relevance_score`（仅搜索有效）、`cited_by_count`、`publication_date`、`created_date`。

### ListResponse[T] — 列表响应

```go
type ListResponse[T any] struct {
    Meta    Meta      // Meta.Count = 结果总数（所有页的总量，不只是当前页）
    Results []T       // 当前页数据
    GroupBy []GroupBy // 仅 group-by 聚合查询时有值，普通搜索/分页查询为空
}
```

### APIError

```go
type APIError struct {
    StatusCode int    // HTTP 状态码：404 未找到、429 速率限制、403 认证失败等
    Message    string // API 返回的错误信息
    URL        string // 请求的完整 URL
}
```

非 200 响应返回 `*APIError`，可通过 `errors.As(err, &apiErr)` 提取。网络和 JSON 解析错误返回 wrapped error。

---

## work 包 — 论文

### work.Search

```go
func work.Search(ctx context.Context, c *openalex.Client, params work.SearchParams) (*openalex.ListResponse[work.Work], error)
```

**SearchParams**

| 字段 | 类型 | 说明 |
|------|------|------|
| Query | `string` | 搜索关键词，全文检索标题和摘要 |
| Page | `int` | 页码，默认 1 |
| PerPage | `int` | 每页条数，默认 25，上限 200 |
| Sort | `*openalex.SortOption` | 排序。`nil` 则按相关性降序 |
| Select | `[]string` | 返回字段白名单，用逗号拼成 `select` 参数，减少响应体积 |
| Filters | `map[string]string` | 过滤条件，自动拼接为 `key:value,key2:value2` |
| GroupBy | `string` | 聚合字段，用于 `work.GroupBy`，如 `"type"`、`"publication_year"` |

```go
resp, err := work.Search(ctx, c, work.SearchParams{
    Query:   "large language model",
    PerPage: 20,
    Sort:    &openalex.SortOption{Field: "cited_by_count"},
    Select:  []string{"id", "display_name", "cited_by_count"},
    Filters: map[string]string{"publication_year": "2024", "type": "article"},
})
```

**Filters 常用键**

| 键 | 值格式 | 示例 | 含义 |
|----|--------|------|------|
| `publication_year` | `"YYYY"` 或 `"YYYY-YYYY"` | `"2024"` / `"2020-2024"` | 发表年份，支持范围 |
| `type` | 枚举值 | `"article"` / `"preprint"` | 文献类型，见 Work.Type 取值 |
| `primary_location.source.id` | Source 短 ID | `"S137773608"` | 发表期刊 |
| `author.id` | Author 短 ID | `"A5023898321"` | 作者 |
| `authorships.institutions.id` | Institution 短 ID | `"I136199984"` | 作者所属机构 |
| `has_doi` | `"true"` / `"false"` | `"true"` | 是否有 DOI |
| `is_oa` | `"true"` / `"false"` | `"true"` | 是否开放获取 |
| `cited_by_count` | 数字或范围或比较 | `">100"` / `"50-200"` | 被引次数 |
| `concepts.id` | Concept 短 ID | `"C154945302"` | 关联概念 |
| `open_access.oa_status` | OA 枚举值 | `"gold"` / `"diamond"` | OA 类型，见 OpenAccess.OaStatus 取值 |
| `default.search_filter` | 枚举值 | `"journal"` | 限定搜索范围为期刊论文 |

---

### work.GroupBy

```go
func work.GroupBy(ctx context.Context, c *openalex.Client, params work.SearchParams) ([]openalex.GroupBy, error)
```

返回匹配查询和过滤条件的论文的聚合桶。`params.GroupBy` 指定聚合字段（如 `"type"`、`"publication_year"`、`"primary_topic.field.id"`）。

```go
buckets, err := work.GroupBy(ctx, c, work.SearchParams{
    Filters: map[string]string{"publication_year": "2024"},
    GroupBy: "type",
})
// buckets: [{Key: "article", KeyDisplayName: "Article", Count: 1234}, ...]
```

**返回值**：`[]openalex.GroupBy`。

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `Key` | `string` | `key` | 聚合键值 |
| `KeyDisplayName` | `string` | `key_display_name` | 聚合键的可读名称 |
| `Count` | `int` | `count` | 该桶的计数 |

---

### work.Get

```go
func work.Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*work.Work, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| id | `string` | 论文 ID。接受短 ID（`"W2626778328"`）或完整 URL（`"https://openalex.org/W2626778328"`） |
| selectFields | `...string` | 可选。指定返回字段，不传则返回全部字段 |

**返回值**：`*work.Work`。未找到时返回 `fmt.Errorf("work %s not found", id)`。

---

### work.GetByIDs

```go
func work.GetByIDs(ctx context.Context, c *openalex.Client, ids []string, selectFields ...string) ([]work.Work, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| ids | `[]string` | 论文短 ID 列表。单次上限 **200 个**，超出直接返回错误 |
| selectFields | `...string` | 可选。指定返回字段 |

**返回值**：`[]work.Work`。

---

### work.GetCitedBy

```go
func work.GetCitedBy(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回**引用了** `workID` 的论文列表。固定排序 `cited_by_count:desc`（被引最多的排在前面）。

---

### work.GetReferencedWorks

```go
func work.GetReferencedWorks(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回 `workID` **引用的**参考文献列表。固定排序 `cited_by_count:desc`。

---

### work.GetRelated

```go
func work.GetRelated(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回与 `workID` 相关的论文。基于 N-gram 算法计算概念重合度。固定排序 `cited_by_count:desc`。

---

### work.GetByAuthor

```go
func work.GetByAuthor(ctx context.Context, c *openalex.Client, authorID string, page openalex.PageParams, extraFilter string, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回指定作者的全部论文。固定排序 `cited_by_count:desc`。

`extraFilter` 为空字符串时不追加额外过滤；非空时会追加到 filter 查询中（如 `"concepts.id:C154945302"`）。

---

### work.GetBySource

```go
func work.GetBySource(ctx context.Context, c *openalex.Client, sourceID string, page openalex.PageParams, sort *openalex.SortOption, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回指定期刊/源的全部论文。`sort` 为 nil 时默认 `cited_by_count:desc`，可覆盖为其他排序（如 `publication_date`）。

---

### work.Work

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | 论文 OpenAlex ID，完整 URL，如 `"https://openalex.org/W2626778328"`。可用 `util.ShortID()` 提取短 ID `"W2626778328"` |
| `Doi` | `string` | `doi` | DOI，完整 URL，如 `"https://doi.org/10.7717/peerj.4375"`。纯 DOI 需 `strings.TrimPrefix(w.Doi, "https://doi.org/")` |
| `DisplayName` | `string` | `display_name` | 论文标题 |
| `PubYear` | `int` | `publication_year` | 发表年份 |
| `PubDate` | `string` | `publication_date` | 发表日期，ISO 8601 格式，如 `"2018-02-13"`。有多个日期时取最早的电子发表日期 |
| `Type` | `string` | `type` | 文献类型。常见值：`"article"`（期刊/会议论文，包含最多）、`"preprint"`（预印本）、`"paratext"`（非学术内容如封面目录）、`"book-chapter"`、`"dissertation"`、`"dataset"`、`"erratum"`（勘误）、`"editorial"`（社论）、`"letter"`、`"review"`（综述期刊论文）、`"grant"`、`"peer-review"` 等 |
| `Language` | `string` | `language` | 论文语言，ISO 639-1 代码，如 `"en"`、`"zh"`、`"fr"`。基于摘要/标题自动检测，可能为空 |
| `IndexedIn` | `[]string` | `indexed_in` | 被哪些数据库收录。取值：`"arxiv"`、`"crossref"`、`"doaj"`、`"pubmed"` |
| `OpenAccess` | `OpenAccess` | `open_access` | 开放获取信息，见下表 |
| `Authorships` | `[]Authorship` | `authorships` | 作者列表，最多 100 位。包含每位作者的身份、机构和通讯作者标记 |
| `PrimaryLoc` | `*PrimaryLocation` | `primary_location` | 主要发表位置（最接近正式版本的副本所在）。**可能为 nil** |
| `BestOALoc` | `*PrimaryLocation` | `best_oa_location` | 最佳 OA 位置（综合评分最高的开放获取副本）。**可能为 nil**。评分依据：publisher > repository，publishedVersion > acceptedVersion > submittedVersion，有 PDF 链接优先 |
| `Topics` | `[]WorkTopic` | `topics` | 关联主题，最多 3 个。每个 Topic 包含学科层级（子领域→领域→学科门类），见 WorkTopic |
| `Concepts` | `[]Concept` | `concepts` | 关联概念标签。`Score >= 0.3` 表示强关联，低分可能是祖先概念的继承。见 Concept |
| `Keywords` | `[]Keyword` | `keywords` | 基于 Topic 提取的关键词短语 |
| `Refs` | `[]string` | `referenced_works` | 该论文引用的文献 ID 列表（此文献 → 其他文献） |
| `Related` | `[]string` | `related_works` | 相关论文 ID 列表，基于概念重合度计算 |
| `CountsByYear` | `[]CountByYear` | `counts_by_year` | 近十年逐年被引次数。**不保证按年份排序**，使用前需排序。零引用的年份可能缺失 |
| `CitedByCount` | `int` | `cited_by_count` | 总被引次数（其他论文引用此论文的次数） |
| `AbstractInv` | `map[string][]int` | `abstract_inverted_index` | 摘要倒排索引（法律原因不存储明文）。键为单词，值为该单词出现的位置编号。**必须调用 `util.RestoreAbstract()` 还原为可读文本** |
| `Biblio` | `*Biblio` | `biblio` | 书目信息：卷号、期号、页码。**可能为 nil** |
| `IDs` | `WorkIDs` | `ids` | 多平台外部 ID，见 WorkIDs |

### work.OpenAccess

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `IsOA` | `bool` | `is_oa` | 是否开放获取（有免费可读的全文链接） |
| `OaStatus` | `string` | `oa_status` | OA 状态：`"gold"`（发表在完全 OA 期刊）、`"green"`（存放在机构/学科仓库）、`"hybrid"`（混合型期刊的 OA 文章）、`"bronze"`（出版商页面免费但无明确许可证）、`"diamond"`（完全 OA 期刊且免 APC）、`"closed"`（非 OA） |
| `OaURL` | `*string` | `oa_url` | 最佳 OA 链接（最接近正式版本的免费全文 URL）。**可能为 nil**。可能是 PDF 直链或落地页 |
| `AnyRepoHasFulltext` | `bool` | `any_repository_has_fulltext` | 是否有任何仓库托管了全文（即使 oa_status 不是 green，也可能存在"影子绿色 OA"） |

### work.Authorship

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `AuthorPosition` | `string` | `author_position` | 作者序位：`"first"`（第一作者）、`"middle"`（中间作者）、`"last"`（末位作者/通讯作者） |
| `Author` | `AuthorRef` | `author` | 作者基本信息（姓名、OpenAlex ID、ORCID） |
| `Institutions` | `[]Institution` | `institutions` | 该作者在此论文中的机构列表 |
| `Countries` | `[]string` | `countries` | 机构所在国家，ISO 二字母代码如 `"US"`、`"CN"` |
| `IsCorresponding` | `bool` | `is_corresponding` | 是否为通讯作者 |

### work.PrimaryLocation

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `IsOA` | `bool` | `is_oa` | 此位置的全文是否免费 |
| `LandingPageURL` | `string` | `landing_page_url` | 落地页 URL（通常是 DOI 链接） |
| `PdfURL` | `*string` | `pdf_url` | 直接 PDF 下载链接。**可能为 nil** |
| `Source` | `*LocationSource` | `source` | 发表源（期刊或仓库）。**可能为 nil** |
| `License` | `string` | `license` | 开放许可证，如 `"cc-by"`、`"cc-by-nc"`、`"cc0"` 等。空字符串表示无明确许可证 |
| `Version` | `string` | `version` | 全文版本：`"publishedVersion"`（正式发表版）、`"acceptedVersion"`（同行评审后接受版）、`"submittedVersion"`（投稿版） |

### work.LocationSource

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | 来源的 OpenAlex ID |
| `DisplayName` | `string` | `display_name` | 来源名称，如期刊名或仓库名 |
| `ISSN` | `any` | `issn` | ISSN，可能是 `string`、`[]string` 或 `nil`（API 返回类型不固定） |
| `Type` | `string` | `type` | 来源类型：`"journal"`、`"repository"` 等，见 Source.Type |
| `IsOA` | `bool` | `is_oa` | 该来源是否为 OA |

### work.WorkTopic

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | Topic 的 OpenAlex ID |
| `DisplayName` | `string` | `display_name` | Topic 名称，如 `"Large Language Models"` |
| `Count` | `int` | `count` | 该 Topic 关联的论文总数 |
| `Subfield` | `TopicRef` | `subfield` | 子领域，如 `"Artificial Intelligence"` |
| `Field` | `TopicRef` | `field` | 领域，如 `"Computer Science"` |
| `Domain` | `TopicRef` | `domain` | 学科门类，如 `"Engineering"` |

### work.WorkIDs

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `OpenAlex` | `string` | `openalex` | OpenAlex ID，完整 URL |
| `Doi` | `string` | `doi` | DOI，完整 URL 如 `"https://doi.org/10.7717/peerj.4375"` |
| `Mag` | `string` | `mag` | Microsoft Academic Graph ID（纯数字字符串）。MAG 已停止更新 |
| `PMID` | `string` | `pmid` | PubMed ID，完整 URL 如 `"https://pubmed.ncbi.nlm.nih.gov/29456894"`。**可能为空** |
| `ArXiv` | `string` | `arxiv` | arXiv ID，完整 URL 如 `"https://arxiv.org/abs/2101.00001"`。**可能为空** |

---

## author 包 — 作者

### author.Search

```go
func author.Search(ctx context.Context, c *openalex.Client, params author.SearchParams) (*openalex.ListResponse[author.Author], error)
```

固定排序 `relevance_score:desc`。参数同 work.Search 但无 Filters 和 Sort。

```go
resp, err := author.Search(ctx, c, author.SearchParams{Query: "Andrew Ng", PerPage: 10})
```

### author.Get

```go
func author.Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*author.Author, error)
```

未找到返回 `fmt.Errorf("author %s not found", id)`。

### author.Author

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | OpenAlex ID，完整 URL |
| `DisplayName` | `string` | `display_name` | 作者姓名 |
| `Orcid` | `string` | `orcid` | ORCID，完整 URL 如 `"https://orcid.org/0000-0001-6187-6610"`。**覆盖率较低，老作者尤其缺失** |
| `WorksCount` | `int` | `works_count` | 该作者的论文总数 |
| `CitedByCount` | `int` | `cited_by_count` | 该作者所有论文的总被引次数 |
| `SummaryStats` | `SummaryStats` | `summary_stats` | 引用指标，见 SummaryStats |
| `LastKnownInsts` | `[]Institution` | `last_known_institutions` | 最近已知机构。取该作者所有论文中按发表日期最新的一篇的机构信息。可能有多个（多机构署名） |
| `Topics` | `[]AuthorTopic` | `topics` | 研究主题。每个 Topic 含学科层级和关联度，见 AuthorTopic |
| `XConcepts` | `[]Concept` | `x_concepts` | 关联概念（即将废弃，建议用 Topics 替代）。Score 范围 0-100，表示该概念在作者论文中的出现强度 |
| `CountsByYear` | `[]CountByYear` | `counts_by_year` | 近十年逐年统计。**不保证有序**。`WorksCount` 为当年发表论文数，`CitedByCount` 为当年被引次数。零值年份可能缺失 |
| `WorksAPIURL` | `string` | `works_api_url` | 获取该作者全部论文的 API URL，如 `"https://api.openalex.org/works?filter=author.id:A5023898321"` |

---

## source 包 — 期刊

### source.Search

```go
func source.Search(ctx context.Context, c *openalex.Client, params source.SearchParams) (*openalex.ListResponse[source.Source], error)
```

固定排序 `relevance_score:desc`。参数同 author.Search。

### source.Get

```go
func source.Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*source.Source, error)
```

未找到返回 `fmt.Errorf("source %s not found", id)`。

### source.Source

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | OpenAlex ID，完整 URL |
| `DisplayName` | `string` | `display_name` | 期刊/源名称 |
| `ISSN` | `[]string` | `issn` | ISSN 列表（不同媒介版本有不同 ISSN） |
| `ISSNL` | `string` | `issn_l` | ISSN-L（Linking ISSN），唯一标识该刊。通常是印刷版 ISSN |
| `IsOA` | `bool` | `is_oa` | 是否为完全 OA 期刊/仓库。注意：状态可能随时间变化，老文章可能不免费 |
| `Type` | `string` | `type` | 类型：`"journal"`（期刊）、`"repository"`（仓库，如 arXiv、PubMed Central）、`"conference"`（会议）、`"ebook platform"`、`"book series"`、`"metadata"`、`"other"` |
| `WorksCount` | `int` | `works_count` | 该期刊收录的论文总数 |
| `CitedByCount` | `int` | `cited_by_count` | 该期刊所有论文的总被引次数 |
| `SummaryStats` | `SummaryStats` | `summary_stats` | 引用指标，见 SummaryStats |
| `HomepageURL` | `*string` | `homepage_url` | 期刊官网首页 URL。**可能为 nil** |
| `HostOrgName` | `*string` | `host_organization_name` | 出版机构名称，如 `"Elsevier BV"`、`"Springer Nature"`。**可能为 nil** |
| `APCUSD` | `*float64` | `apc_usd` | 文章处理费（美元），数据来自 DOAJ。**可能为 nil** |
| `CountryCode` | `string` | `country_code` | 出版国家，ISO 二字母代码如 `"US"`、`"GB"`、`"CN"` |
| `Topics` | `[]AuthorTopic` | `topics` | 该期刊主要发表的主题 |
| `TopicShare` | `[]AuthorTopic` | `topic_share` | 各主题在期刊中的占比。`Value` 字段表示该主题论文占比（0-1），乘 100 得百分比 |
| `CountsByYear` | `[]CountByYear` | `counts_by_year` | 近十年逐年统计。`WorksCount` 为当年新收录论文数，`CitedByCount` 为当年该期刊论文被引次数。**不保证有序** |
| `WorksAPIURL` | `string` | `works_api_url` | 获取该期刊全部论文的 API URL |

---

## util 包 — 工具函数

### util.ShortID

```go
func util.ShortID(openalexURL string) string
```

从完整 OpenAlex URL 中提取短 ID。所有 API 返回的 ID 都是完整 URL，展示或拼接 filter 时通常需要短 ID。

```
"https://openalex.org/W2626778328" → "W2626778328"
```

### util.JoinPipe

```go
func util.JoinPipe(ids []string) string
```

用 `|` 连接多个 ID，用于 OpenAlex 的 filter 批量查询。

```
[]string{"W1", "W2", "W3"} → "W1|W2|W3"
```

### util.RestoreAbstract

```go
func util.RestoreAbstract(idx map[string][]int) string
```

将 OpenAlex 倒排索引 `{单词: [位置1, 位置2, ...]}` 还原为按位置排列的纯文本。`nil` 返回空字符串。

```go
text := util.RestoreAbstract(w.AbstractInv)
// "Despite growing interest in Open Access, ..."
```

### util.ResolvePDF

```go
func util.ResolvePDF(w util.PDFWork) util.PDFResult
```

按优先级解析论文 PDF 链接：

1. **arXiv**：如有 arXiv ID → `https://arxiv.org/pdf/{arxivID}`（最可靠）
2. **OpenAlex best_oa_location.pdf_url**：最佳 OA 位置的 PDF 直链
3. **OpenAlex open_access.oa_url**：OA 链接（可能是落地页而非 PDF）
4. **DOI**：`https://doi.org/{doi}`（跳转到出版商页面，**不一定是 PDF**）
5. **无**：`PDFSourceNone`，URL 为空

```go
type PDFResult struct {
    URL    string     // PDF 地址，可能为空
    Source PDFSource  // 来源类型
}

pdf := util.ResolvePDF(&w)
if pdf.URL != "" {
    fmt.Println(pdf.URL, util.PDFSourceName(pdf.Source))
}
```

### util.PDFSource

| 常量 | 值 | `PDFSourceName()` 返回 | 说明 |
|------|----|----------------------|------|
| `PDFSourceArXiv` | 1 | `"arxiv"` | arXiv 直链，100% 可用 |
| `PDFSourceOpenAlex` | 2 | `"openalex"` | OpenAlex 记录的 OA 链接 |
| `PDFSourceUnpaywall` | 3 | `"unpaywall"` | 预留，未实现 |
| `PDFSourceDOI` | 4 | `"doi"` | DOI 跳转，通常到出版商页面而非 PDF |
| `PDFSourceNone` | 5 | `""` | 无可用 PDF |

---

## 共享类型

### SummaryStats — 引用指标（Author / Source）

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `HIndex` | `int` | `h_index` | h 指数：有 h 篇论文被引至少 h 次 |
| `I10Index` | `int` | `i10_index` | i10 指数：被引 10 次以上的论文数 |
| `TwoYrMeanCitedness` | `float64` | `2yr_mean_citedness` | 近两年平均被引次数，类似影响因子（Impact Factor）。用去年引用数除以前年和大前年发表的论文数 |

### AuthorTopic — 研究主题（Author.Topics / Source.Topics / Source.TopicShare）

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | Topic 的 OpenAlex ID |
| `DisplayName` | `string` | `display_name` | Topic 名称 |
| `Count` | `int` | `count` | 关联的论文总数 |
| `Value` | `float64` | `value` | 占比分数（0-1）。在 `Source.TopicShare` 中表示该主题论文占比；在 `Author.Topics` 中表示关联度 |
| `Subfield` | `TopicRef` | `subfield` | 子领域，如 `"Artificial Intelligence"` |
| `Field` | `TopicRef` | `field` | 领域，如 `"Computer Science"` |
| `Domain` | `TopicRef` | `domain` | 学科门类，如 `"Formal Sciences"` |

### CountByYear — 逐年统计（Work / Author / Source）

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `Year` | `int` | `year` | 年份 |
| `WorksCount` | `int` | `works_count` | 当年发表/收录的论文数 |
| `CitedByCount` | `int` | `cited_by_count` | 当年被引次数 |

**注意**：仅包含近十年数据，零值年份已移除，**不保证按年份排序**。

### Concept — 概念标签（Work.Concepts / Author.XConcepts）

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | 概念 OpenAlex ID |
| `DisplayName` | `string` | `display_name` | 概念名称，如 `"Machine Learning"`、`"COVID-19"` |
| `Score` | `float64` | `score` | 关联强度。Work 中 `>= 0.3` 为强关联；Author 中 `0-100` 表示出现频率 |

### TopicRef — 学科引用（WorkTopic / AuthorTopic 内嵌）

| 字段 | 类型 | JSON | 说明 |
|------|------|------|------|
| `ID` | `string` | `id` | 学科 OpenAlex ID |
| `DisplayName` | `string` | `display_name` | 学科名称 |
