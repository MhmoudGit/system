package main

import (
	"errors"
	"fmt"
	"time"
)

var DB = []Transaction{}

type Loyalty struct {
	Active     bool
	Earn       bool
	Redeem     bool
	EarnRate   float64
	RedeemRate float64
	Expiration int64
}

type Transaction struct {
	ID        int64
	RefNo     string
	UserID    string
	Amount    float64
	Points    int64
	Earn      bool
	Redeem    bool
	Rate      float64
	ExpiredAt time.Time
	CreatedAt time.Time
}

func (l *Loyalty) CreateTransaction(refNo, userID string, amount float64, earn, redeem bool) error {
	if !l.Active {
		return errors.New("loyalty program is not active")
	}

	if !l.Earn && earn {
		return errors.New("earn transactions are not allowed")
	}

	if !l.Redeem && redeem {
		return errors.New("redeem transactions are not allowed")
	}

	t := Transaction{
		ID:        int64(len(DB) + 1),
		RefNo:     refNo,
		UserID:    userID,
		Amount:    amount,
		Earn:      earn,
		Redeem:    redeem,
		ExpiredAt: time.Now().Add(time.Duration(l.Expiration) * time.Hour),
	}

	if earn {
		t.Rate = l.EarnRate
		t.Points = int64(amount * t.Rate)
	}

	if redeem {
		t.Rate = l.RedeemRate
		t.Points = -int64(amount / t.Rate)
	}

	DB = append(DB, t)

	return nil
}

func main() {
	l := Loyalty{
		Active:     true,
		Earn:       true,
		Redeem:     true,
		EarnRate:   1,
		RedeemRate: 1,
		Expiration: 24,
	}

	l.CreateTransaction("123", "user1", 100, true, false)
	l.CreateTransaction("124", "user1", 100, true, false)
	l.CreateTransaction("125", "user1", 100, false, true)

	totalPoints := 0

	for _, t := range DB {
		color := "\033[32m"
		if t.Redeem {
			color = "\033[31m"
		}
		fmt.Printf(color+"Transaction: %+v\033[0m\n", t)
		totalPoints += int(t.Points)
		fmt.Println("\033[33mTotal Points: ", totalPoints, "\033[0m")
	}

	fmt.Println(totalPoints)
}
