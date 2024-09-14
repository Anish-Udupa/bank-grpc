package transaction_service

import (
	pb "bank-grpc/pb"
	"bank-grpc/services/helper"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const MAX_TRANSACTION_AMT uint32 = 50000

type TransactionServiceServer struct {
	db *gorm.DB
	pb.TransactionServiceServer
}

func NewTransactionServiceServer(db *gorm.DB) *TransactionServiceServer {
	return &TransactionServiceServer{
		db: db,
	}
}

func (t *TransactionServiceServer) MakeTransaction(ctx context.Context, req *pb.MakeTransactionRequest) (*pb.Status, error) {
	fromAcc := req.GetAccountNumber()
	fromAccPasswd := req.GetPassword()
	toAcc := req.GetTransaction().GetToAccountNumber()
	transactionAmt := req.GetTransaction().GetAmount()
	transactionMsg := req.GetTransaction().GetMessage()

	// Check if to account exists
	toAccExists, err := helper.CheckAccExists(toAcc, t.db)
	if err != nil {
		return generateInternalServerErrMsg()
	}
	if !toAccExists {
		return &pb.Status{
			Success: false,
			Message: "To Account Number is invalid",
		}, nil
	}

	// Authenticate from account
	authSuccess, err := helper.AuthenticateAccount(fromAcc, fromAccPasswd, t.db)
	if err != nil {
		return generateInternalServerErrMsg()
	}
	if !authSuccess {
		return &pb.Status{
			Success: false,
			Message: "Failed to authenticate",
		}, nil
	}

	// Check if amount entered is correct
	if !(transactionAmt > 0 && transactionAmt < MAX_TRANSACTION_AMT) {
		return &pb.Status{
			Success: false,
			Message: fmt.Sprintf("Transaction amount must be between 0 and %d", MAX_TRANSACTION_AMT),
		}, nil
	}

	// Check if balance is necessary
	hasSufficientBal, err := helper.CheckSufficientBalance(fromAcc, transactionAmt, t.db)
	if err != nil {
		return generateInternalServerErrMsg()
	}
	if !hasSufficientBal {
		return &pb.Status{
			Success: false,
			Message: "Insufficient balance",
		}, nil
	}

	// Initiate transfer
	err = helper.MakeTransaction(fromAcc, toAcc, transactionAmt, transactionMsg, t.db)
	if err != nil {
		return &pb.Status{
			Success: false,
			Message: "Transaction failed due to some internal server error",
		}, nil
	}
	return &pb.Status{
		Success: true,
		Message: "Transfer successful",
	}, nil
}

func (t *TransactionServiceServer) MakeDeposit(ctx context.Context, req *pb.MakeDepositMakeWithdrawRequest) (*pb.Status, error) {
	accountNum := req.GetAccountNumber()
	accountPass := req.GetPassword()
	amount := req.GetAmount()
	message := req.GetMessage()

	// Authenticate
	isAuthenticated, err := helper.AuthenticateAccount(accountNum, accountPass, t.db)
	if err != nil {
		return generateInternalServerErrMsg()
	}

	if !isAuthenticated {
		return &pb.Status{
			Success: false,
			Message: "Failed to authenticate",
		}, nil
	}

	// Make deposit
	err = helper.MakeDeposit(accountNum, amount, message, t.db)
	if err != nil {
		return &pb.Status{
			Success: false,
			Message: "Transaction failed due to some internal server error",
		}, nil
	}

	return &pb.Status{
		Success: true,
		Message: "Deposit successful",
	}, nil
}

func (t *TransactionServiceServer) MakeWithdraw(ctx context.Context, req *pb.MakeDepositMakeWithdrawRequest) (*pb.Status, error) {
	accountNum := req.GetAccountNumber()
	accountPass := req.GetPassword()
	amount := req.GetAmount()
	message := req.GetMessage()

	// Authenticate
	isAuthenticated, err := helper.AuthenticateAccount(accountNum, accountPass, t.db)
	if err != nil {
		return generateInternalServerErrMsg()
	}

	if !isAuthenticated {
		return &pb.Status{
			Success: false,
			Message: "Failed to authenticate",
		}, nil
	}

	// Check sufficient balance
	hasSufficientBal, err := helper.CheckSufficientBalance(accountNum, amount, t.db)
	if err != nil {
		return generateInternalServerErrMsg()
	}
	if !hasSufficientBal {
		return &pb.Status{
			Success: false,
			Message: "Insufficient balance",
		}, nil
	}

	// Make withdraw
	err = helper.MakeWithdraw(accountNum, amount, message, t.db)
	if err != nil {
		return &pb.Status{
			Success: false,
			Message: "Transaction failed due to some internal server error",
		}, nil
	}

	return &pb.Status{
		Success: true,
		Message: "Withdraw successful",
	}, nil
}

func generateInternalServerErrMsg() (*pb.Status, error) {
	return &pb.Status{
		Success: false,
		Message: "Internal Server Error",
	}, status.Error(codes.Internal, "Internal server error. Please try again later")
}
