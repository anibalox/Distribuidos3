package main

import (
	"bufio"
	"fmt"
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

var numberToIp = map[int]string{
	0: "10003",    //Colocar IP de Tierra
	1: "100002",   // IP de Titan
	2: "23231231", // IP Marte
}

func (s *server) AgregarBase(ctx context.Context, req *pb.Peticion) (*pb.Reloj, error) {

	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre
	cantidadSoldados := req.Valor

	agregarBase(nombreSector, nombreBase, cantidadSoldados)

	file, err := os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
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

func (s *server) RenombrarBase(ctx context.Context, req *pb.Peticion) (*pb.Reloj, error) {

	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre
	nuevoNombre := req.Valor

	renombrarBase(nombreSector, nombreBase, nuevoNombre)

	file, err := os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
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

	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre
	nuevoSoldados := req.Valor

	actualizarValor(nombreSector, nombreBase, nuevoSoldados)

	file, err := os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
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
	nombreSector := req.Base.Sector
	nombreBase := req.Base.Nombre

	borrarBase(nombreSector, nombreBase)

	file, err := os.OpenFile("LogRegistro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el Log de registros para anadir el cambio
	if err != nil {
		log.Fatal(err)
	}

	//Guardamos la informacion en el Log
	file.WriteString("BorrarBase " + nombreSector + " " + nombreBase + "\n")

	//Cerramos el archivo
	file.Close()

	reloj[nroServidor] += 1

	stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])

	return &pb.Reloj{RelojServidor: stringReloj}, nil
}

func traducirReloj(reloj string) [3]int {

	var nuevoReloj [3]int

	datos := strings.Split(reloj, "-")

	for i := 0; i < 3; i++ {
		nuevoReloj[i], _ = strconv.Atoi(datos[i])
	}

	return nuevoReloj
}

func (s *server) GetSoldados(ctx context.Context, req *pb.DatosBase) (*pb.SoldadosBase, error) {

	var partes []string
	var linea string
	cantidadSoldados := -1

	file, err := os.Open(req.Sector + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	//Buscar en el archivo linea a reemplazar
	for scanner.Scan() {
		linea = scanner.Text()
		partes = strings.Split(linea, " ")
		if partes[1] == req.Nombre {
			cantidadSoldados, _ = strconv.Atoi(partes[2])
			break
		}
	}

	file.Close()
	stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])

	return &pb.SoldadosBase{CantidadSoldados: int32(cantidadSoldados),
		RelojServidor:  stringReloj,
		NombreServidor: strconv.Itoa(nroServidor)}, nil
}

func (s *server) IniciarMerge(context.Context, *pb.MensajeSimple) (*pb.MensajeSimple, error) {

	//Este codigo obtiene la informacion de los otros servidores y actualiza la propia
	for nroServ, direccionServidorPlanetario := range numberToIp { //Se itera sobre numberToIp
		if nroServ != nroServidor {
			connS, err := grpc.Dial(direccionServidorPlanetario, grpc.WithInsecure())
			if err != nil {
				panic("cannot create tcp connection" + err.Error())
			}
			serviceCliente := pb.NewServidorPlanetarioClient(connS)

			res, _ := serviceCliente.Merge(context.Background(), &pb.MensajeMerge{
				RelojServidor: &pb.Reloj{RelojServidor: ""},
				Logs:          "",
			})

			//Calculamos cuantos logs inicialmente habia
			cantidadInicialLogs := reloj[nroServ]

			//Actualizamos el reloj
			reloj[nroServ] = traducirReloj(res.RelojServidor.RelojServidor)[nroServ]

			//Aplicamos los cambios a el servidor actual
			aplicarCambios("LogRegistro.txt", res.Logs, cantidadInicialLogs)

			connS.Close()
		}
	}

	//Este codigo propaga los cambios
	for nroServ, direccionServidorPlanetario := range numberToIp { //Se itera sobre numberToIp
		if nroServ != nroServidor {
			connS, err := grpc.Dial(direccionServidorPlanetario, grpc.WithInsecure())
			if err != nil {
				panic("cannot create tcp connection" + err.Error())
			}
			serviceCliente := pb.NewServidorPlanetarioClient(connS)

			stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])
			contenido, err := os.ReadFile("LogRegistro.txt")

			serviceCliente.Merge(context.Background(), &pb.MensajeMerge{
				RelojServidor: &pb.Reloj{RelojServidor: stringReloj},
				Logs:          string(contenido),
			})

			connS.Close()
		}
	}

	return nil, status.Errorf(codes.Unimplemented, "method IniciarMerge not implemented")
}

func (s *server) Merge(ctx context.Context, req *pb.MensajeMerge) (*pb.MensajeMerge, error) {

	if req.Logs == "" { //Si Logs es vacio, significa que se le pide los logs y los relojes

		//Se crea reloj y contenido para enviar los datos al servidor dominante
		stringReloj := strconv.Itoa(reloj[0]) + "-" + strconv.Itoa(reloj[1]) + "-" + strconv.Itoa(reloj[2])
		contenido, err := os.ReadFile("LogRegistro.txt")

		if err != nil {
			log.Fatal(err)
		}
		return &pb.MensajeMerge{RelojServidor: &pb.Reloj{RelojServidor: stringReloj}, Logs: string(contenido)}, nil

	} else { // Si no lo es, significa que se tiene que actualizar el propio Logs

		//Calculamos cuantos logs inicialmente habia
		cantidadInicialLogs := reloj[0] + reloj[1] + reloj[2]

		//Actualizamos el reloj
		reloj = traducirReloj(req.RelojServidor.RelojServidor)

		//Aplicamos los cambios a el servidor actual
		aplicarCambios("LogRegistro.txt", req.Logs, cantidadInicialLogs)

		return nil, nil
	}

}
func (s *server) Finalizar(ctx context.Context, req *pb.MensajeSimple) (*pb.MensajeSimple, error) {
	fmt.Println("Llego senal termino, cerrando servidor...")
	defer os.Exit(1)
	return &pb.MensajeSimple{Valor: "1"}, nil
}

func main() {

	//Por mientras, dale el nroServidor a mano. 0: Tierra, 1: Titan, 2: Marte
	//Copialo en diferentes carpetas y pruebalo.
	//Tienes tambien que darle el puerto a los servers a mano. Ve tu cual a cual.
	//Ademas, cambia la var numberToIp que esta al principio pa ponerle
	//las direcciones que vas a usar por mientras.
	nroServidor = 0
	port := ":50051" // Puerto de DataNode

	listner, err := net.Listen("tcp", port)

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	pb.RegisterServidorPlanetarioServer(serv, &server{})

	//FALTA AGREGAR GO FUNCTION PARA HACER MERGE CADA 1 MIN

	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}

}
