package utils

import "fmt"
func firstN(s string, n int) string {
    i := 0
    for j := range s {
        if i == n {
            return s[:j]
        }
        i++
    }
    return s
}

func Ellipsis(s string, n int) string {
    if len(s) < n {
        return s
    }
    return fmt.Sprintf("%s%s", firstN(s, n - 1), "...")
}
