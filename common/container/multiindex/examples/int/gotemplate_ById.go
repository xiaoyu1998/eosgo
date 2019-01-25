// Code generated by gotemplate. DO NOT EDIT.

package int

import (
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/multiindex"
)

// template type HashedUniqueIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Hash,KeyFunc)

type ById struct {
	super *ByNum            // index on the HashedUniqueIndex, IndexBase is the last super index
	final *TestIndex        // index under the HashedUniqueIndex, MultiIndex is the final index
	inner map[int]*ByIdNode // use hashmap to safe HashedUniqueIndex's k/v(HashedUniqueIndexNode)
}

func (i *ById) init(final *TestIndex) {
	i.final = final
	i.inner = map[int]*ByIdNode{}
	i.super = &ByNum{}
	i.super.init(final)
}

/*generic class*/

/*generic class*/

type ByIdNode struct {
	super *ByNumNode     // index-node on the HashedUniqueIndexNode, IndexBaseNode is the last super node
	final *TestIndexNode // index-node under the HashedUniqueIndexNode, MultiIndexNode is the final index
	hash  int            // k of hashmap
}

/*generic class*/

/*generic class*/

func (i *ById) GetSuperIndex() interface{} { return i.super }
func (i *ById) GetFinalIndex() interface{} { return i.final }

func (n *ByIdNode) GetSuperNode() interface{} { return n.super }
func (n *ByIdNode) GetFinalNode() interface{} { return n.final }

func (n *ByIdNode) value() *int {
	return n.super.value()
}

func (i *ById) Size() int {
	return len(i.inner)
}

func (i *ById) Empty() bool {
	return len(i.inner) == 0
}

func (i *ById) clear() {
	i.inner = map[int]*ByIdNode{}
	i.super.clear()
}

func (i *ById) Insert(v int) (IteratorById, bool) {
	fn, res := i.final.insert(v)
	if res {
		return i.makeIterator(fn), true
	}
	return i.End(), false
}

func (i *ById) insert(v int, fn *TestIndexNode) (*ByIdNode, bool) {
	hash := ByIdHashFunc(v)
	node := ByIdNode{hash: hash}
	if _, ok := i.inner[hash]; ok {
		container.Logger.Warn("#hash index insert failed")
		return nil, false
	}
	i.inner[hash] = &node
	sn, res := i.super.insert(v, fn)
	if res {
		node.final = fn
		node.super = sn
		return &node, true
	}
	delete(i.inner, hash)
	return nil, false
}

func (i *ById) Find(k int) (IteratorById, bool) {
	node, res := i.inner[k]
	if res {
		return IteratorById{i, node, betweenById}, true
	}
	return i.End(), false
}

func (i *ById) Each(f func(key int, obj int)) {
	for k, v := range i.inner {
		f(k, *v.value())
	}
}

func (i *ById) Erase(iter IteratorById) {
	i.final.erase(iter.node.final)
}

func (i *ById) erase(n *ByIdNode) {
	delete(i.inner, n.hash)
	i.super.erase(n.super)
}

func (i *ById) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(IteratorById); ok {
		i.Erase(itr)
	} else {
		i.super.erase_(iter)
	}
}

func (i *ById) Modify(iter IteratorById, mod func(*int)) bool {
	if _, b := i.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (i *ById) modify(n *ByIdNode) (*ByIdNode, bool) {
	delete(i.inner, n.hash)

	hash := ByIdHashFunc(*n.value())
	if _, exist := i.inner[hash]; exist {
		container.Logger.Warn("#hash index modify failed")
		i.super.erase(n.super)
		return nil, false
	}

	i.inner[hash] = n

	if sn, res := i.super.modify(n.super); !res {
		delete(i.inner, hash)
		return nil, false
	} else {
		n.super = sn
	}

	return n, true
}

func (i *ById) modify_(iter multiindex.IteratorType, mod func(*int)) bool {
	if itr, ok := iter.(IteratorById); ok {
		return i.Modify(itr, mod)
	} else {
		return i.super.modify_(iter, mod)
	}
}

func (i *ById) Values() []int {
	vs := make([]int, 0, i.Size())
	i.Each(func(key int, obj int) {
		vs = append(vs, obj)
	})
	return vs
}

type IteratorById struct {
	index    *ById
	node     *ByIdNode
	position posById
}

type posById byte

const (
	//begin   = 0
	betweenById = 1
	endById     = 2
)

func (i *ById) makeIterator(fn *TestIndexNode) IteratorById {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*ByIdNode); ok {
			return IteratorById{i, n, betweenById}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

func (i *ById) End() IteratorById {
	return IteratorById{i, nil, endById}
}

func (iter IteratorById) Value() (v int) {
	if iter.position == betweenById {
		return *iter.node.value()
	}
	return
}

func (iter IteratorById) HasNext() bool {
	container.Logger.Warn("hashed index iterator is unmoveable")
	return false
}

func (iter IteratorById) IsEnd() bool {
	return iter.position == endById
}