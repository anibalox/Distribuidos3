package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	pb "Tarea3/Proto"

	"google.golang.org/grpc"
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

	hostNameNode := "localhost" // Host de NameNode
	port := ":50051"            // Puerto de NameNode
	connS, err := grpc.Dial(hostNameNode+port, grpc.WithInsecure())

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serviceCliente := pb.NewServidorPlanetarioClient(connS)

	// res, err := serviceCliente.AgregarBase(context.Background(),
	// 	&pb.Peticion{Base: &pb.DatosBase{Sector: "Calama", Nombre: "Base1"},
	// 		Valor: "1500"})

	// res, err = serviceCliente.AgregarBase(context.Background(),
	// 	&pb.Peticion{Base: &pb.DatosBase{Sector: "Calama", Nombre: "Base2"},
	// 		Valor: "153300"})

	// res, err = serviceCliente.AgregarBase(context.Background(),
	// 	&pb.Peticion{Base: &pb.DatosBase{Sector: "Copiapo", Nombre: "Base1"},
	// 		Valor: "17000"})

	// res, err = serviceCliente.RenombrarBase(context.Background(),
	// 	&pb.Peticion{Base: &pb.DatosBase{Sector: "Copiapo", Nombre: "Base1"},
	// 		Valor: "NuevaCopiapoBase1"})

	// res, err = serviceCliente.ActualizarValor(context.Background(),
	// 	&pb.Peticion{Base: &pb.DatosBase{Sector: "Calama", Nombre: "Base1"},
	// 		Valor: "50000"})

	// res, err := serviceCliente.BorrarBase(context.Background(),
	// 	&pb.Peticion{Base: &pb.DatosBase{Sector: "Calama", Nombre: "Base1"},
	// 		Valor: ""})

	res, err := serviceCliente.GetSoldados(context.Background(),
		&pb.DatosBase{Sector: "Calama", Nombre: "Base2"})

	fmt.Println(res.CantidadSoldados)
	fmt.Println(res.NombreServidor)
	fmt.Println(res.RelojServidor)

	connS.Close()

	// contenido, err := os.ReadFile("NewLogs.txt")

	// serviceCliente.Merge(context.Background(), &pb.MensajeMerge{
	// 	RelojServidor: &pb.Reloj{RelojServidor: "3-5-0"},
	// 	Logs:          string(contenido),
	// })

}
