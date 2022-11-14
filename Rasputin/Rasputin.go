package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"context"

	pb "Tarea3/Proto"
)

type server struct {
	pb.UnimplementedBrokerRasputinServer
}

var numberToIp = map[string]string{
	"1": "10003",    //Colocar IP de Tierra
	"2": "100002",   // IP de Titan
	"3": "23231231", // IP Marte
}

func randomNumber(max int, min int) int {
	return rand.Intn(max-min) + min
}

func (s *server) DerivarConsulta(ctx context.Context, req *pb.MensajeSimple) (*pb.Direccion, error) {
	return &pb.Direccion{DireccionServidor: numberToIp[strconv.Itoa(randomNumber(4, 1))]}, nil
}
func (s *server) GetSoldados(ctx context.Context, req *pb.DatosBase) (*pb.SoldadosBase, error) {

	nroServer := strconv.Itoa(randomNumber(4, 1))
	connS, err := grpc.Dial(numberToIp[nroServer], grpc.WithInsecure())

	fmt.Println("Se inicio conexion con " + nroServer + " Esperando respuesta...")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}
	serviceCliente := pb.NewServidorPlanetarioClient(connS)

	res, _ := serviceCliente.GetSoldados(context.Background(), req)

	defer connS.Close()

	fmt.Println("Llego respuestas. Cerrando conexion")

	return res, nil
}
func (s *server) Finalizar(ctx context.Context, req *pb.MensajeSimple) (*pb.MensajeSimple, error) {

	defer os.Exit(1)
	fmt.Println("Llego senal termino. Enviando senal de termino a todos los servidores...")

	for nroServidor, direccionServidorPlanetario := range numberToIp { //Se itera sobre numberToIp
		connS, err := grpc.Dial(direccionServidorPlanetario, grpc.WithInsecure())
		if err != nil {
			panic("cannot create tcp connection" + err.Error())
		}
		serviceCliente := pb.NewServidorPlanetarioClient(connS)

		serviceCliente.Finalizar(context.Background(), &pb.MensajeSimple{Valor: "1"})

		fmt.Println("Se finalizo servidor planetario: " + nroServidor)

		connS.Close()
	}
	fmt.Println("Se finalizaron todos los servidores. Cerrando Broker Rasputin")

	return &pb.MensajeSimple{Valor: "1"}, nil
}

func main() {

	port := ":50051" // Dale el puerto que te sirva

	listner, err := net.Listen("tcp", port)

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	pb.RegisterBrokerRasputinServer(serv, &server{})

	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}

}
