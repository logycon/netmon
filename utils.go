package main

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func IIF(cond bool, ifTrue string, ifFalse string) string {
	if cond {
		return ifTrue
	} else {
		return ifFalse
	}
}
