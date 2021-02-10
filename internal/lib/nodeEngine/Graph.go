package nodeEngine

import (
	"fmt"
)

type Graph struct {
	Blocks map[int]*Block
	Links  map[int]*Link
}

type Block struct {
	Id         int
	Type       string
	Properties map[string]interface{}

	Inputs []int
	Output [][]int

	Status string
}

type Port struct {
	Id   int
	Slot int
}

type Link struct {
	Id      int
	Source  Port
	Target  Port
	IsReady bool
}

func (b *Block) setProperty(key string, value string) {
	if b.Properties == nil {
		b.Properties = make(map[string]interface{})
	}
	b.Properties[key] = value
}

func (g *Graph) AddBlock(id int, tp string, properties map[string]interface{}, inputs []int, outputs [][]int) *Block {
	g.Blocks[id] = &Block{
		Id:         id,
		Type:       tp,
		Properties: properties,
		Inputs:     inputs,
		Output:     outputs,

		Status: "pending",
	}
	return g.Blocks[id]
}

func (g *Graph) AddLink(id, sourceId, sourceSlot, targetId, targetSlot int) *Link {
	g.Links[id] = &Link{
		id,
		Port{
			Id:   sourceId,
			Slot: sourceSlot,
		},
		Port{
			Id:   targetId,
			Slot: targetSlot,
		},
		false,
	}
	//fmt.Println("==>", g.Links[id])
	return g.Links[id]
}

func (g *Graph) FindBlockById(id int) *Block {
	if block, ok := g.Blocks[id]; ok {
		return block
	} else {
		return nil
	}
}

func (g *Graph) findLinkById(id int) *Link {
	if link, ok := g.Links[id]; ok {
		return link
	} else {
		return nil
	}
}

func (g *Graph) FindLinkBySourcePort(id, slot int) []*Link {
	var res []*Link
	for _, v := range g.Links {
		if v.Source.Id == id && v.Source.Slot == slot {
			res = append(res, v)
		}
	}
	return res
}
func (g *Graph) findLinkByTargetPort(id, slot int) *Link {
	for _, v := range g.Links {
		if v.Target.Id == id && v.Target.Slot == slot {
			fmt.Println(v)
			return v
		}
	}
	return nil
}

func (b *Block) Done() {
	b.Status = "done"
}
func (b *Block) ReSet() {
	b.Status = "pending"
}

func (g *Graph) Run() []*Block {

	var res []*Block

	for _, block := range g.Blocks {
		flag := true

		if block.Status != "pending" {
			continue
		}

		for _, inputLink := range block.Inputs {
			link := g.findLinkById(inputLink)
			if link == nil {
				panic("Wrong file")
			}

			source := g.FindBlockById(link.Source.Id)
			if source == nil {
				panic("Wrong file")
			}

			if source.Status != "done" {
				flag = false
				break
			}
		}

		if flag {
			block.Status = "in queue"
			res = append(res, block)
		}
	}

	return res
}

func New() *Graph {
	return &Graph{
		Blocks: map[int]*Block{},
		Links:  map[int]*Link{},
	}
}
