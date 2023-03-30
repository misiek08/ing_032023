package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

type Accounts []Account

func (accs Accounts) Len() int {
	return len(accs)
}

func (accs Accounts) Less(i, j int) bool {
	return strings.Compare(accs[i].Number, accs[j].Number) < 0
}

func (accs Accounts) Swap(i, j int) {
	accs[i], accs[j] = accs[j], accs[i]
}

type Account struct {
	Balance     decimal.Decimal `json:"balance"`
	DebitCount  int64           `json:"debitCount"`
	CreditCount int64           `json:"creditCount"`
	Number      string          `json:"account"`
}

type Transaction struct {
	DebitAccount  string
	CreditAccount string
	Amount        decimal.Decimal
}

func parseTransactions(data []byte) ([]Transaction, error) {
	var transactions []Transaction
	err := json.Unmarshal(data, &transactions)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse transactions")
	}
	return transactions, nil
}

func calculateTransactions(transactions []Transaction) Accounts {
	accountsMap := make(map[string]*Account)
	for _, transaction := range transactions {
		if _, ok := accountsMap[transaction.DebitAccount]; !ok {
			accountsMap[transaction.DebitAccount] = &Account{Number: transaction.DebitAccount}
		}
		if _, ok := accountsMap[transaction.CreditAccount]; !ok {
			accountsMap[transaction.CreditAccount] = &Account{Number: transaction.CreditAccount}
		}

		creditAccount := accountsMap[transaction.CreditAccount]
		creditAccount.CreditCount = accountsMap[transaction.CreditAccount].CreditCount + 1
		creditAccount.Balance = creditAccount.Balance.Add(transaction.Amount)
		debitAccount := accountsMap[transaction.DebitAccount]
		debitAccount.DebitCount += 1
		debitAccount.Balance = debitAccount.Balance.Sub(transaction.Amount)
	}

	accountList := make(Accounts, 0, len(accountsMap))

	for _, account := range accountsMap {
		accountList = append(accountList, *account)
	}

	sort.Sort(accountList)

	return accountList
}

func main() {
	decimal.MarshalJSONWithoutQuotes = true

	fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		transactions, err := parseTransactions(ctx.PostBody())
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBody([]byte(fmt.Sprintf("unable to parse request: %v", err)))
		}
		accountsStats := calculateTransactions(transactions)
		accountsJSON, err := json.Marshal(accountsStats)
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBody([]byte(fmt.Sprintf("unable to serialize accounts: %v", err)))
		} else {
			ctx.Response.SetBody(accountsJSON)
		}
	})
}
