package main

import (
	gobackcloud "GoBackCloud/proto"
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	args := os.Args
	bkorigin := args[1]
	bkdestino := args[2]

	if len(args) >= 4 {
		log.Fatalln("error leyendo los args!!")
	}

	clientgrpc_conection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		status.Errorf(
				codes.Code(code.Code_ABORTED),
				fmt.Sprintf("Error al conectarse al servidor!!!"),
			)
	}
	defer clientgrpc_conection.Close()
	conection := gobackcloud.NewBackupServiceClient(clientgrpc_conection)

	pth, err := os.Getwd()
	if err != nil {
		log.Println("No se pudo obtener el Working Directory")
	}

	backupfile, err := ioutil.ReadFile(path.Join(pth, "Client",bkorigin))
	if err != nil {
		log.Fatalln("Error al leer el archivo de backup: ", err)
	}

	req := &gobackcloud.BackupRequest{
		DatabaseBackup: &gobackcloud.Backup{
			Backupfile: fmt.Sprintf("%s", backupfile),
			FileName:   bkdestino+".sql",
		},
	}

	res, err := conection.CreateBackup(context.Background(), req)
	if err != nil {
		log.Println("Error llamando la funcion para crear el backup")
	}
	fmt.Printf("Backup Guardado con los siguientes Datos: %v\n", res)

}
