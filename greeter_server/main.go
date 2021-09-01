/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "linkCutterTask/helloworld"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	//pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var (
	dbMap        = make(map[string]string)
	compareSlice []string
)

func RandomizeString(link string) string {
	alphabet := "a1b2c3d4f5z6x7c8v9b0mnbhj"
	alphabetHelp := "abc1def2ghi3jkl4mnop5qrs6tuv7wxy8zAB9CDEF0GHIJKLMNOPQRSTUVWXYZ5qrs6tuv7wxy8zAB9CDEF0GHEF0GHIJKLMNO4mnop5qrs6tdef2ghi36tuv7wxy8LMNO4mn"
	alphabetDig := "0123456789"
	alphabetLower := "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	originalLink := link
	originalLinkLen := utf8.RuneCountInString(originalLink)
	link = strings.ToLower(link)
	//наличие https/http
	link = strings.ReplaceAll(link, "http://", "")
	link = strings.ReplaceAll(link, "https://", "")
	if utf8.RuneCountInString(link) < 9 {
		link += alphabet[utf8.RuneCountInString(link):9]

	} else if utf8.RuneCountInString(link) > 9 {
		link = link[0:9]
	} else {
		link = link[0:9]
	}
	rLink := []rune(link)
	//преобразуем ссылки
	//здесь баг, если использовать @@, то есть два спец символа подряд и более, то ссылка становится не уникальной
	for i, j := 0, len(rLink)-1; i < j; i, j = i+1, j-1 {
		rLink[i], rLink[j] = rLink[j], rLink[i]
		if rLink[i]%2 == 0 {
			rLink[i] = unicode.ToUpper(rLink[i])
		}
	}
	//преобразование ссылки без спец символов, если в неё поместили спец.символ
	for i, r := range rLink {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			rLink[i] = rune(alphabet[i])
		} else {
			rLink[i] = r
		}
	}
	//для сравнения, что такая линка не юзается
	compareSlice = append(compareSlice, string(rLink))
	//проверка перед отправкой дальше, что соответствует шаблону
	//длина 10 символов
	//содержит буквы up down
	//содержит цифры
	//сожержит _
	//уникальна
	if value, ok := dbMap[string(rLink)+"_"]; ok {
		log.Println("Ссылка уже существует в MAP!")
		if value != originalLink[0:originalLinkLen] {
			log.Println("Оригинальные ссылки разные", value, originalLink[0:originalLinkLen])
			//пишем логику доработки изменённой ссылки
			tmp := []rune(link[0:1])
			for i, _ := range rLink {
				rLink[i] = rune(alphabetHelp[rune(i)+tmp[0]])
			}

		} else {
			log.Println("Оригинальные ссылки одинаковые", value, originalLink[0:originalLinkLen])
		}

	}
	//проверка на количество апперов и лоуверов
	countDig := 0
	countUpper := 0
	countLower := 0
	for _, r := range rLink {
		if unicode.IsDigit(r) {
			countDig += 1
		} else if unicode.IsUpper(r) {
			countUpper += 1
		} else if unicode.IsLower(r) {
			countLower += 1
		}
	}

	if countDig == 0 {
		rLink[0] = rune(alphabetDig[originalLinkLen%10])
		rLink[8] = rune(alphabetDig[originalLinkLen%5])
	}
	if countUpper == 0 {
		rLink[1] = rune(alphabetUpper[originalLinkLen%10])
		rLink[7] = rune(alphabetUpper[originalLinkLen%5])
	}
	if countLower == 0 {
		rLink[2] = rune(alphabetLower[originalLinkLen%10])
		rLink[6] = rune(alphabetLower[originalLinkLen%5])
	}

	for _, r := range rLink {
		if unicode.IsDigit(r) {
			countDig += 1
		}
	}
	log.Printf("countDig: %d\tcountUpper: %d\tcountLower: %d", countDig, countUpper, countLower)
	log.Println("len url short:", len(rLink), rLink)

	return string(rLink) + "_"
}

func CuttingLink(link string) string {
	//создаём укороченную линку и вносим в мап
	linkOriginal := link
	link = RandomizeString(link)
	dbMap[link] = linkOriginal
	for key, value := range dbMap {
		log.Printf("DBMap  | Short [%s]: Orig [%s]\n", key, value)
	}
	log.Println("DBMap len:", len(dbMap))
	return link
}

// SayHello implements helloworld.GreeterServer
func (s *server) Create(ctx context.Context, in *pb.LinkRequest) (*pb.LinkReply, error) {
	tmp := in.GetName()
	log.Printf("Server | Received from client origLink: %v", in.GetName())
	tmp = CuttingLink(tmp)
	return &pb.LinkReply{Message: "Server | Client get short link: " + tmp}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *server) Get(ctx context.Context, in *pb.LinkRequest) (*pb.LinkReply, error) {
	tmp := in.GetName()
	_, ok := dbMap[tmp]
	if utf8.RuneCountInString(tmp) < 10 {
		log.Println("Длина ссылки меньше 10 символов")
		return &pb.LinkReply{Message: "Длина ссылки меньше 10"}, nil
	} else if !ok {
		log.Println("Укороченная ссылка не найдена", ok)
		return &pb.LinkReply{Message: "Укороченная ссылка не найдена"}, nil
	} else {
		log.Printf("Received from client short link: %v", in.GetName())
		tmp = dbMap[tmp]
	}
	return &pb.LinkReply{Message: "Server | Original Link: " + tmp}, nil
}

func main() {
	rand.Seed(time.Now().Unix())
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	//для тестов
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
