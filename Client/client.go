package main

import (
	gobackcloud "GoBackCloud/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"path"
)

func main() {
	bkorigin := flag.String("bkorigen", "backup.sql", "Especifique el nombre del archivo local al que le hara backup")
	bkdestino := flag.String("bkdestino", "backup.sql", "Especifique el nombre de como se guardara el backup en el servidor de storage")
	server := flag.String("servidor", "localhost", "Especifique la IP del servidor al que se conectara por el puerto 50051")
	flag.Parse()

	clientgrpc_conection, err := grpc.Dial(*server+":50051", grpc.WithInsecure())
	if err != nil {
		status.Errorf(
				codes.Code(code.Code_ABORTED),
				fmt.Sprintf("Error al conectarse al servidor!!!"),
			)
	}
	defer clientgrpc_conection.Close()
	conection := gobackcloud.NewBackupServiceClient(clientgrpc_conection)

	if err != nil {
		log.Println("No se pudo obtener el Working Directory")
	}

	backupfile, err := ioutil.ReadFile(path.Join(*bkorigin))
	if err != nil {
		log.Fatalln("Error al leer el archivo de backup: ", err)
	}

	req := &gobackcloud.BackupRequest{
		DatabaseBackup: &gobackcloud.Backup{
			Backupfile: fmt.Sprintf("%s", backupfile),
			FileName:   *bkdestino+".sql",
		},
	}

	res, err := conection.CreateBackup(context.Background(), req)
	if err != nil {
		log.Fatalln("Error llamando la funcion para crear el backup", err)
	}
	fmt.Printf("Backup Guardado con los siguientes Datos: %v\n", res)

}
