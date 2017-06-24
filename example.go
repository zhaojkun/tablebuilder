package main

type Page struct {
	ID   int    `json:"id"`   // 主键
	URL  string `json:"url"`  // 名称
	Name string `json:"name"` // 别名
	Hash string `json:"hash"` // 哈希
	Body string `json:"body"` // 内容
}
