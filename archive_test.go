package main

type Archiver interface {
	AddFile(*ArchivedFile) error
	FindFileByFilename(string) *ArchivedFile
	FindAmazonIdByHash(string) *string
	DeleteFile(string, string) error
}

type mockArchiver struct {
	files []*ArchivedFile
}

func (m *mockArchiver) AddFile(a *ArchivedFile) error {
	m.files = append(m.files, a)
	return nil
}

func (m *mockArchiver) FindFileByFilename(filename string) *ArchivedFile {
	for _, file := range m.files {
		if file.Filename() == filename {
			return file
		}
	}
	return nil
}

func (m *mockArchiver) FindAmazonIdByHash(hash string) *string {
	for _, file := range m.files {
		if file.Hash() == hash {
			amazonId := file.AmazonId()
			return &amazonId
		}
	}
	return nil
}

func (m *mockArchiver) DeleteFile(hash, filename string) error {
	for i, file := range m.files {
		if file.Hash() == hash && file.Filename() == filename {
			m.files = append(m.files[:i], m.files[i+1:]...)
		}
	}
	return nil
}
