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

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	pb "biomorph/db"

	"bytes"
	"encoding/json"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	port = ":50051"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "bio.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&CreatureModel{})
}

type JsonModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Json      []byte
	ID        uint64 `gorm:"primary_key auto_increment"`
}

func (j *JsonModel) Decode(v interface{}) error {
	return json.NewDecoder(bytes.NewReader(j.Json)).Decode(&v)
}

func (j *JsonModel) Encode(v interface{}) error {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(v)
	if err != nil {
		return err
	}
	j.Json = b.Bytes()
	return nil
}

type CreatureModel struct {
	JsonModel
}

type values map[string]float64

type value_map struct {
	VMap    values
	Parents []uint64
}

type server struct{}

func (s *server) GetCreature(ctx context.Context, in *pb.GetCreatureRequest) (*pb.GetCreatureReply, error) {
	r := pb.GetCreatureReply{}
	var m CreatureModel
	db.First(&m, in.GetId())
	if len(m.Json) == 0 {
		return &r, errors.New(fmt.Sprintf("Could not find creature ID %d", in.GetId()))
	}
	var vm value_map
	e := m.Decode(&vm)
	if e != nil {
		return &r, e
	}

	r.Parents = vm.Parents
	r.Values = vm.VMap
	return &r, nil
}

func (s *server) SaveCreature(ctx context.Context, in *pb.SaveCreatureRequest) (*pb.SaveCreatureReply, error) {
	r := pb.SaveCreatureReply{}
	m := CreatureModel{}
	e := m.Encode(value_map{in.GetValues(), in.GetParents()})
	if e != nil {
		return &r, e
	}
	db.Create(&m)
	r.Id = m.ID
	return &r, nil
}

func main() {
	defer db.Close()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDbServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
