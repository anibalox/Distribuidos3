package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func escribirArchivo(linea string) {

	file, err := os.OpenFile("HakunaMatata.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	//Guardamos la informacion en SECTOR.txt
	file.WriteString(linea + "\n")

	file.Close()
}

func reemplazarNombreBase(base string, nuevoNombre string) {

	var partes []string
	var linea string
	var lineaCambiar string
	var nuevaLinea string

	file, _ := os.Open("HakunaMatata.txt")
	scanner := bufio.NewScanner(file)
	//Buscar en el archivo linea a reemplazar
	for scanner.Scan() {
		linea = scanner.Text()
		partes = strings.Split(linea, " ")
		if partes[1] == base {
			lineaCambiar = linea
			nuevaLinea = partes[0] + " " + nuevoNombre + " " + partes[2]
			break
		}
	}

	file.Close()

	input, _ := ioutil.ReadFile("HakunaMatata.txt")

	output := bytes.Replace(input, []byte(lineaCambiar), []byte(nuevaLinea), -1)

	ioutil.WriteFile("HakunaMatata.txt", output, 0666)

}

func main() {

	escribirArchivo("CostaEnredada Campamento1 4000")
	escribirArchivo("CostaEnredada Camp 40550")
	escribirArchivo("CostaEnredada Bolivia 405501110")
	escribirArchivo("CostaEnredada Ruperto 4111000")
	escribirArchivo("CostaEnredada Lamama 408282900")
	escribirArchivo("CostaEnredada Koala 4002220")

	reemplazarNombreBase("Bolivia", "Chile")

}
