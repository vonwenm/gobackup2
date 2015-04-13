package main

type ArchivedFile struct {
	filename  string
	hash      string
	amazonId  string
	isDeleted bool
}

func (a *ArchivedFile) Filename() string {
	return a.filename
}

func (a *ArchivedFile) Hash() string {
	return a.hash
}

func (a *ArchivedFile) AmazonId() string {
	return a.amazonId
}

func (a *ArchivedFile) IsDeleted() bool {
	return a.isDeleted
}
