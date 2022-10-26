package repository

import (
	"Avito/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
)

//goland:noinspection ALL
const requestByIdQuery = `  SELECT id
	FROM users
	WHERE user_id = $1 LIMIT 1
`

//goland:noinspection ALL
const getIdBalanceQuery = `  SELECT id, curr_amount
	FROM users
	WHERE user_id = $1 LIMIT 1
`

//goland:noinspection ALL

const logTransactionQuery = ` INSERT INTO 
	transactions(user_id_from, user_id_to, transaction_sum, status, event_type, created_at, updated_at)
	VALUES ((SELECT id FROM users WHERE user_id = $1), 
	        (SELECT id FROM users WHERE user_id = $2), 
	        $3, $4, $5, current_timestamp, current_timestamp)
`

//goland:noinspection ALL
const getTransactionQuery = `SELECT user_id_from, user_id_to, transaction_sum, status, event_type, created_at
	FROM transactions
	WHERE user_id_from = (SELECT id FROM users WHERE user_id = $1)
	AND created_at = (select created_at from transactions order by created_at desc limit 1)
`

//goland:noinspection ALL
const addFundsQuery = ` UPDATE users
	SET curr_amount = curr_amount + $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
	RETURNING user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const withdrawFundsQuery = ` UPDATE users
	SET curr_amount = curr_amount - $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
	RETURNING user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const decreasePendingAmountQuery = ` UPDATE users
	SET pending_amount = pending_amount - $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
	RETURNING user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const createAccQuery = ` INSERT INTO 
	users(user_id, curr_amount, pending_amount, last_updated)
	VALUES ($1, 0, 0, current_timestamp)
	RETURNING id
`

//goland:noinspection ALL
const logServiceOrderQuery = ` INSERT INTO
	service(user_id, invoice, service_id, order_id, status, created_at, updated_at)
	VALUES ((SELECT id FROM users WHERE user_id = $1), $2, $3, $4, $5, current_timestamp, current_timestamp)
	RETURNING user_id, service_id, invoice, status, created_at
`

//goland:noinspection ALL
const changeServiceStatusQuery = ` UPDATE service
	SET status = $5, 
	updated_at = current_timestamp
	WHERE user_id = (SELECT id FROM users WHERE user_id = $1) 
	AND order_id = $2 
	AND service_id = $3
	AND invoice = $4
`

//goland:noinspection ALL
const reserveAmountQuery = ` UPDATE users
	SET pending_amount = pending_amount + $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
`

//goland:noinspection ALL
const getLastServiceQuery = ` SELECT user_id, service_id, order_id, invoice, status, created_at, updated_at
	FROM service
	WHERE user_id = (SELECT id FROM users WHERE user_id = $1)
	AND order_id = $2 
	AND service_id = $3
	AND invoice = $4
`

//goland:noinspection ALL
const getLastServiceStatusQuery = ` SELECT status
	FROM service
	WHERE user_id = (SELECT id FROM users WHERE user_id = $1)
	AND order_id = $2 
	AND service_id = $3
	AND invoice = $4
