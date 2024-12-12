package sort

import (
	"sort"
)

// sort labels alphabeticaly by label name
type byKey []interface{}

func (k byKey) Len() int           { return len(k) / 2 }
func (k byKey) Less(i, j int) bool { return k[i*2].(string) < k[j*2].(string) }
func (k byKey) Swap(i, j int) {
	k[i*2], k[j*2] = k[j*2], k[i*2]
	k[i*2+1], k[j*2+1] = k[j*2+1], k[i*2+1]
}

func Uniq(labels []interface{}) []interface{} {
	if len(labels)%2 == 1 {
		labels = labels[:len(labels)-1]
	}

	if len(labels) > 2 {
		sort.Sort(byKey(labels))

		idx := 0
		for {
			if labels[idx] == labels[idx+2] {
				copy(labels[idx:], labels[idx+2:])
				labels = labels[:len(labels)-2]
			} else {
				idx += 2
			}
			if idx+2 >= len(labels) {
				break
			}
		}
	}
	return labels
}
