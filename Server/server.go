package main

import (
	gobackcloud "GoBackCloud/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type server struct {
	gobackcloud.UnimplementedBackupServiceServer
}

func (s *server) CreateBackup(ctx context.Context, req *gobackcloud.BackupRequest) (*gobackcloud.BackupResponse, error) {
	// Avisando los calls a la funcion
	log.Println("Funcion Llamada para almacenar backup!")

	//Alistando status
	stats := 1

	// Obteniendo los parametros para armar el backup y la respuesta
	bkfile := req.GetDatabaseBackup().GetBackupfile()
	bkname := req.GetDatabaseBackup().GetFileName()

	// Creando el archivo para escribir el backup
	path, err :=  os.Getwd()
	file, err := os.Create(path+"\\"+bkname)
	if err != nil {
		log.Println("Error al crear el archivo del backup")
		stats = 2
	}
	defer file.Close()
	nbyte, err := file.Write([]byte(bkfile))
	if err != nil {
		log.Println("Error al escribir el archivo de backup")
		stats = 3
	}

	response := &gobackcloud.BackupResponse{
		StoragePath:  path+"\\"+bkname,
		StatusBackup: gobackcloud.Status(stats),
		WritedBytes:  int32(nbyte),
	}

	return response, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tcp_listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error Al escuchar por el puerto 50051"),
			)
	}

	log.Println("Servidor Levantado y escuchando en 0.0.0.0:50051")

	grpc_server := grpc.NewServer()
	gobackcloud.RegisterBackupServiceServer(grpc_server, &server{})

	if err := grpc_server.Serve(tcp_listener); err != nil {
		status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error levantando el servidor"),
			)
	}

}