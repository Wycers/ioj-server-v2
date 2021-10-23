package models

type Problem struct {
	Model

	Name      string `json:"name" gorm:"unique_index:idx2"`
	Title     string `json:"title"`
	ProcessId uint64

	PublicVolume  string `json:"public_volume"`
	PrivateVolume string `json:"-"`

	RankLists []RankList `json:"rank_lists"`
}
