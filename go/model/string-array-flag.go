package model

type StringArrayFlags []string

func (saf *StringArrayFlags) String() string {
	out := "["
	for _, s := range *saf {
		out = out + s + ","
	}
	return out[:len(out)-1] + "]"
}

func (saf *StringArrayFlags) Set(value string) error {
	*saf = append(*saf, value)
	return nil
}
