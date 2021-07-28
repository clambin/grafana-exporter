package mock

type Writer struct {
	files map[string]map[string]string
}

func (w *Writer) WriteFiles(directory string, files map[string]string) (err error) {
	for name, content := range files {
		_ = w.WriteFile(directory, name, content)
	}
	return
}

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

func (w *Writer) GetFile(directory, file string) (content string, ok bool) {
	if w.files != nil {
		var dir map[string]string

		if dir, ok = w.files[directory]; ok {
			content, ok = dir[file]
		}
	}

	return
}
