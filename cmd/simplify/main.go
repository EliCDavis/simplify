package main

import (
	"flag"
	"fmt"
	"log"
	"path"

	"github.com/fogleman/simplify"
)

var factor float64

func init() {
	flag.Float64Var(&factor, "f", 0.5, "percentage of faces in the output")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: simplify [-f FACTOR] input.stl output.stl")
		return
	}
	fmt.Printf("Loading %s\n", args[0])

	ext := path.Ext(args[0])
	var mesh *simplify.Mesh
	var err error
	if ext == ".stl" {
		mesh, err = simplify.LoadBinarySTL(args[0])
	} else if ext == ".obj" {
		mesh, err = simplify.LoadOBJ(args[0])
	}

	if err != nil {
		log.Fatal(err)
	}

	if mesh == nil {
		log.Fatalf("Unsupported filetype: %s\n", ext)
	}

	fmt.Printf("Input mesh contains %d faces\n", len(mesh.Triangles))
	fmt.Printf("Simplifying to %d%% of original...\n", int(factor*100))
	mesh = mesh.Simplify(factor)
	fmt.Printf("Output mesh contains %d faces\n", len(mesh.Triangles))
	fmt.Printf("Writing %s\n", args[1])
	mesh.SaveBinarySTL(args[1])
}
