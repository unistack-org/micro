package metadata

import "sort"

type Iterator struct {
	md   Metadata
	keys []string
	cur  int
	cnt  int
}

// Next advances the iterator to the next element.
func (iter *Iterator) Next(k *string, v *[]string) bool {
	if iter.cur+1 > iter.cnt {
		return false
	}

	if k != nil && v != nil {
		*k = iter.keys[iter.cur]
		vv := iter.md[*k]
		*v = make([]string, len(vv))
		copy(*v, vv)
		iter.cur++
	}
	return true
}

// Iterator returns an iterator for iterating over metadata in sorted order.
func (md Metadata) Iterator() *Iterator {
	iter := &Iterator{md: md, cnt: len(md)}
	iter.keys = make([]string, 0, iter.cnt)
	for k := range md {
		iter.keys = append(iter.keys, k)
	}
	sort.Strings(iter.keys)
	return iter
}
