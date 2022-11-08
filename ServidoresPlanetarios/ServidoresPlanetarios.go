package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	pb "Tarea3/Proto"
)

type server struct {
	pb.UnimplementedServidorPlanetarioServer
}

var reloj [3]int = [3]int{0, 0, 0}
var nroServidor int

func (s *server) AgregarBase(ctx context.Context, req *pb.Peticion) (*pb.Reloj, error) {

	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre
	cantidadSoldados := req.Valor

	//Creamos y/o abrimos el .txt correspodiente
	file, err := os.OpenFile(nombreSector+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//Guardamos la informacion en SECTOR.txt
	file.WriteString(nombreSector + " " + nombreBase + " " + cantidadSoldados + "\n")

	//Cerramos el archivo
	file.Close()

	file, err = os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
	if err != nil {
		log.Fatal(err)
	}

	//Guardamos la informacion en el Log
	file.WriteString("AgregarBase " + nombreSector + " " + nombreBase + " " + cantidadSoldados + "\n")

	//Cerramos el archivo
	file.Close()

	reloj[nroServidor] += 1

	stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])

	return &pb.Reloj{RelojServidor: stringReloj}, nil
}

//Hay que ver si esto funciona
func (s *server) RenombrarBase(ctx context.Context, req *pb.Peticion) (*pb.Reloj, error) {

	var partes []string
	var linea string
	var lineaCambiar string
	var nuevaLinea string

	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre
	nuevoNombre := req.Valor

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

	file, err = os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
	if err != nil {
		log.Fatal(err)
	}

	//Guardamos la informacion en el Log
	file.WriteString("RenombrarBase " + nombreSector + " " + nombreBase + " " + nuevoNombre + "\n")

	//Cerramos el archivo
	file.Close()

	reloj[nroServidor] += 1

	stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])

	return &pb.Reloj{RelojServidor: stringReloj}, nil
}
func (s *server) ActualizarValor(ctx context.Context, req *pb.Peticion) (*pb.Reloj, error) {
	var partes []string
	var linea string
	var lineaCambiar string
	var nuevaLinea string

	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre
	nuevoSoldados := req.Valor

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

	file, err = os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
	if err != nil {
		log.Fatal(err)
	}

	//Guardamos la informacion en el Log
	file.WriteString("ActualizarValor " + nombreSector + " " + nombreBase + " " + nuevoSoldados + "\n")

	//Cerramos el archivo
	file.Close()

	reloj[nroServidor] += 1

	stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])

	return &pb.Reloj{RelojServidor: stringReloj}, nil
}

func (s *server) BorrarBase(ctx context.Context, req *pb.Peticion) (*pb.Reloj, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BorrarBase not implemented")
}

func main() {

	nroServidor, _ = strconv.Atoi(os.Args[1]) //1: Tierra , 2: Titan, 3: Marte Le puedes dar un numero Felipe?
	port := ":50051"                          // Puerto de DataNode

	listner, err := net.Listen("tcp", port)

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	pb.RegisterServidorPlanetarioServer(serv, &server{})

	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}

}
