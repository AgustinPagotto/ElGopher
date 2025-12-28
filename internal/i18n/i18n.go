package i18n

type Translator struct {
	Messages map[string]string
}

func (t Translator) T(key string) string {
	text := t.Messages[key]
	if text == "" {
		return "error"
	}
	return text
}
