package main

type Meta struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Count   int    `json:"count,omitempty"`
}

type DataItem struct {
	Kind string       `json:"kind"`
	Data DataItemData `json:"data"`
}

type DataItemData struct {
	Members *[]DataItemMember `json:"members"`
	GroupId string            `json:"groupId"`
}

type DataItemMember struct {
	GroupId string               `json:"groupId"`
	Member  DataItemMemberMember `json:"member"`
}

type DataItemMemberMember struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}
