// Code generated by gotemplate. DO NOT EDIT.

package transaction_state

import (
	"fmt"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/container/multiindex"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply)

// OrderedIndex holds elements of the red-black tree
type ByExpiry struct {
	super *ByBlockNum            // index on the OrderedIndex, IndexBase is the last super index
	final *TransactionStateIndex // index under the OrderedIndex, MultiIndex is the final index

	Root *ByExpiryNode
	size int
}

func (tree *ByExpiry) init(final *TransactionStateIndex) {
	tree.final = final
	tree.super = &ByBlockNum{}
	tree.super.init(final)
}

func (tree *ByExpiry) clear() {
	tree.Clear()
	tree.super.clear()
}

/*generic class*/

/*generic class*/

// OrderedIndexNode is a single element within the tree
type ByExpiryNode struct {
	Key    common.TimePointSec
	super  *ByBlockNumNode
	final  *TransactionStateIndexNode
	color  colorByExpiry
	Left   *ByExpiryNode
	Right  *ByExpiryNode
	Parent *ByExpiryNode
}

/*generic class*/

/*generic class*/

func (node *ByExpiryNode) value() *multi_index.TransactionState {
	return node.super.value()
}

type colorByExpiry bool

const (
	blackByExpiry, redByExpiry colorByExpiry = true, false
)

func (tree *ByExpiry) Insert(v multi_index.TransactionState) (IteratorByExpiry, bool) {
	fn, res := tree.final.insert(v)
	if res {
		return tree.makeIterator(fn), true
	}
	return tree.End(), false
}

func (tree *ByExpiry) insert(v multi_index.TransactionState, fn *TransactionStateIndexNode) (*ByExpiryNode, bool) {
	key := ByExpiryFunc(v)

	node, res := tree.put(key)
	if !res {
		log.Warn("#ordered index insert failed")
		return nil, false
	}
	sn, res := tree.super.insert(v, fn)
	if res {
		node.super = sn
		node.final = fn
		return node, true
	}
	tree.remove(node)
	return nil, false
}

func (tree *ByExpiry) Erase(iter IteratorByExpiry) (itr IteratorByExpiry) {
	itr = iter
	itr.Next()
	tree.final.erase(iter.node.final)
	return
}

func (tree *ByExpiry) Erases(first, last IteratorByExpiry) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *ByExpiry) erase(n *ByExpiryNode) {
	tree.remove(n)
	tree.super.erase(n.super)
	n.super = nil
	n.final = nil
}

func (tree *ByExpiry) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(IteratorByExpiry); ok {
		tree.Erase(itr)
	} else {
		tree.super.erase_(iter)
	}
}

