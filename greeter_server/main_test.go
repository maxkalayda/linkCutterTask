package main

import (
	"fmt"
	"testing"
	"unicode"
)

func Test(t *testing.T) {
	//тестовые данные
	countDig := 0
	countUpper := 0
	countLower := 0
	countUnder := 0
	testData := []string{
		"https://maxkalayda.com",
		"http://maxkalayda.com",
		"https://maxkalayda.com/profile",
		"https://=-(*&%~±3",
		"@@@@@@",
		"!!!!!!",
		"YyYyYyYyY",
		"y$an#de@x.%ru",
		"  ",
		"",
		"https://m.com",
		"http://123.ru",
		"@#",
		"@!",
		"23456.ru",
		"https://ya.ru",
		"aaa",
		"bbb",
		"ccc",
		"AAA",
		"CCC",
		"BBB",
		"https://developers.google.com/protocol-buffers/docs/reference/go-generated",
		"https://grpc.io/docs/languages/go/quickstart/",
		"https://gitlab.com/users/sign_in?__cf_chl_jschl_tk__=pmd_m.i_sJo3lYAtvPP8q17OcGq0lPfHPQD3OGhswCgInuM-1630398294-0-gqNtZGzNAiWjcnBszQl9",
	}
	//тестируемый код
	for _, r := range testData {
		result := RandomizeString(r)
		for _, r := range result {
			if unicode.IsDigit(r) {
				countDig += 1
			} else if unicode.IsUpper(r) {
				countUpper += 1
			} else if unicode.IsLower(r) {
				countLower += 1
			} else if r == '_' {
				countUnder += 1
			}
		}
		if countDig == 0 || countLower == 0 || countUpper == 0 {
			t.Errorf("Тест не пройден! Ожидаем CountDig, CountUpper, CountLower > 0, получили"+
				"CountDig = %d, CountUpper = %d, CountLower = %d", countDig, countUpper, countLower)
		} else {
			fmt.Printf("CountDig = %d, CountUpper = %d, CountLower = %d, countUnder = %d\n", countDig, countUpper, countLower, countUnder)
			countDig = 0
			countUpper = 0
			countLower = 0
			countUnder = 0
		}
	}
}
