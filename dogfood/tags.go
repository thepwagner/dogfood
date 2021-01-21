package dogfood

import "math/rand"

type Tags map[string]string

func (t Tags) Tags() Tags { return t }

type HasTags interface {
	Tags() Tags
}

func NewMergedTags(tags ...HasTags) HasTags {
	return mergedTags(tags)
}

type mergedTags []HasTags

func (m mergedTags) Tags() Tags {
	ret := Tags{}
	for _, i := range m {
		if t := i.Tags(); t != nil {
			for k, v := range t {
				ret[k] = v
			}
		}
	}
	return ret
}

func NewWeightedTag(name string, weights map[string]int) HasTags {
	w := &weightedTag{name: name}
	for value, weight := range weights {
		for i := 0; i < weight; i++ {
			w.values = append(w.values, value)
		}
	}
	return w
}

type weightedTag struct {
	name   string
	values []string
}

func (w weightedTag) Tags() Tags {
	i := rand.Intn(len(w.values))
	value := w.values[i]
	return Tags{w.name: value}
}
