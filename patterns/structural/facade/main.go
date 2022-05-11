package main

import (
	"fmt"
	"log"
)

type walletFacade struct {
	account      *account
	wallet       *wallet
	securityCode *securityCode
	notification *notification
	ledger       *ledger
}

func newWalletFacade(accountID string, code int) *walletFacade {
	fmt.Println("Creating account started")
	walletFacade := &walletFacade{
		account:      newAccount(accountID),
		securityCode: newSecurityCode(code),
		wallet:       newWallet(),
		notification: &notification{},
		ledger:       &ledger{},
	}
	fmt.Println("Account created")
	return walletFacade
}

func (w *walletFacade) addMoneyToWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Adding money to wallet started")
	err := w.account.checkAccount(accountID)
	if err != nil {
		return err
	}
	err = w.securityCode.checkSecurityCode(securityCode)
	if err != nil {
		return err
	}

	w.wallet.creditBalance(amount)
	w.notification.sendWalletCreditNotification()
	w.ledger.makeEntry(accountID, "credit", amount)
	return nil
}

func (w *walletFacade) deductMoneyFromWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Deduction money from wallet started")
	err := w.account.checkAccount(accountID)
	if err != nil {
		return err
	}
	err = w.securityCode.checkSecurityCode(securityCode)
	if err != nil {
		return err
	}
	err = w.wallet.debitBalance(amount)
	if err != nil {
		return err
	}

	w.notification.sendWalletDebitNotification()
	w.ledger.makeEntry(accountID, "credit", amount)
	return nil
}

type account struct {
	name string
}

func newAccount(accountName string) *account {
	return &account{
		name: accountName,
	}
}

func (a *account) checkAccount(accountName string) error {
	if a.name != accountName {
		return fmt.Errorf("account name is incorrect")
	}
	fmt.Println("Account Verified")
	return nil
}

type securityCode struct {
	code int
}

func newSecurityCode(code int) *securityCode {
	return &securityCode{
		code: code,
	}
}

func (s *securityCode) checkSecurityCode(inCode int) error {
	if s.code != inCode {
		return fmt.Errorf("security code is incorrect")
	}
	fmt.Println("Security code verified")
	return nil
}

type wallet struct {
	balance int
}

func newWallet() *wallet {
	return &wallet{
		balance: 0,
	}
}

func (w *wallet) creditBalance(amount int) {
	w.balance += amount
	fmt.Println("Wallet balance fullfilled successfully")
}

func (w *wallet) debitBalance(amount int) error {
	if w.balance < amount {
		return fmt.Errorf("balance is not sufficient")
	}
	fmt.Println("Wallet balance is sufficient")
	w.balance = w.balance - amount
	return nil
}

type ledger struct {
}

func (l *ledger) makeEntry(accountID, txnType string, amount int) {
	fmt.Printf("Making ledger entry dor account ID %s with txnType %s for amount %d\n", accountID, txnType, amount)
}

type notification struct {
}

func (n *notification) sendWalletCreditNotification() {
	fmt.Println("Sending notification on crediting wallet")
}

func (n *notification) sendWalletDebitNotification() {
	fmt.Println("Sending notification on debiting wallet")
}

func main() {
	fmt.Println()
	walletFacade := newWalletFacade("Ann", 1234)
	fmt.Println()

	err := walletFacade.addMoneyToWallet("Ann", 1234, 10)
	if err != nil {
		log.Fatalf("Error:%s\n", err.Error())
	}

	fmt.Println()
	err = walletFacade.deductMoneyFromWallet("Ann", 1234, 5)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}
