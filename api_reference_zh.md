# openalex-go API 参考

`go get github.com/APX103/openalex-go`

## 初始化

```go
c := openalex.New(opts...)
```

| Option | 签名 | 说明 |
|--------|------|------|
| `WithAPIKey` | `string` | 设置 API Key（优先级高于 Mailto） |
| `WithMailto` | `string` | 设置邮箱，无 Key 时进入礼貌池 |
| `WithTimeout` | `time.Duration` | 请求超时，默认 15s |
| `WithBaseURL` | `string` | 自定义 API 地址，默认 `https://api.openalex.org` |
| `WithHTTPClient` | `*http.Client` | 替换默认 HTTP 客户端 |

认证逻辑：`apiKey` 非空 → `?api_key=xxx`；否则 `mailto` 非空 → `?mailto=xxx`；否则无认证。

---

## 通用类型

### PageParams

```go
type PageParams struct {
    Page    int // 页码，默认 1
    PerPage int // 每页条数，默认 25，上限 200
}
```

### SortOption

```go
type SortOption struct {
    Field string // 字段名
    Order string // "desc"（默认）或 "asc"
}
```

### ListResponse[T]

```go
type ListResponse[T any] struct {
    Meta    Meta      // Count: 结果总数
    Results []T       // 当前页数据
    GroupBy []GroupBy // group-by 查询时有值
}
```

### APIError

```go
type APIError struct {
    StatusCode int
    Message    string
    URL        string
}
```

---

## work 包 — 论文

### work.Search

```go
func work.Search(ctx context.Context, c *openalex.Client, params work.SearchParams) (*openalex.ListResponse[work.Work], error)
```

**SearchParams**

| 字段 | 类型 | 说明 |
|------|------|------|
| Query | `string` | 搜索关键词 |
| Page | `int` | 页码 |
| PerPage | `int` | 每页条数 |
| Sort | `*openalex.SortOption` | 排序，nil 按相关性 |
| Select | `[]string` | 返回字段白名单 |
| Filters | `map[string]string` | 过滤条件 |

```go
resp, err := work.Search(ctx, c, work.SearchParams{
    Query:   "large language model",
    PerPage: 20,
    Sort:    &openalex.SortOption{Field: "cited_by_count"},
    Select:  []string{"id", "display_name", "cited_by_count"},
    Filters: map[string]string{"publication_year": "2024", "type": "article"},
})
// → resp.Results []work.Work, resp.Meta.Count 总数
```

**Filters 常用键**

| 键 | 值示例 | 说明 |
|----|--------|------|
| `publication_year` | `"2024"` / `"2020-2024"` | 年份 |
| `type` | `"article"` / `"preprint"` | 文献类型 |
| `primary_location.source.id` | `"S137773608"` | 期刊 |
| `author.id` | `"A5023898321"` | 作者 |
| `authorships.institutions.id` | `"I136199984"` | 机构 |
| `has_doi` | `"true"` | 有 DOI |
| `is_oa` | `"true"` | 开放获取 |
| `cited_by_count` | `">100"` / `"50-200"` | 被引次数 |
| `concepts.id` | `"C154945302"` | 概念 |
| `open_access.oa_status` | `"gold"` / `"green"` | OA 类型 |
| `default.search_filter` | `"journal"` | 搜索范围 |

---

### work.Get

