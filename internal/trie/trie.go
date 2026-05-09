package trie


type node struct {
	children map[rune]*node
	isWord   bool
}

type Trie struct {
	root *node
}

func newNode() *node {
	return &node{children: make(map[rune]*node)}
}

func NewTrie() *Trie {
	return &Trie{root: newNode()}
}

func (t *Trie) Insert(w string) {
	curr := t.root
	for _, r := range w {
		next, ok := curr.children[r]
		if !ok {
			next = newNode()
			curr.children[r] = next
		}
		curr = next
	}
	curr.isWord = true
}

func (t *Trie) walk(s string) *node {
	curr := t.root
	for _, r := range s {
		curr = curr.children[r]
		if curr == nil {
			return nil
		}
	}
	return curr
}

func (t *Trie) IsWord(w string) bool {
	final := t.walk(w)
	return final != nil && final.isWord
}

func (t *Trie) collect(n *node, buf []rune, res *[]string) {
	if n.isWord {
		*res = append(*res, string(buf))
	}

	for r, next := range n.children {
		buf = append(buf, r)
		t.collect(next, buf, res)
		buf = buf[:len(buf)-1]
	}
}

func (t *Trie) WithPrefix(p string) []string {
	n := t.walk(p)
	if n == nil {
		return nil
	}

	var res []string
	t.collect(n, []rune(p), &res)
	return res
}
