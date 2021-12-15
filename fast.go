package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	// "log"
)


// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	lind := 0
	var line []byte
	var user = make(map[string]interface{})
	var isAndroid uint8
	var isMSIE uint8
	var indAndroid int
	var indMSIE int
	var email string
	for scanner.Scan() {
		line = scanner.Bytes()
		err := json.Unmarshal(line, &user)
		if err != nil {
			panic(err)
		}

		isAndroid = 0
		isMSIE = 0

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

			indAndroid = strings.Index(browser, "Android")
			indMSIE = strings.Index(browser, "MSIE")

			if indAndroid >= 0 {
				isAndroid++
			}
			if indMSIE >= 0 {
				isMSIE++
			}

			if indAndroid >= 0 || indMSIE >= 0 {
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid > 0 && isMSIE > 0) {
			lind++
			continue
		}

		email = strings.Replace(user["email"].(string), "@", " [at] ", -1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", lind, user["name"], email)

		lind++
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}

	_, err = fmt.Fprintln(out, "found users:\n"+foundUsers)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
	if err != nil {
		panic(err)
	}
}