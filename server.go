package main

import (
	db_pkg "bank-grpc/db"
	"bank-grpc/pb"
	"fmt"
	"net"

	"bank-grpc/services/account_service"
	"bank-grpc/services/transaction_service"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Setup viper
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error in viper setup: %v", err)
	}

	// Connect to db
	db, err := db_pkg.ConnectAndMigrate()
	if err != nil {
		log.Fatalf("Db connection issue: %v", err)
	}

	
	// Create grpc servers
	address := fmt.Sprintf("localhost:%d", 8080)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	
	pb.RegisterAccountServiceServer(grpcServer, account_service.NewAccountServiceServer(db))
	pb.RegisterTransactionServiceServer(grpcServer, transaction_service.NewTransactionServiceServer(db))
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}