`

type RequestPostgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *RequestPostgres {
	return &RequestPostgres{db: db}
}

func (r *RequestPostgres) GetBalance(userid models.GetBalanceRequest, ctx *gin.Context) (models.GetBalanceResponse, error) {
	var balanceRes models.GetBalanceResponse

	fail := func(err error) (models.GetBalanceResponse, error) {
		return balanceRes, fmt.Errorf("GetBalance: %v", err)
	}

	if userid.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	query := fmt.Sprintf(
		"SELECT ac.curr_amount, ac.pending_amount FROM users ac " +
			"WHERE user_id = $1")
	row := r.db.QueryRowContext(ctx, query, userid.UserId)

	if err := row.Scan(
		&balanceRes.Balance,
		&balanceRes.Pending,
	); err != nil {
		return models.GetBalanceResponse{}, err
	}

	return models.GetBalanceResponse{Balance: balanceRes.Balance, Pending: balanceRes.Pending}, nil
}

func (r *RequestPostgres) Deposit(depositReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error) {
	var depositResponse models.UpdateBalanceDepositResponse

	fail := func(err error) (models.UpdateBalanceDepositResponse, error) {
		return depositResponse, fmt.Errorf("DepositMoney: %v", err)
	}

	if depositReq.Sum <= 0 {
		err := errors.New("can't add negative or zero funds")
		return fail(err)
	}

	if depositReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	var exists int64

	if err := r.db.QueryRow(requestByIdQuery, depositReq.UserId).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			rows := r.db.QueryRow(createAccQuery, depositReq.UserId)
			if err := rows.Scan(
				&depositResponse.UserId,
			); err != nil {
				return depositResponse, err
			}
			log.Print("created new user ", depositReq.UserId, " in database")
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Println("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	_, err = tx.ExecContext(ctx, addFundsQuery, depositReq.UserId, depositReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, depositReq.UserId, depositReq.UserId, depositReq.Sum, "Completed", "Deposit")

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	var holder int

	rows := r.db.QueryRow(getTransactionQuery, depositReq.UserId)
	if err := rows.Scan(
		&depositResponse.UserId,
		&holder,
		&depositResponse.Sum,
		&depositResponse.Status,
		&depositResponse.EventType,
		&depositResponse.CreatedAt,
	); err != nil {
		return depositResponse, err
	}

	return models.UpdateBalanceDepositResponse{
		UserId:    depositResponse.UserId,
		Sum:       depositResponse.Sum,
		Status:    depositResponse.Status,
		EventType: depositResponse.EventType,
		CreatedAt: depositResponse.CreatedAt,
	}, nil
}

func (r *RequestPostgres) Withdrawal(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error) {
	var withdrawResponse models.UpdateBalanceWithdrawResponse

	fail := func(err error) (models.UpdateBalanceWithdrawResponse, error) {
		return withdrawResponse, fmt.Errorf("WithdrawMoney: %v", err)
	}

	if withdrawReq.Sum < 0 {
		err := errors.New("can't withdraw negative funds")
		return fail(err)
	}

	if withdrawReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Println("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalanceHolder struct {
		Id      int64
		Balance int64
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, withdrawReq.UserId).Scan(&idBalanceHolder.Id, &idBalanceHolder.Balance); err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user with that user id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}
	if idBalanceHolder.Balance < withdrawReq.Sum {
		err = errors.New("not enough funds")
		log.Println("not enough funds on the user")
		return fail(err)
	}
	_, err = tx.ExecContext(ctx, withdrawFundsQuery, withdrawReq.UserId, withdrawReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, withdrawReq.UserId, withdrawReq.UserId, withdrawReq.Sum, "Completed", "Withdrawal")

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	var holder int
	rows := r.db.QueryRow(getTransactionQuery, withdrawReq.UserId)
	if err := rows.Scan(
		&withdrawResponse.UserId,
		&holder,
		&withdrawResponse.Sum,
		&withdrawResponse.Status,
		&withdrawResponse.EventType,
		&withdrawResponse.CreatedAt,
	); err != nil {
		return withdrawResponse, err
	}
	log.Print("found acc ", withdrawReq.UserId, " in database, withdrew ", withdrawReq.Sum, " funds")
	return models.UpdateBalanceWithdrawResponse{
		UserId:    withdrawResponse.UserId,
		Sum:       withdrawResponse.Sum,
		Status:    withdrawResponse.Status,
		EventType: withdrawResponse.EventType,
		CreatedAt: withdrawResponse.CreatedAt,
	}, nil
}

func (r *RequestPostgres) Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error) {
	var transferRes models.TransferResponse

	fail := func(err error) (models.TransferResponse, error) {
		return transferRes, fmt.Errorf("TransferMoney: %v", err)
	}

	if transferReq.Sum < 0 {
		err := errors.New("can't transfer negative amount")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Println("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalanceHolder struct {
		Id      int64
		Balance int64
	}

	if err = tx.QueryRowContext(ctx, requestByIdQuery, transferReq.ReceiverId).Scan(&idBalanceHolder.Id); err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user with that receiver id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, transferReq.SenderId).Scan(&idBalanceHolder.Id, &idBalanceHolder.Balance); err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user with that sender id: add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	}
	if idBalanceHolder.Balance < transferReq.Sum {
		err = errors.New("not enough funds to transfer")
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, withdrawFundsQuery, transferReq.SenderId, transferReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, addFundsQuery, transferReq.ReceiverId, transferReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, transferReq.SenderId, transferReq.ReceiverId, transferReq.Sum, "Completed", "Withdrawn-Transfer")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, transferReq.ReceiverId, transferReq.SenderId, transferReq.Sum, "Completed", "Received-Transfer")

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	rows := r.db.QueryRow(getTransactionQuery, transferReq.SenderId)
	if err := rows.Scan(
		&transferRes.UserIdFrom,
		&transferRes.UserIdTo,
		&transferRes.Amount,
		&transferRes.Status,
		&transferRes.EventType,
		&transferRes.Timecode,
	); err != nil {
		return transferRes, err
	}

	return models.TransferResponse{
		UserIdTo:   transferRes.UserIdTo,
		UserIdFrom: transferRes.UserIdFrom,
		Amount:     transferRes.Amount,
		Status:     transferRes.Status,
		EventType:  transferRes.EventType,
		Timecode:   transferRes.Timecode,
	}, nil
}

func (r *RequestPostgres) ReserveService(reserveServiceReq models.ReserveServiceRequest, ctx *gin.Context) (models.ReserveServiceResponse, error) {
	var reserveRes models.ReserveServiceResponse

	fail := func(err error) (models.ReserveServiceResponse, error) {
		return reserveRes, fmt.Errorf("ReserveService: %v", err)
	}

	if reserveServiceReq.Payment < 0 {
		err := errors.New("can't reserve negative sum")
		return fail(err)
	}

	if reserveServiceReq.ServiceId < 0 {
		err := errors.New("illegal service ID")
		return fail(err)
	}

	if reserveServiceReq.OrderId < 0 {
		err := errors.New("illegal order ID")
		return fail(err)
	}

	if reserveServiceReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Println("")
	}

	// Defer a rollback in case anything fails.
	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var exists int64

	if err = tx.QueryRowContext(ctx, requestByIdQuery, reserveServiceReq.UserId).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			log.Println("no account with that user id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logServiceOrderQuery, reserveServiceReq.UserId, reserveServiceReq.Payment, reserveServiceReq.ServiceId,
		reserveServiceReq.OrderId, "Pending")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, reserveAmountQuery, reserveServiceReq.UserId, reserveServiceReq.Payment)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastServiceQuery, reserveServiceReq.UserId,
		reserveServiceReq.OrderId, reserveServiceReq.ServiceId, reserveServiceReq.Payment)
	if err := logServiceOrderRes.Scan(
		&reserveRes.UserId,
		&reserveRes.ServiceId,
		&reserveRes.OrderId,
		&reserveRes.Invoice,
		&reserveRes.Status,
		&reserveRes.CreatedAt,
		&reserveRes.UpdatedAt,
	); err != nil {
		return reserveRes, err
	}
	return reserveRes, nil
}

func (r *RequestPostgres) ApproveService(approveServiceReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error) {
	var approvalServiceResponse models.StatusServiceResponse

	fail := func(err error) (models.StatusServiceResponse, error) {
		return approvalServiceResponse, fmt.Errorf("ApproveService: %v", err)
	}

	if approveServiceReq.Payment < 0 {
		err := errors.New("can't withdraw negative funds")
		return fail(err)
	}

	if approveServiceReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Print("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalanceHolder struct {
		Id      int64
		Balance int64
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, approveServiceReq.UserId).Scan(
		&idBalanceHolder.Id,
		&idBalanceHolder.Balance,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Print("No user with that user id. Add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	}

	if idBalanceHolder.Balance < approveServiceReq.Payment {
		err = errors.New("not enough funds")
		return fail(err)
	}

	var status string

	if err = tx.QueryRowContext(ctx, getLastServiceStatusQuery, approveServiceReq.UserId, approveServiceReq.OrderId,
		approveServiceReq.ServiceId, approveServiceReq.Payment).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user with that user id. Add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	} else {
		if status == "Approved" {
			err = errors.New("this payment has already been approved")
			return fail(err)
		}
	}

	_, err = tx.ExecContext(ctx, changeServiceStatusQuery, approveServiceReq.UserId, approveServiceReq.OrderId,
		approveServiceReq.ServiceId, approveServiceReq.Payment, "Approved")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, withdrawFundsQuery, approveServiceReq.UserId, approveServiceReq.Payment)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, decreasePendingAmountQuery, approveServiceReq.UserId, approveServiceReq.Payment)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastServiceQuery, approveServiceReq.UserId, approveServiceReq.OrderId,
		approveServiceReq.ServiceId, approveServiceReq.Payment)
	if err := logServiceOrderRes.Scan(
		&approvalServiceResponse.UserId,
		&approvalServiceResponse.ServiceId,
		&approvalServiceResponse.OrderId,
		&approvalServiceResponse.Invoice,
		&approvalServiceResponse.Status,
		&approvalServiceResponse.CreatedAt,
		&approvalServiceResponse.UpdatedAt,
	); err != nil {
		return approvalServiceResponse, err
	}
	return approvalServiceResponse, nil
}

func (r *RequestPostgres) FailedService(failedServiceReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error) {
	var failedService models.StatusServiceResponse

	fail := func(err error) (models.StatusServiceResponse, error) {
		return failedService, fmt.Errorf("FailedService: %v", err)
	}

	if failedServiceReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Println("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalance struct {
		Id  int64
		Bal int64
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, failedServiceReq.UserId).Scan(
		&idBalance.Id,
		&idBalance.Bal,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user with that user-id: create a new one by depositing money")
			err = errors.New("no user with that user-id")
			return fail(err)
		}
		return fail(err)
	}

	var status string

	if err = tx.QueryRowContext(ctx, getLastServiceStatusQuery, failedServiceReq.UserId, failedServiceReq.OrderId,
		failedServiceReq.ServiceId, failedServiceReq.Payment).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			log.Println("no service log with that parameters")
			err = errors.New("no service log with that parameters")
			return fail(err)
		}
		return fail(err)
	} else {
		if status == "Approved" || status == "Cancelled" {
			err = fmt.Errorf("this payment has already been %s", status)
			return fail(err)
		}
	}

	_, err = tx.ExecContext(ctx, changeServiceStatusQuery, failedServiceReq.UserId, failedServiceReq.OrderId,
		failedServiceReq.ServiceId, failedServiceReq.Payment, "Cancelled")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, decreasePendingAmountQuery, failedServiceReq.UserId, failedServiceReq.Payment)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastServiceQuery, failedServiceReq.UserId, failedServiceReq.OrderId,
		failedServiceReq.ServiceId, failedServiceReq.Payment)
	if err := logServiceOrderRes.Scan(
		&failedService.UserId,
		&failedService.ServiceId,
		&failedService.OrderId,
		&failedService.Invoice,
		&failedService.Status,
		&failedService.CreatedAt,
		&failedService.UpdatedAt,
	); err != nil {
		return failedService, err
	}
	return failedService, nil
}
