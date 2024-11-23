package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	users := make([]User, 0, 1000)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		line = bytes.TrimSpace(line)

		user := User{}
		err = user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	seenBrowsers := make(map[string]struct{})
	uniqueBrowsers := 0
	sb := strings.Builder{}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers
		for _, browser := range browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
		}
		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", -1)
		sb.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}

	fmt.Fprintln(out, "found users:\n"+sb.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))

}
