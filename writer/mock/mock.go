package mock

// Writer mocks the Writer to record files written during unit testing
type Writer struct {
	files map[string]map[string]string
}

// WriteFiles saves files to the specified directory
func (w *Writer) WriteFiles(directory string, files map[string]string) (err error) {
	for name, content := range files {
		_ = w.WriteFile(directory, name, content)
	}
	return
}

// WriteFile saves one file to the specified directory
func (w *Writer) WriteFile(directory, file, content string) (err error) {
	if w.files == nil {
		w.files = make(map[string]map[string]string)
	}

	if _, ok := w.files[directory]; ok == false {
		w.files[directory] = make(map[string]string)
	}

	w.files[directory][file] = content

	return
}

// GetFile returns the content of the file, written be WriteFiles or WriteFile. Returns false if the file does not exist.
func (w *Writer) GetFile(directory, file string) (content string, ok bool) {
	if w.files != nil {
		var dir map[string]string

		if dir, ok = w.files[directory]; ok {
			content, ok = dir[file]
		}
	}

	return
}
