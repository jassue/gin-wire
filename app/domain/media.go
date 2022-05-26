package domain

type Media struct {
    ID uint64 `json:"id"`
    DiskType string `json:"-"`
    SrcType int8 `json:"-"`
    Src string `json:"src"`
    Url string `json:"url"`
}
