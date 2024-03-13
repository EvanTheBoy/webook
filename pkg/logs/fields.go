package logs

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Error(err error) Field {
	return Field{
		Key:   "err",
		Value: err,
	}
}
