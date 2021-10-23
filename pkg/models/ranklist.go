package models

type RankList struct {
	Model

	ProblemID uint64 `json:"-"`

	Name    string           `json:"name"`
	Title   string           `json:"title"`
	Models  []RankListModel  `json:"metrics"`
	Records []RankListRecord `json:"records"`
}

type RankListModel struct {
	Model
	RankListID uint64 `json:"-"`

	Key      string `json:"key"`
	Priority uint   `json:"priority"`
	Order    string `json:"order"`
}

type RankListRecord struct {
	Model
	RankListID uint64 `json:"-"`

	AccountID uint64  `json:"-"`
	Account   Account `json:"account"`

	Key   string  `json:"key"`
	Value float64 `json:"value"`
}
