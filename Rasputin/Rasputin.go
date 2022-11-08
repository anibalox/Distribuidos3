package main

import (
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	pb "Tarea3/Proto"
)

type server struct {
	pb.UnimplementedBrokerRasputinServer
}

var reloj [3]int = [3]int{0, 0, 0}
var nroServidor int

func (s *server) AgregarBase(context.Context, *pb.Peticion) (*pb.Direccion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AgregarBase not implemented")
}
func (s *server) RenombrarBase(context.Context, *pb.Peticion) (*pb.Direccion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenombrarBase not implemented")
}
func (s *server) ActualizarValor(context.Context, *pb.Peticion) (*pb.Direccion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActualizarValor not implemented")
}
func (s *server) BorrarBase(context.Context, *pb.Peticion) (*pb.Direccion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BorrarBase not implemented")
}
func (s *server) GetSoldados(context.Context, *pb.DatosBase) (*pb.SoldadosBase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSoldados not implemented")
}

func main() {

	nroServidor, _ = strconv.Atoi(os.Args[1]) //1: Tierra , 2: Titan, 3: Marte Le puedes dar un numero Felipe?
	port := ":50051"                          // Puerto de DataNode

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
