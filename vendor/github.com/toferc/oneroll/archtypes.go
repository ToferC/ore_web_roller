package oneroll

// Archtype is a grouping of Sources, Permissions & Intrinsics that defines what powers a character can use
type Archtype struct {
	Sources     []*Source
	Permissions []*Permission
	Intrinsics  []*Intrinsic
}

// Source is a source of a Character's powers
type Source struct {
	Type        string
	Cost        int
	Description string
}

// Permission is the type of powers a Character can purchase
type Permission struct {
	Type        string
	Cost        int
	Description string
}

// Intrinsic is a modification from the human standard
type Intrinsic struct {
	Name        string
	Cost        int
	Description string
}
