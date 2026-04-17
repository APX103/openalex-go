# openalex-go 接口文档

面向调用方的 Go SDK 使用说明，基于 [OpenAlex](https://openalex.org/) 学术知识图谱 API。

## 安装

```bash
go get github.com/APX103/openalex-go
```

## 快速上手

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/APX103/openalex-go"
    "github.com/APX103/openalex-go/work"
    "github.com/APX103/openalex-go/util"
)

func main() {
    c := openalex.New(openalex.WithAPIKey("your-key"))

    resp, err := work.Search(context.Background(), c, work.SearchParams{
        Query:   "large language model",
        PerPage: 10,
        Sort:    &openalex.SortOption{Field: "cited_by_count"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("共 %d 条结果\n", resp.Meta.Count)
}
```

## 客户端配置

```go
c := openalex.New(
    openalex.WithAPIKey("key"),           // API Key，速率限制 10→100 次/秒
    openalex.WithMailto("you@example.com"), // 无 Key 时使用礼貌池
    openalex.WithTimeout(30 * time.Second), // 请求超时，默认 15s
    openalex.WithBaseURL("https://proxy"),   // 自定义代理地址
    openalex.WithHTTPClient(customClient),   // 完全替换 HTTP 客户端
)
```

无 API Key 时务必设置 `WithMailto`，否则共享池仅 100 credits/天。

## 通用类型

### PageParams — 分页参数

所有列表接口的分页控制：

```go
type PageParams struct {
    Page    int // 页码，默认 1
    PerPage int // 每页条数，默认 25，上限 200
}
```

零值 `openalex.PageParams{}` 等同于第 1 页、每页 25 条。

### SortOption — 排序

```go
type SortOption struct {
    Field string // 字段名
    Order string // "desc"（默认）或 "asc"
}
```

常用排序字段：`relevance_score`、`cited_by_count`、`publication_date`。

### ListResponse — 列表响应

所有列表接口返回同一泛型结构：

```go
type ListResponse[T any] struct {
    Meta    Meta      // Meta.Count 为结果总数
    Results []T       // 当前页数据
    GroupBy []GroupBy // 仅 group-by 查询时有值
}
```

### 错误处理

```go
type APIError struct {
    StatusCode int
    Message    string
    URL        string
}
```

HTTP 非 200 返回 `*APIError`；网络/解析错误返回 wrapped error。所有 API 函数均返回 `error`，调用方应始终检查：

```go
resp, err := work.Search(ctx, c, params)
if err != nil {
    // 判断是否为 API 错误
    var apiErr *openalex.APIError
    if errors.As(err, &apiErr) {
        log.Printf("API 错误: %d %s", apiErr.StatusCode, apiErr.Message)
    }
    return err
}
```

## work 包 — 论文

### Search — 搜索论文

```go
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Work], error)
```

```go
type SearchParams struct {
    Query   string                 // 搜索关键词
    Page    int                    // 页码
    PerPage int                    // 每页条数
    Sort    *openalex.SortOption   // 排序，nil 则按相关性
    Select  []string               // 返回字段白名单，减少传输量
    Filters map[string]string      // 过滤条件
}
```

示例：

```go
// 基本搜索
resp, _ := work.Search(ctx, c, work.SearchParams{
    Query: "transformer attention mechanism",
    PerPage: 20,
})

// 带过滤和排序
resp, _ := work.Search(ctx, c, work.SearchParams{
    Query:   "large language model",
    PerPage: 50,
    Sort:    &openalex.SortOption{Field: "cited_by_count"},
    Filters: map[string]string{
        "publication_year": "2024",
        "type":             "article",
        "is_oa":            "true",
    },
})

// 指定返回字段（生产环境推荐，减少响应体积）
resp, _ := work.Search(ctx, c, work.SearchParams{
    Query:  "deep learning",
    Select: []string{"id", "doi", "display_name", "publication_year", "cited_by_count"},
})
```

常用过滤字段：

| 过滤条件 | 示例值 | 说明 |
|---------|--------|------|
| `publication_year` | `"2024"` / `"2020-2024"` | 年份，支持范围 |
| `type` | `"article"` / `"preprint"` | 文献类型 |
| `primary_location.source.id` | `"S137773608"` | 期刊 ID |
| `author.id` | `"A5023898321"` | 作者 ID |
| `authorships.institutions.id` | `"I136199984"` | 机构 ID |
| `has_doi` | `"true"` | 是否有 DOI |
| `is_oa` | `"true"` | 是否开放获取 |
| `cited_by_count` | `">100"` / `"50-200"` | 被引次数，支持比较运算符和范围 |
| `concepts.id` | `"C154945302"` | 概念 ID |
| `open_access.oa_status` | `"gold"` / `"green"` | OA 类型 |
| `default.search_filter` | `"journal"` | 限定搜索范围 |

多条件用逗号分隔，SDK 会自动拼接：`Filters: map[string]string{"publication_year": "2024", "type": "article"}` → `filter=publication_year:2024,type:article`。

### Get — 获取单篇论文

```go
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Work, error)
```

```go
w, err := work.Get(ctx, c, "W2626778328", "id", "display_name", "abstract_inverted_index")
```

未找到时返回 `fmt.Errorf("work %s not found", id)`。

### GetByIDs — 批量获取论文

```go
func GetByIDs(ctx context.Context, c *openalex.Client, ids []string, selectFields ...string) ([]Work, error)
```

**单次上限 200 个 ID**，超出直接返回错误。

```go
works, err := work.GetByIDs(ctx, c, []string{"W1", "W2", "W3"}, "id", "display_name")
```

### GetCitedBy — 获取引用该论文的文献

```go
func GetCitedBy(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

默认按 `cited_by_count:desc` 排序。

```go
cited, _ := work.GetCitedBy(ctx, c, "W2626778328", openalex.PageParams{Page: 1, PerPage: 20})
```

### GetReferencedWorks — 获取该论文的参考文献

```go
func GetReferencedWorks(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

默认按 `cited_by_count:desc` 排序。

### GetRelated — 获取相关论文

```go
func GetRelated(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

基于 N-gram 相关性，默认按 `cited_by_count:desc` 排序。

### GetByAuthor — 获取作者的全部论文

```go
func GetByAuthor(ctx context.Context, c *openalex.Client, authorID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

默认按 `cited_by_count:desc` 排序。

### GetBySource — 获取期刊的全部论文

```go
func GetBySource(ctx context.Context, c *openalex.Client, sourceID string, page openalex.PageParams, sort *openalex.SortOption, selectFields ...string) (*openalex.ListResponse[Work], error)
```

默认按 `cited_by_count:desc` 排序，可通过 `sort` 参数覆盖（如按 `publication_date` 排序）。

### Work 结构体关键字段

```go
type Work struct {
    ID           string                // OpenAlex ID（完整 URL）
    Doi          string                // DOI（完整 URL，如 https://doi.org/10.APX103）
    DisplayName  string                // 论文标题
    PubYear      int                   // 发表年份
    PubDate      string                // 发表日期（YYYY-MM-DD）
    Type         string                // 类型：article, preprint, ...
    OpenAccess   OpenAccess            // 开放获取信息
    Authorships  []Authorship          // 作者列表及机构
    PrimaryLoc   *PrimaryLocation      // 主要发表位置（期刊等）
    BestOALoc    *PrimaryLocation      // 最佳 OA 位置（含 PDF 链接）
    Concepts     []Concept             // 关联概念
    Keywords     []Keyword             // 关键词
    CitedByCount int                   // 被引次数
    CountsByYear []CountByYear         // 逐年引用统计
    AbstractInv  map[string][]int      // 摘要倒排索引（需 RestoreAbstract 还原）
    Biblio       *Biblio               // 卷期页码信息
    IDs          WorkIDs               // 多平台 ID（arXiv, PMID, MAG 等）
}
```

注意：`Doi` 返回的是完整 URL（`https://doi.org/10.APX103`），如需纯 DOI 值需自行 `strings.TrimPrefix`。

### OpenAccess

```go
type OpenAccess struct {
    IsOA   bool    // 是否开放获取
    OaStatus string // OA 类型：gold, green, hybrid, bronze
    OaURL   *string // OA 链接
}
```

### Authorship

```go
type Authorship struct {
    AuthorPosition  string        // 作者序位：first, middle, last
    Author          AuthorRef     // 作者基本信息
    Institutions    []Institution // 作者机构
    IsCorresponding bool          // 是否通讯作者
}
```

## author 包 — 作者

### Search — 搜索作者

```go
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Author], error)
```

```go
type SearchParams struct {
    Query   string   // 搜索关键词
    Page    int      // 页码
    PerPage int      // 每页条数
    Select  []string // 返回字段
}
```

默认按 `relevance_score:desc` 排序。

```go
resp, _ := author.Search(ctx, c, author.SearchParams{
    Query:   "Andrew Ng",
    PerPage: 10,
})
```

### Get — 获取单个作者

```go
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Author, error)
```

```go
a, _ := author.Get(ctx, c, "A5023898321", "id", "display_name", "works_count", "cited_by_count")
```

### Author 结构体关键字段

```go
type Author struct {
    ID             string         // OpenAlex ID
    DisplayName    string         // 作者姓名
    Orcid          string         // ORCID（完整 URL）
    WorksCount     int            // 论文总数
    CitedByCount   int            // 总被引次数
    SummaryStats   SummaryStats   // h_index, i10_index 等
    LastKnownInsts []Institution  // 最近已知机构
    Topics         []AuthorTopic  // 研究主题
    XConcepts      []Concept      // 关联概念
    CountsByYear   []CountByYear  // 逐年引用统计
}
```

### SummaryStats

```go
type SummaryStats struct {
    HIndex             int     // h 指数
    I10Index           int     // i10 指数
    TwoYrMeanCitedness float64 // 近两年平均被引次数
}
```

## source 包 — 期刊

### Search — 搜索期刊

```go
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Source], error)
```

参数与 author.Search 相同，默认按 `relevance_score:desc` 排序。

```go
resp, _ := source.Search(ctx, c, source.SearchParams{
    Query:   "Nature",
    PerPage: 10,
})
```

### Get — 获取单个期刊

```go
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Source, error)
```

```go
s, _ := source.Get(ctx, c, "S137773608")
```

### Source 结构体关键字段

```go
type Source struct {
    ID           string         // OpenAlex ID
    DisplayName  string         // 期刊名称
    ISSN         []string       // ISSN 列表
    ISSNL        string         // ISSN-L（链接 ISSN）
    IsOA         bool           // 是否 OA 期刊
    Type         string         // 类型：journal, repository, ...
    WorksCount   int            // 收录论文总数
    CitedByCount int            // 总被引次数
    SummaryStats SummaryStats   // h_index 等
    HomepageURL  *string        // 期刊主页（可为 nil）
    HostOrgName  *string        // 出版机构（可为 nil）
    APCUSD       *float64       // APC 费用（美元，可为 nil）
    CountryCode  string         // 国家代码
    TopicShare   []AuthorTopic  // 主题占比（含 value 分数）
    CountsByYear []CountByYear  // 逐年统计
}
```

注意：`HomepageURL`、`HostOrgName`、`APCUSD` 为指针类型，使用前需判空。

## util 包 — 工具函数

### ShortID — 提取短 ID

```go
func ShortID(openalexURL string) string
```

```go
util.ShortID("https://openalex.org/W2626778328") // → "W2626778328"
```

所有 API 返回的 ID 均为完整 URL，需提取短 ID 用于过滤和展示。

### JoinPipe — 拼接批量查询 ID

```go
func JoinPipe(ids []string) string
```

```go
util.JoinPipe([]string{"W1", "W2", "W3"}) // → "W1|W2|W3"
```

### RestoreAbstract — 还原摘要文本

```go
func RestoreAbstract(idx map[string][]int) string
```

OpenAlex 摘要以倒排索引格式存储（`{word: [position1, position2, ...]}`），此函数还原为纯文本：

```go
text := util.RestoreAbstract(w.AbstractInv)
```

### ResolvePDF — 解析 PDF 链接

```go
func ResolvePDF(w PDFWork) PDFResult

type PDFResult struct {
    URL    string     // PDF URL（可能为空）
    Source PDFSource  // 来源类型
}
```

优先级链：arXiv → OpenAlex best_oa_location → open_access.oa_url → DOI 跳转 → 无可用 PDF。

```go
pdf := util.ResolvePDF(w)
if pdf.URL != "" {
    fmt.Printf("PDF: %s (来源: %s)\n", pdf.URL, util.PDFSourceName(pdf.Source))
}
```

来源枚举：

| 常量 | 值 | 说明 |
|-----|---|------|
| `util.PDFSourceArXiv` | 1 | arXiv 直链 |
| `util.PDFSourceOpenAlex` | 2 | OpenAlex OA 记录 |
| `util.PDFSourceUnpaywall` | 3 | 预留（未实现） |
| `util.PDFSourceDOI` | 4 | DOI 跳转 |
| `util.PDFSourceNone` | 5 | 无可用 PDF |

```go
util.PDFSourceName(util.PDFSourceArXiv) // → "arxiv"
```

## 生产环境注意事项

### Select 字段白名单

生产环境所有请求都应指定 `select`，只拉取必要字段。OpenAlex 完整 Work 记录可达数十 KB，而列表页通常只需要少量字段：

```go
// 列表页：不需要摘要和逐年统计
Select: []string{"id", "doi", "display_name", "publication_year", "publication_date",
    "type", "open_access", "authorships", "primary_location", "concepts", "cited_by_count", "ids"}

// 详情页：需要完整信息
Select: []string{"id", "doi", "display_name", "publication_year", "publication_date",
    "type", "open_access", "authorships", "primary_location", "best_oa_location",
    "concepts", "cited_by_count", "counts_by_year", "abstract_inverted_index", "ids", "biblio"}

// 批量查询：最小字段集
Select: []string{"id", "doi", "display_name", "publication_year", "publication_date",
    "type", "open_access", "authorships", "primary_location", "best_oa_location",
    "concepts", "cited_by_count", "ids"}
```

### 速率限制

| 场景 | 速率 |
|------|------|
| 无认证 | 10 次/秒，100 credits/天 |
| 仅 Mailto | 礼貌池，约 10 次/秒 |
| API Key | 100 次/秒 |

API Key 在 [openalex.org](https://openalex.org/) 免费申请。

### 分页限制

- `per_page` 上限 **200**（API 硬限制）
- `page` 分页在约 10,000 页后性能下降，大数据量场景建议用 filter 缩小范围
- `GetByIDs` 单次最多 **200** 个 ID，超出需分批

### OpenAlex ID 格式

API 返回的 ID 是完整 URL（`https://openalex.org/W123`），但 SDK 所有接口**接受短 ID**：

```go
work.Get(ctx, c, "W2626778328")          // 短 ID，可用
work.Get(ctx, c, "https://openalex.org/W2626778328") // 完整 URL，也可用
```

### counts_by_year 无序警告

`CountsByYear` 字段在 API 响应中**不保证按年份排序**。前端使用前需按 `Year` 排序：

```go
slices.SortFunc(w.CountsByYear, func(a, b CountByYear) int {
    return a.Year - b.Year
})
```

### 摘要还原

`AbstractInv` 是 `map[string][]int`，不能直接使用。必须调用 `util.RestoreAbstract()` 还原为可读文本。还原结果为英文纯文本。

### DOI 格式

`Work.Doi` 返回完整 URL（`https://doi.org/10.APX103`）。如需纯 DOI 值：

```go
doi := strings.TrimPrefix(w.Doi, "https://doi.org/")
```

### 指针类型字段

以下字段在 API 响应中可能为 null，SDK 用指针类型表示，使用前**必须判空**：

| 字段 | 类型 | 所在结构体 |
|-----|------|----------|
| `BestOALoc` | `*PrimaryLocation` | Work |
| `PrimaryLoc` | `*PrimaryLocation` | Work |
| `Biblio` | `*Biblio` | Work |
| `OaURL` | `*string` | OpenAccess |
| `HomepageURL` | `*string` | Source |
| `HostOrgName` | `*string` | Source |
| `APCUSD` | `*float64` | Source |
| `PdfURL` | `*string`（在 `PrimaryLocation` 内） | PrimaryLocation |

### 错误处理最佳实践

```go
resp, err := work.Search(ctx, c, params)
if err != nil {
    var apiErr *openalex.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case 404:
            // 资源不存在
        case 429:
            // 速率限制，应加入退避重试
        case 403:
            // 认证失败
        default:
            // 其他服务端错误
        }
    }
    return err
}
```

SDK 本身**不内置重试逻辑**，速率限制处理需调用方自行实现。