```go
func work.Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*work.Work, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| id | `string` | 论文 ID（短 ID 或完整 URL） |
| selectFields | `...string` | 可选，返回字段白名单 |

```go
w, err := work.Get(ctx, c, "W2626778328", "id", "display_name", "abstract_inverted_index")
// → *work.Work
// 未找到返回 fmt.Errorf("work %s not found", id)
```

---

### work.GetByIDs

```go
func work.GetByIDs(ctx context.Context, c *openalex.Client, ids []string, selectFields ...string) ([]work.Work, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| ids | `[]string` | 论文 ID 列表，上限 **200** |
| selectFields | `...string` | 可选，返回字段白名单 |

```go
works, err := work.GetByIDs(ctx, c, []string{"W1", "W2", "W3"}, "id", "display_name")
// → []work.Work
// 超过 200 返回 fmt.Errorf("GetByIDs: max 200 IDs per request, got %d", n)
```

---

### work.GetCitedBy

```go
func work.GetCitedBy(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回引用 `workID` 的论文。固定排序 `cited_by_count:desc`。

```go
cited, err := work.GetCitedBy(ctx, c, "W2626778328", openalex.PageParams{Page: 1, PerPage: 20})
```

---

### work.GetReferencedWorks

```go
func work.GetReferencedWorks(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回 `workID` 的参考文献。固定排序 `cited_by_count:desc`。

---

### work.GetRelated

```go
func work.GetRelated(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回与 `workID` 相关的论文（N-gram）。固定排序 `cited_by_count:desc`。

---

### work.GetByAuthor

```go
func work.GetByAuthor(ctx context.Context, c *openalex.Client, authorID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回指定作者的论文。固定排序 `cited_by_count:desc`。

---

### work.GetBySource

```go
func work.GetBySource(ctx context.Context, c *openalex.Client, sourceID string, page openalex.PageParams, sort *openalex.SortOption, selectFields ...string) (*openalex.ListResponse[work.Work], error)
```

返回指定期刊的论文。`sort` 为 nil 时默认 `cited_by_count:desc`，可覆盖。

```go
works, err := work.GetBySource(ctx, c, "S137773608",
    openalex.PageParams{Page: 1, PerPage: 20},
    &openalex.SortOption{Field: "publication_date"})
```

---

### work.Work 结构体

```go
type Work struct {
    ID           string              `json:"id"`                        // 完整 URL
    Doi          string              `json:"doi"`                       // 完整 URL
    DisplayName  string              `json:"display_name"`              // 标题
    PubYear      int                 `json:"publication_year"`          // 年份
    PubDate      string              `json:"publication_date"`          // YYYY-MM-DD
    Type         string              `json:"type"`                      // article/preprint/...
    Language     string              `json:"language,omitempty"`
    IndexedIn    []string            `json:"indexed_in,omitempty"`
    OpenAccess   OpenAccess          `json:"open_access"`
    Authorships  []Authorship        `json:"authorships"`
    PrimaryLoc   *PrimaryLocation    `json:"primary_location"`          // 可为 nil
    BestOALoc    *PrimaryLocation    `json:"best_oa_location,omitempty"`// 可为 nil
    Topics       []WorkTopic         `json:"topics,omitempty"`
    Concepts     []Concept           `json:"concepts,omitempty"`
    Keywords     []Keyword           `json:"keywords,omitempty"`
    Refs         []string            `json:"referenced_works,omitempty"`
    Related      []string            `json:"related_works,omitempty"`
    CountsByYear []CountByYear       `json:"counts_by_year,omitempty"`  // 不保证有序
    CitedByCount int                 `json:"cited_by_count"`
    AbstractInv  map[string][]int    `json:"abstract_inverted_index,omitempty"` // 需 util.RestoreAbstract
    Biblio       *Biblio             `json:"biblio,omitempty"`          // 可为 nil
    IDs          WorkIDs             `json:"ids,omitempty"`
}
```

### work.OpenAccess

```go
type OpenAccess struct {
    IsOA               bool    `json:"is_oa"`
    OaStatus           string  `json:"oa_status"`           // gold/green/hybrid/bronze
    OaURL              *string `json:"oa_url"`              // 可为 nil
    AnyRepoHasFulltext bool    `json:"any_repository_has_fulltext"`
}
```

### work.Authorship

```go
type Authorship struct {
    AuthorPosition  string        `json:"author_position"`  // first/middle/last
    Author          AuthorRef     `json:"author"`
    Institutions    []Institution `json:"institutions"`
    Countries       []string      `json:"countries,omitempty"`
    IsCorresponding bool          `json:"is_corresponding"`
}
```

### work.AuthorRef

```go
type AuthorRef struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
    Orcid       string `json:"orcid,omitempty"`
}
```

### work.Institution

```go
type Institution struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
    Ror         string `json:"ror,omitempty"`
    CountryCode string `json:"country_code,omitempty"`
    Type        string `json:"type,omitempty"`
}
```

### work.PrimaryLocation

```go
type PrimaryLocation struct {
    IsOA           bool            `json:"is_oa"`
    LandingPageURL string          `json:"landing_page_url"`
    PdfURL         *string         `json:"pdf_url,omitempty"`       // 可为 nil
    Source         *LocationSource `json:"source,omitempty"`       // 可为 nil
    License        string          `json:"license,omitempty"`
    Version        string          `json:"version,omitempty"`
}
```

### work.LocationSource

```go
type LocationSource struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
    ISSN        any    `json:"issn,omitempty"`
    Type        string `json:"type"`
    IsOA        bool   `json:"is_oa"`
}
```

### work.WorkTopic

```go
type WorkTopic struct {
    ID          string   `json:"id"`
    DisplayName string   `json:"display_name"`
    Count       int      `json:"count"`
    Subfield    TopicRef `json:"subfield"`
    Field       TopicRef `json:"field"`
    Domain      TopicRef `json:"domain"`
}
```

### work.Keyword

```go
type Keyword struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
}
```

### work.Biblio

```go
type Biblio struct {
    Volume    string `json:"volume,omitempty"`
    Issue     string `json:"issue,omitempty"`
    FirstPage string `json:"first_page,omitempty"`
    LastPage  string `json:"last_page,omitempty"`
}
```

### work.WorkIDs

```go
type WorkIDs struct {
    OpenAlex string `json:"openalex"`
    Doi      string `json:"doi"`
    Mag      string `json:"mag"`
    PMID     string `json:"pmid,omitempty"`
    ArXiv    string `json:"arxiv,omitempty"`
}
```

---

## author 包 — 作者

### author.Search

```go
func author.Search(ctx context.Context, c *openalex.Client, params author.SearchParams) (*openalex.ListResponse[author.Author], error)
```

固定排序 `relevance_score:desc`。

| 参数 | 类型 | 说明 |
|------|------|------|
| Query | `string` | 搜索关键词 |
| Page | `int` | 页码 |
| PerPage | `int` | 每页条数 |
| Select | `[]string` | 返回字段白名单 |

```go
resp, err := author.Search(ctx, c, author.SearchParams{Query: "Andrew Ng", PerPage: 10})
```

### author.Get

```go
func author.Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*author.Author, error)
```

```go
a, err := author.Get(ctx, c, "A5023898321", "id", "display_name", "works_count")
```

### author.Author 结构体

```go
type Author struct {
    ID             string            `json:"id"`
    DisplayName    string            `json:"display_name"`
    Orcid          string            `json:"orcid,omitempty"`        // 完整 URL
    WorksCount     int               `json:"works_count"`
    CitedByCount   int               `json:"cited_by_count"`
    SummaryStats   SummaryStats      `json:"summary_stats"`
    LastKnownInsts []Institution     `json:"last_known_institutions"`
    Topics         []AuthorTopic     `json:"topics,omitempty"`
    XConcepts      []Concept         `json:"x_concepts,omitempty"`
    CountsByYear   []CountByYear     `json:"counts_by_year,omitempty"` // 不保证有序
    WorksAPIURL    string            `json:"works_api_url,omitempty"`
}
```

---

## source 包 — 期刊

### source.Search

```go
func source.Search(ctx context.Context, c *openalex.Client, params source.SearchParams) (*openalex.ListResponse[source.Source], error)
```

固定排序 `relevance_score:desc`。参数同 author.Search。

```go
resp, err := source.Search(ctx, c, source.SearchParams{Query: "Nature", PerPage: 10})
```

### source.Get

```go
func source.Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*source.Source, error)
```

```go
s, err := source.Get(ctx, c, "S137773608")
```

### source.Source 结构体

```go
type Source struct {
    ID           string         `json:"id"`
    DisplayName  string         `json:"display_name"`
    ISSN         []string       `json:"issn,omitempty"`
    ISSNL        string         `json:"issn_l,omitempty"`
    IsOA         bool           `json:"is_oa"`
    Type         string         `json:"type"`
    WorksCount   int            `json:"works_count"`
    CitedByCount int            `json:"cited_by_count"`
    SummaryStats SummaryStats   `json:"summary_stats"`
    HomepageURL  *string        `json:"homepage_url,omitempty"`           // 可为 nil
    HostOrgName  *string        `json:"host_organization_name,omitempty"`// 可为 nil
    APCUSD       *float64       `json:"apc_usd,omitempty"`               // 可为 nil
    CountryCode  string         `json:"country_code,omitempty"`
    Topics       []AuthorTopic  `json:"topics,omitempty"`
    TopicShare   []AuthorTopic  `json:"topic_share,omitempty"`
    CountsByYear []CountByYear  `json:"counts_by_year,omitempty"`        // 不保证有序
    WorksAPIURL  string         `json:"works_api_url,omitempty"`
}
```

---

## util 包 — 工具函数

### util.ShortID

```go
func util.ShortID(openalexURL string) string
```

```
util.ShortID("https://openalex.org/W2626778328") → "W2626778328"
```

### util.JoinPipe

```go
func util.JoinPipe(ids []string) string
```

```
util.JoinPipe([]string{"W1", "W2"}) → "W1|W2"
```

### util.RestoreAbstract

```go
func util.RestoreAbstract(idx map[string][]int) string
```

将 OpenAlex 倒排索引还原为纯文本。`nil` 返回空字符串。

```go
text := util.RestoreAbstract(w.AbstractInv)
```

### util.ResolvePDF

```go
func util.ResolvePDF(w util.PDFWork) util.PDFResult
```

解析 PDF 链接，优先级：arXiv → best_oa_location.pdf_url → open_access.oa_url → DOI → 无。

```go
type PDFResult struct {
    URL    string     // PDF 地址，可能为空
    Source PDFSource  // 来源
}

pdf := util.ResolvePDF(&w)
fmt.Println(pdf.URL, pdf.Source)
```

### util.PDFSource 常量

| 常量 | 值 | PDFSourceName |
|------|----|---------------|
| `PDFSourceArXiv` | 1 | `"arxiv"` |
| `PDFSourceOpenAlex` | 2 | `"openalex"` |
| `PDFSourceUnpaywall` | 3 | `"unpaywall"` (预留) |
| `PDFSourceDOI` | 4 | `"doi"` |
| `PDFSourceNone` | 5 | `""` |

```go
util.PDFSourceName(util.PDFSourceArXiv) → "arxiv"
```

---

## 共享类型

### SummaryStats（Author / Source 共用）

```go
type SummaryStats struct {
    HIndex             int     `json:"h_index"`
    I10Index           int     `json:"i10_index"`
    TwoYrMeanCitedness float64 `json:"2yr_mean_citedness"`
}
```

### AuthorTopic（Author.Topics / Source.Topics / Source.TopicShare 共用）

```go
type AuthorTopic struct {
    ID          string   `json:"id"`
    DisplayName string   `json:"display_name"`
    Count       int      `json:"count"`
    Value       float64  `json:"value,omitempty"`
    Subfield    TopicRef `json:"subfield"`
    Field       TopicRef `json:"field"`
    Domain      TopicRef `json:"domain"`
}
```

### CountByYear（Work / Author / Source 共用）

```go
type CountByYear struct {
    Year         int `json:"year"`
    WorksCount   int `json:"works_count"`
    CitedByCount int `json:"cited_by_count"`
}
```

### Concept（Work.Concepts / Author.XConcepts 共用）

```go
type Concept struct {
    ID          string  `json:"id"`
    DisplayName string  `json:"display_name"`
    Score       float64 `json:"score"`
}
```

### TopicRef（WorkTopic / AuthorTopic 共用）

```go
type TopicRef struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
}
```
