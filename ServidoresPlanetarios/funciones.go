package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func agregarBase(nombreSector string, nombreBase string, cantidadSoldados string) {

	//Creamos y/o abrimos el .txt correspodiente
	file, err := os.OpenFile(nombreSector+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//Guardamos la informacion en SECTOR.txt
	file.WriteString(nombreSector + " " + nombreBase + " " + cantidadSoldados + "\n")

	//Cerramos el archivo
	file.Close()

}

func renombrarBase(nombreSector string, nombreBase string, nuevoNombre string) {
	var partes []string
	var linea string
	var lineaCambiar string
	var nuevaLinea string

	file, err := os.Open(nombreSector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	//Buscar en el archivo linea a reemplazar
	for scanner.Scan() {
		linea = scanner.Text()
		partes = strings.Split(linea, " ")
		if partes[1] == nombreBase {
			lineaCambiar = linea
			nuevaLinea = partes[0] + " " + nuevoNombre + " " + partes[2]
			break
		}
	}

	file.Close()

	input, err := ioutil.ReadFile(nombreSector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	output := bytes.Replace(input, []byte(lineaCambiar), []byte(nuevaLinea), -1)

	if err = ioutil.WriteFile(nombreSector+".txt", output, 0666); err != nil {
		log.Fatal(err)
	}

}

func actualizarValor(nombreSector string, nombreBase string, nuevoSoldados string) {
	var partes []string
	var linea string
	var lineaCambiar string
	var nuevaLinea string

	file, err := os.Open(nombreSector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	//Buscar en el archivo linea a reemplazar
	for scanner.Scan() {
		linea = scanner.Text()
		partes = strings.Split(linea, " ")
		if partes[1] == nombreBase {
			lineaCambiar = linea
			nuevaLinea = partes[0] + " " + partes[1] + " " + nuevoSoldados
			break
		}
	}

	file.Close()

	//Reemplazamos linea
	input, err := ioutil.ReadFile(nombreSector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	output := bytes.Replace(input, []byte(lineaCambiar), []byte(nuevaLinea), -1)

	if err = ioutil.WriteFile(nombreSector+".txt", output, 0666); err != nil {
		log.Fatal(err)
	}

}

func borrarBase(nombreSector string, nombreBase string) {
	var partes []string
	var linea string
	var lineaCambiar string

	file, err := os.Open(nombreSector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	//Buscar en el archivo linea a reemplazar
	for scanner.Scan() {
		linea = scanner.Text()
		partes = strings.Split(linea, " ")
		if partes[1] == nombreBase {
			lineaCambiar = linea
			break
		}
	}

	file.Close()

	//Reemplazamos linea
	input, err := ioutil.ReadFile(nombreSector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	output := bytes.Replace(input, []byte(lineaCambiar+"\n"), []byte(""), -1)

	if err = ioutil.WriteFile(nombreSector+".txt", output, 0666); err != nil {
		log.Fatal(err)
	}
}

func aplicarCambios(nombreArchivoLogs string, nuevoContenido string, cantidadInicial int) {

	var partes []string

	d1 := []byte(nuevoContenido)
	err := os.WriteFile(nombreArchivoLogs, d1, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(nombreArchivoLogs)
	if err != nil {
		log.Fatal(err)
	}

	i := 1
	scanner := bufio.NewScanner(file)

	//Buscar en el archivo lineas nuevas para aplicar los cambios
	for scanner.Scan() {
		if cantidadInicial < i { // Cuando ocurre esto, significa que estamos en lineas nuevas
			partes = strings.Split(scanner.Text(), " ")
			if partes[0] == "AgregarBase" {
				agregarBase(partes[1], partes[2], partes[3])

			} else if partes[0] == "RenombrarBase" {
				renombrarBase(partes[1], partes[2], partes[3])

			} else if partes[0] == "ActualizarValor" {
				actualizarValor(partes[1], partes[2], partes[3])
			} else { //Este else representa BorrarBase

			}
		}

		i += 1
	}
	file.Close()

}