func (tree *ByExpiry) Modify(iter IteratorByExpiry, mod func(*multi_index.TransactionState)) bool {
	if _, b := tree.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (tree *ByExpiry) modify(n *ByExpiryNode) (*ByExpiryNode, bool) {
	n.Key = ByExpiryFunc(*n.value())

	if !tree.inPlace(n) {
		tree.remove(n)
		node, res := tree.put(n.Key)
		if !res {
			log.Warn("#ordered index modify failed")
			tree.super.erase(n.super)
			return nil, false
		}

		//n.Left = node.Left
		//if n.Left != nil {
		//	n.Left.Parent = n
		//}
		//n.Right = node.Right
		//if n.Right != nil {
		//	n.Right.Parent = n
		//}
		//n.Parent = node.Parent
		//if n.Parent != nil {
		//	if n.Parent.Left == node {
		//		n.Parent.Left = n
		//	} else {
		//		n.Parent.Right = n
		//	}
		//} else {
		//	tree.Root = n
		//}
		node.super = n.super
		node.final = n.final
		n = node
	}

	if sn, res := tree.super.modify(n.super); !res {
		tree.remove(n)
		return nil, false
	} else {
		n.super = sn
	}

	return n, true
}

func (tree *ByExpiry) modify_(iter multiindex.IteratorType, mod func(*multi_index.TransactionState)) bool {
	if itr, ok := iter.(IteratorByExpiry); ok {
		return tree.Modify(itr, mod)
	} else {
		return tree.super.modify_(iter, mod)
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByExpiry) Find(key common.TimePointSec) IteratorByExpiry {
	if true {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && ByExpiryCompare(key, lower.Key()) == 0 {
			return lower
		}
		return tree.End()
	} else {
		if node := tree.lookup(key); node != nil {
			return IteratorByExpiry{tree, node, betweenByExpiry}
		}
		return tree.End()
	}
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (tree *ByExpiry) LowerBound(key common.TimePointSec) IteratorByExpiry {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByExpiryCompare(key, node.Key) > 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenByExpiry
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (tree *ByExpiry) UpperBound(key common.TimePointSec) IteratorByExpiry {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByExpiryCompare(key, node.Key) >= 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenByExpiry
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByExpiry) Remove(key common.TimePointSec) {
	if true {
		for lower := tree.LowerBound(key); lower.position != endByExpiry; {
			if ByExpiryCompare(lower.Key(), key) == 0 {
				node := lower.node
				lower.Next()
				tree.remove(node)
			} else {
				break
			}
		}
	} else {
		node := tree.lookup(key)
		tree.remove(node)
	}
}

func (tree *ByExpiry) put(key common.TimePointSec) (*ByExpiryNode, bool) {
	var insertedNode *ByExpiryNode
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		ByExpiryCompare(key, key)
		tree.Root = &ByExpiryNode{Key: key, color: redByExpiry}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		if true {
			for loop {
				compare := ByExpiryCompare(key, node.Key)
				switch {
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByExpiryNode{Key: key, color: redByExpiry}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare >= 0:
					if node.Right == nil {
						node.Right = &ByExpiryNode{Key: key, color: redByExpiry}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := ByExpiryCompare(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					return node, false
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByExpiryNode{Key: key, color: redByExpiry}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare > 0:
					if node.Right == nil {
						node.Right = &ByExpiryNode{Key: key, color: redByExpiry}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++

	return insertedNode, true
}

func (tree *ByExpiry) swapNode(node *ByExpiryNode, pred *ByExpiryNode) {
	if node == pred {
		return
	}

	tmp := ByExpiryNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

	pred.color = node.color
	node.color = tmp.color

	pred.Right = node.Right
	if pred.Right != nil {
		pred.Right.Parent = pred
	}
	node.Right = tmp.Right
	if node.Right != nil {
		node.Right.Parent = node
	}

	if pred.Parent == node {
		pred.Left = node
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}
		node.Parent = pred

	} else {
		pred.Left = node.Left
		if pred.Left != nil {
			pred.Left.Parent = pred
		}
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}

		node.Parent = tmp.Parent
		if node.Parent != nil {
			if node.Parent.Left == pred {
				node.Parent.Left = node
			} else {
				node.Parent.Right = node
			}
		} else {
			tree.Root = node
		}
	}
}

func (tree *ByExpiry) remove(node *ByExpiryNode) {
	var child *ByExpiryNode
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		tree.swapNode(node, pred)
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == blackByExpiry {
			node.color = nodeColorByExpiry(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = blackByExpiry
		}
	}
	tree.size--
}

func (tree *ByExpiry) lookup(key common.TimePointSec) *ByExpiryNode {
	node := tree.Root
	for node != nil {
		compare := ByExpiryCompare(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

// Empty returns true if tree does not contain any nodes
func (tree *ByExpiry) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *ByExpiry) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *ByExpiry) Keys() []common.TimePointSec {
	keys := make([]common.TimePointSec, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *ByExpiry) Values() []multi_index.TransactionState {
	values := make([]multi_index.TransactionState, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *ByExpiry) Left() *ByExpiryNode {
	var parent *ByExpiryNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *ByExpiry) Right() *ByExpiryNode {
	var parent *ByExpiryNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *ByExpiry) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *ByExpiry) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		outputByExpiry(tree.Root, "", true, &str)
	}
	return str
}

func (node *ByExpiryNode) String() string {
	if !node.color {
		return fmt.Sprintf("(%v,%v)", node.Key, "red")
	}
	return fmt.Sprintf("(%v)", node.Key)
}

func outputByExpiry(node *ByExpiryNode, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		outputByExpiry(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		outputByExpiry(node.Left, newPrefix, true, str)
	}
}

func (node *ByExpiryNode) grandparent() *ByExpiryNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *ByExpiryNode) uncle() *ByExpiryNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *ByExpiryNode) sibling() *ByExpiryNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (node *ByExpiryNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *ByExpiry) rotateLeft(node *ByExpiryNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *ByExpiry) rotateRight(node *ByExpiryNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *ByExpiry) replaceNode(old *ByExpiryNode, new *ByExpiryNode) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *ByExpiry) insertCase1(node *ByExpiryNode) {
	if node.Parent == nil {
		node.color = blackByExpiry
	} else {
		tree.insertCase2(node)
	}
}

func (tree *ByExpiry) insertCase2(node *ByExpiryNode) {
	if nodeColorByExpiry(node.Parent) == blackByExpiry {
		return
	}
	tree.insertCase3(node)
}

func (tree *ByExpiry) insertCase3(node *ByExpiryNode) {
	uncle := node.uncle()
	if nodeColorByExpiry(uncle) == redByExpiry {
		node.Parent.color = blackByExpiry
		uncle.color = blackByExpiry
		node.grandparent().color = redByExpiry
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *ByExpiry) insertCase4(node *ByExpiryNode) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *ByExpiry) insertCase5(node *ByExpiryNode) {
	node.Parent.color = blackByExpiry
	grandparent := node.grandparent()
	grandparent.color = redByExpiry
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *ByExpiryNode) maximumNode() *ByExpiryNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *ByExpiry) deleteCase1(node *ByExpiryNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *ByExpiry) deleteCase2(node *ByExpiryNode) {
	sibling := node.sibling()
	if nodeColorByExpiry(sibling) == redByExpiry {
		node.Parent.color = redByExpiry
		sibling.color = blackByExpiry
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *ByExpiry) deleteCase3(node *ByExpiryNode) {
	sibling := node.sibling()
	if nodeColorByExpiry(node.Parent) == blackByExpiry &&
		nodeColorByExpiry(sibling) == blackByExpiry &&
		nodeColorByExpiry(sibling.Left) == blackByExpiry &&
		nodeColorByExpiry(sibling.Right) == blackByExpiry {
		sibling.color = redByExpiry
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *ByExpiry) deleteCase4(node *ByExpiryNode) {
	sibling := node.sibling()
	if nodeColorByExpiry(node.Parent) == redByExpiry &&
		nodeColorByExpiry(sibling) == blackByExpiry &&
		nodeColorByExpiry(sibling.Left) == blackByExpiry &&
		nodeColorByExpiry(sibling.Right) == blackByExpiry {
		sibling.color = redByExpiry
		node.Parent.color = blackByExpiry
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *ByExpiry) deleteCase5(node *ByExpiryNode) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColorByExpiry(sibling) == blackByExpiry &&
		nodeColorByExpiry(sibling.Left) == redByExpiry &&
		nodeColorByExpiry(sibling.Right) == blackByExpiry {
		sibling.color = redByExpiry
		sibling.Left.color = blackByExpiry
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColorByExpiry(sibling) == blackByExpiry &&
		nodeColorByExpiry(sibling.Right) == redByExpiry &&
		nodeColorByExpiry(sibling.Left) == blackByExpiry {
		sibling.color = redByExpiry
		sibling.Right.color = blackByExpiry
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *ByExpiry) deleteCase6(node *ByExpiryNode) {
	sibling := node.sibling()
	sibling.color = nodeColorByExpiry(node.Parent)
	node.Parent.color = blackByExpiry
	if node == node.Parent.Left && nodeColorByExpiry(sibling.Right) == redByExpiry {
		sibling.Right.color = blackByExpiry
		tree.rotateLeft(node.Parent)
	} else if nodeColorByExpiry(sibling.Left) == redByExpiry {
		sibling.Left.color = blackByExpiry
		tree.rotateRight(node.Parent)
	}
}

func nodeColorByExpiry(node *ByExpiryNode) colorByExpiry {
	if node == nil {
		return blackByExpiry
	}
	return node.color
}

//////////////iterator////////////////

func (tree *ByExpiry) makeIterator(fn *TransactionStateIndexNode) IteratorByExpiry {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*ByExpiryNode); ok {
			return IteratorByExpiry{tree: tree, node: n, position: betweenByExpiry}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

// Iterator holding the iterator's state
type IteratorByExpiry struct {
	tree     *ByExpiry
	node     *ByExpiryNode
	position positionByExpiry
}

type positionByExpiry byte

const (
	beginByExpiry, betweenByExpiry, endByExpiry positionByExpiry = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *ByExpiry) Iterator() IteratorByExpiry {
	return IteratorByExpiry{tree: tree, node: nil, position: beginByExpiry}
}

func (tree *ByExpiry) Begin() IteratorByExpiry {
	itr := IteratorByExpiry{tree: tree, node: nil, position: beginByExpiry}
	itr.Next()
	return itr
}

func (tree *ByExpiry) End() IteratorByExpiry {
	return IteratorByExpiry{tree: tree, node: nil, position: endByExpiry}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *IteratorByExpiry) Next() bool {
	if iterator.position == endByExpiry {
		goto end
	}
	if iterator.position == beginByExpiry {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Left {
				goto between
			}
			node = iterator.node
		}
	}

end:
	iterator.node = nil
	iterator.position = endByExpiry
	return false

between:
	iterator.position = betweenByExpiry
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *IteratorByExpiry) Prev() bool {
	if iterator.position == beginByExpiry {
		goto begin
	}
	if iterator.position == endByExpiry {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Right {
				goto between
			}
			node = iterator.node
			//if iterator.tree.Comparator(node.Key, iterator.node.Key) >= 0 {
			//	goto between
			//}
		}
	}

begin:
	iterator.node = nil
	iterator.position = beginByExpiry
	return false

between:
	iterator.position = betweenByExpiry
	return true
}

func (iterator IteratorByExpiry) HasNext() bool {
	return iterator.position != endByExpiry
}

func (iterator *IteratorByExpiry) HasPrev() bool {
	return iterator.position != beginByExpiry
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorByExpiry) Value() multi_index.TransactionState {
	return *iterator.node.value()
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator IteratorByExpiry) Key() common.TimePointSec {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *IteratorByExpiry) Begin() {
	iterator.node = nil
	iterator.position = beginByExpiry
}

func (iterator IteratorByExpiry) IsBegin() bool {
	return iterator.position == beginByExpiry
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *IteratorByExpiry) End() {
	iterator.node = nil
	iterator.position = endByExpiry
}

func (iterator IteratorByExpiry) IsEnd() bool {
	return iterator.position == endByExpiry
}

// Delete remove the node which pointed by the iterator
// Modifies the state of the iterator.
func (iterator *IteratorByExpiry) Delete() {
	node := iterator.node
	//iterator.Prev()
	iterator.tree.remove(node)
}

func (tree *ByExpiry) inPlace(n *ByExpiryNode) bool {
	prev := IteratorByExpiry{tree, n, betweenByExpiry}
	next := IteratorByExpiry{tree, n, betweenByExpiry}
	prev.Prev()
	next.Next()

	var (
		prevResult int
		nextResult int
	)

	if prev.IsBegin() {
		prevResult = 1
	} else {
		prevResult = ByExpiryCompare(n.Key, prev.Key())
	}

	if next.IsEnd() {
		nextResult = -1
	} else {
		nextResult = ByExpiryCompare(n.Key, next.Key())
	}

	return (true && prevResult >= 0 && nextResult <= 0) ||
		(!true && prevResult > 0 && nextResult < 0)
}
