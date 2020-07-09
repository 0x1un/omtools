package adtools

import "fmt"

type element string

type LinkNode struct {
	Data element
	Next *LinkNode
}

type LinkHead struct {
	Length uint
	Node   *LinkNode
}

type Linkeder interface {
	Push(data element)
	Insert(index uint, data element)
	Remove(index uint) error
	Len() uint
	Search(node element) int
	Get(index uint) *LinkNode
	Print()
}

func NewLinked() Linkeder {
	return &LinkHead{Length: 0, Node: &LinkNode{Data: "", Next: nil}}
}

func (l *LinkHead) Push(data element) {
	p := l.Node
	for p.Next != nil {
		p = p.Next
	}
	p.Next = &LinkNode{Data: data, Next: nil}
	l.Length++
}

func (l *LinkHead) Insert(index uint, data element) {
	if index >= l.Len() {
		return
	}
	var i uint
	p := l.Node
	for i = 0; i <= index; i++ {
		p = p.Next
	}
	p.Data = data
	l.Length++
}
func (l *LinkHead) Remove(index uint) error {
	if index >= l.Len() {
		return fmt.Errorf("index out of range: %d / %d", index, l.Len())
	}
	var i uint
	p := l.Node
	for i = 0; i < index; i++ {
		p = p.Next
	}
	p.Next = p.Next.Next
	l.Length--
	return nil
}
func (l *LinkHead) Len() uint {
	return l.Length
}
func (l *LinkHead) Search(node element) int {
	var i int
	p := l.Node
	for i = 0; p.Next != nil; i++ {
		p = p.Next
		if p.Data == node {
			return i
		}
	}
	return -1
}
func (l *LinkHead) Get(index uint) *LinkNode {
	if index >= l.Len() {
		return nil
	}
	var i uint
	p := l.Node
	for i = 0; i <= index; i++ {
		p = p.Next
	}
	return p
}
func (l *LinkHead) Print() {
	for p := l.Node; p.Next != nil; p = p.Next {
		println(p.Data)
	}
	println()
}
