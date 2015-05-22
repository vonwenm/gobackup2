package main

type filechecker struct {
	archive *archive
}

func NewFileChecker(archive *archive) (*filechecker, error) {
	return &filechecker{
		archive: archive,
	}, nil
}
