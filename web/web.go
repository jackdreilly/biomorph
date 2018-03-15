package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackdreilly/biomorph"
	pb "github.com/jackdreilly/biomorph/db"

	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	n_images = 30
	address  = "localhost:50051"
)

var (
	client pb.DbClient
	conn   *grpc.ClientConn
	logger *log.Logger
)

func init() {
	// Set up a connection to the server.
	log.SetOutput(os.Stderr)
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client = pb.NewDbClient(conn)

	// Sets your Google Cloud Platform project ID.
	projectID := "quiklyrics-go"

	// Creates a client.
	ctx := context.Background()
	c, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name of the log to write to.
	logName := "biomorph"

	// Selects the log to write to.
	logger = c.Logger(logName).StandardLogger(logging.Info)
}

type Image struct {
	Bytes string `json:"bytes"`
	Id    uint64 `json:"id"`
}

type Response struct {
	Images []Image `json:"images"`
	Gif    string  `json:"gif"`
}

type value_map struct {
	values  map[string]float64
	parents []uint64
}

func log_err(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func AddCreature(v *value_map) uint64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SaveCreature(ctx, &pb.SaveCreatureRequest{Values: v.values, Parents: v.parents})
	log_err(err)
	return r.GetId()
}

func Parents(id uint64) []uint64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetCreature(ctx, &pb.GetCreatureRequest{Id: id})
	log_err(err)
	return r.GetParents()
}

func HistoryImages(id uint64) []Image {
	parents := Parents(id)
	parents = append(parents, id)
	images := make([]Image, len(parents))
	for i, cid := range parents {
		nc, _ := GetCreature(cid)
		img := biomorph.DrawTreeCreature(nc)
		var buff bytes.Buffer
		png.Encode(&buff, img)
		images[i].Bytes = base64.StdEncoding.EncodeToString(buff.Bytes())
		images[i].Id = cid
	}
	return images
}

type my_quant struct{}

func (m *my_quant) Quantize(p color.Palette, _ image.Image) color.Palette {
	return []color.Color{
		color.White,
		color.Black,
	}
}

func Gif(images []Image) string {
	outGif := &gif.GIF{}
	for _, i := range images {
		b, _ := base64.StdEncoding.DecodeString(i.Bytes)
		im, _ := png.Decode(bytes.NewReader(b))
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		gif.Encode(w, im, &gif.Options{NumColors: 2, Quantizer: &my_quant{}})
		gim, _ := gif.Decode(bufio.NewReader(&buf))
		outGif.Image = append(outGif.Image, gim.(*image.Paletted))
		outGif.Delay = append(outGif.Delay, 0)
	}
	var buff bytes.Buffer
	w := bufio.NewWriter(&buff)
	gif.EncodeAll(w, outGif)
	return base64.StdEncoding.EncodeToString(buff.Bytes())
}

func NewCreature() (*biomorph.Creature, uint64) {
	c := biomorph.NewCreature(biomorph.NewTreeSpecies())
	return c, AddCreature(&value_map{c.ValuesMap(), []uint64{}})
}

func GetCreature(id uint64) (*biomorph.Creature, []uint64) {
	c := biomorph.NewCreature(biomorph.NewTreeSpecies())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetCreature(ctx, &pb.GetCreatureRequest{Id: id})
	log_err(err)
	c.SetValuesFromMap(r.GetValues())
	return c, r.GetParents()
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	c, id := NewCreature()
	WriteImagesOut(id, c, []uint64{}, w)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	creature, parents := GetCreature(uint64(id))
	WriteHistoryOut(uint64(id), creature, parents, w)
}

func MutateImage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	creature, parents := GetCreature(uint64(id))
	WriteImagesOut(uint64(id), creature, parents, w)
}

func WriteImagesOut(id uint64, creature *biomorph.Creature, parents []uint64, w http.ResponseWriter) {
	parents = append(parents, id)
	response := Response{Images: make([]Image, n_images)}
	var vm value_map
	vm.parents = parents
	for i := 0; i < n_images; i++ {
		nc := biomorph.MutateCreature(creature)
		img := biomorph.DrawTreeCreature(nc)
		var buff bytes.Buffer
		png.Encode(&buff, img)
		response.Images[i].Bytes = base64.StdEncoding.EncodeToString(buff.Bytes())
		vm.values = nc.ValuesMap()
		nid := AddCreature(&vm)
		response.Images[i].Id = nid
	}
	json.NewEncoder(w).Encode(response)
}

func WriteHistoryOut(id uint64, creature *biomorph.Creature, parents []uint64, w http.ResponseWriter) {
	parents = append(parents, id)
	images := HistoryImages(id)
	json.NewEncoder(w).Encode(&Response{images, Gif(images)})
}

func main() {
	defer conn.Close()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/get_images", GetImages)
	http.HandleFunc("/get_image", GetImage)
	http.HandleFunc("/mutate_image", MutateImage)

	http.HandleFunc("/choose_image", func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.URL.Query().Get("id"))
	})

	http.ListenAndServe(":8080", nil)
}
