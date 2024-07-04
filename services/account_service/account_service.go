package account_service

import (
	pb "bank-grpc/pb"
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"bank-grpc/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServiceServer struct {
	db *gorm.DB
	pb.AccountServiceServer
}

func NewAccountServiceServer(db *gorm.DB) *AccountServiceServer {
	return &AccountServiceServer{
		db: db,
	}
}

func (a *AccountServiceServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Status, error) {
	if req == nil || req.AccountHolderName == "" || req.Password == "" {
		return &pb.Status{
			Success: false,
			Message: "Fields cannot be empty",
		}, nil
	}
	
	const accountNumDigits int = 8
	maxLimit := int64(int(math.Pow10(accountNumDigits)-1))

	hasGeneratedAccountNumber := false
	triesRemaining := 10
	var accountNum uint32 = 0

	for !hasGeneratedAccountNumber && triesRemaining > 0 {
		bigAccountNum, err := rand.Int(rand.Reader, big.NewInt(maxLimit))
		if err != nil {
			log.Errorf("Couldnt generate random number: %v", err)
			return &pb.Status{
				Success: false,
				Message: "Internal server error",
			}, status.Error(codes.Internal, "Internal server error. Please try again later")
		}

		accountNum = uint32(bigAccountNum.Int64())

		var existingAccounts []model.Account
		res := a.db.Where(&model.Account{AccountNumber: accountNum}).Find(&existingAccounts)
		if res.Error != nil {
			log.Error(res.Error)
			return &pb.Status{
				Success: false,
				Message: "Internal server error",
			}, status.Error(codes.Internal, "Internal server error. Please try again later")
		}

		if len(existingAccounts) == 0 {
			// Account number is unique
			account := model.Account{
				AccountNumber: accountNum,
				AccountHolderName: req.GetAccountHolderName(),
				AccountPassword: req.GetPassword(),
				Balance: 0,
			}

			a.db.Save(&account)
			hasGeneratedAccountNumber = true
		}
		triesRemaining--
	}
	
	if hasGeneratedAccountNumber {
		return &pb.Status{
			Success: true,
			Message: fmt.Sprintf("Successfully created account with account number: %v", accountNum),
		}, nil
	} else {
		return &pb.Status{
			Success: false,
			Message: "Internal server error",
		}, status.Error(codes.Internal, "Internal server error. Please try again later")
	}

}

func (a *AccountServiceServer) GetAccountBalance(ctx context.Context, req *pb.GetAccountBalanceRequest) (*pb.GetAccountBalanceResponse, error) {
	var account []*model.Account
	result := a.db.Where(&model.Account{AccountNumber: req.GetAccountNumber(), AccountPassword: req.GetPassword()}).Find(&account)
	if result.Error != nil {
		log.Error(result.Error)
		return &pb.GetAccountBalanceResponse{
			Status: &pb.Status{
				Success: false,
				Message: "Internal server error",
			},
			Balance: 0,
		}, status.Error(codes.Internal, "Internal server error. Please try again later")
	}

	if len(account) > 0 {
		return &pb.GetAccountBalanceResponse{
			Status: &pb.Status{
				Success: true,
			},
			Balance: account[0].Balance,
		}, nil
	} else {
		return &pb.GetAccountBalanceResponse{
			Status: &pb.Status{
				Success: false,
				Message: "Invalid account number or password",
			},
			Balance: 0,
		}, nil
	}
}

