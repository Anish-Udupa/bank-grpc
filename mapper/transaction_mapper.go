package mapper

import (
	"bank-grpc/model"
	"bank-grpc/pb"
)

func MapTransactionToAccountStatements(transactions []model.Transaction) []*pb.AccountStatements {
	statements := []*pb.AccountStatements{}

	for _, t := range transactions {
		statements = append(statements, &pb.AccountStatements{
			FromAccountNumber: t.FromAccountNumber,
			ToAccountNumber: t.ToAccountNumber,
			TransactionDate: t.UpdatedAt.Format("2006-01-02"),
			Amount: t.Amount,
			Message: t.Message,
		})
	}

	return statements
}