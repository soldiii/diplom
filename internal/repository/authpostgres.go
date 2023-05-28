package repository

import (
	"fmt"
	"strconv"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user *model.User) (int, error) {
	var id int
	query_usr := fmt.Sprintf("INSERT INTO %s (email, name, surname, patronymic, reg_date_time, encrypted_password, role, is_valid) VALUES ($1, $2, $3, $4, $5, $6, $7, true) RETURNING id", usersTable)
	row := r.db.QueryRow(query_usr, user.Email, user.Name, user.Surname, user.Patronymic, user.RegistrationDateTime, user.EncryptedPassword, user.Role)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateUserTempTable(user *model.UserCode) (int, error) {
	var id int
	sup_id, err := strconv.Atoi(user.SupervisorID)
	if err != nil {
		return 0, err
	}
	query_usr := fmt.Sprintf("INSERT INTO %s (email, name, surname, patronymic, reg_date_time, encrypted_password, role, supervisor_id, initials, code, attempt_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id", userCodesTable)
	row := r.db.QueryRow(query_usr, user.Email, user.Name, user.Surname, user.Patronymic, user.RegistrationDateTime, user.EncryptedPassword, user.Role, sup_id, user.SupervisorInitials, user.Code, user.AttemptNumber)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CheckForEmail(email string) error {
	query := fmt.Sprintf("SELECT email FROM %s WHERE email = $1 AND is_valid=true", usersTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&email); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) GetEmailOfMainSupervisor() (string, error) {
	var email string
	query := fmt.Sprintf("SELECT email FROM %s INNER JOIN %s ON users.id=supervisors.id ORDER BY users.reg_date_time ASC LIMIT 1", usersTable, supervisorsTable)
	row := r.db.QueryRow(query)
	if err := row.Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}

func (r *AuthPostgres) GetSupervisorEmailFromID(id int) (string, error) {
	var email string
	query := fmt.Sprintf("SELECT email FROM %s INNER JOIN %s ON users.id=supervisors.id WHERE supervisors.id=$1", usersTable, supervisorsTable)
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}

func (r *AuthPostgres) IsDBHaveMainSupervisor() (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT(SELECT count(*) FROM %s)=0", supervisorsTable)
	row := r.db.QueryRow(query)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}

	return !flag, nil
}

func (r *AuthPostgres) CreateMainSupervisor(user *model.User, supervisor *model.Supervisor) (int, error) {
	id, err := r.CreateUser(user)
	if err != nil {
		return 0, err
	}
	supervisor.ID = id
	query_sup := fmt.Sprintf("INSERT INTO %s (id, initials) VALUES ($1, $2)", supervisorsTable)
	r.db.QueryRow(query_sup, supervisor.ID, supervisor.SupervisorInitials)
	return id, nil
}

func (r *AuthPostgres) IsRegistrationCodeValid(email string, code string) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT CASE WHEN code = $2 THEN 1 ELSE 0 END as result FROM %s WHERE email = $1", userCodesTable)
	row := r.db.QueryRow(query, email, code)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if !flag {
		return false, nil
	}
	return true, nil
}

func (r *AuthPostgres) GetAttemptNumberByEmail(email string) (int, error) {
	var attemptNumber int
	getQuery := fmt.Sprintf("SELECT attempt_number FROM %s WHERE email = $1", userCodesTable)
	row := r.db.QueryRow(getQuery, email)
	if err := row.Scan(&attemptNumber); err != nil {
		return 0, err
	}
	return attemptNumber, nil
}

func (r *AuthPostgres) IncrementAttemptNumberByEmail(email string) {
	updateQuery := fmt.Sprintf("UPDATE %s SET attempt_number = attempt_number + 1 WHERE email = $1", userCodesTable)
	r.db.QueryRow(updateQuery, email)
}

func (r *AuthPostgres) DeleteFromTempTableByEmail(email string) {
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE email = $1", userCodesTable)
	r.db.QueryRow(deleteQuery, email)
}

func GetRoleByEmail(r *AuthPostgres, email string) (string, error) {
	var role string
	query := fmt.Sprintf("SELECT role FROM %s WHERE email = $1", userCodesTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&role); err != nil {
		return "", err
	}
	return role, nil
}

func UpdateAgent(r *AuthPostgres, email string) (int, error) {
	var id int
	queryUser := fmt.Sprintf("UPDATE %s SET name=uc.name, surname=uc.surname, patronymic=uc.patronymic, reg_date_time = uc.reg_date_time, encrypted_password=uc.encrypted_password, role=uc.role, is_valid=true FROM %s uc JOIN %s u ON u.email=uc.email WHERE uc.email = $1 RETURNING u.id", usersTable, userCodesTable, usersTable)
	row := r.db.QueryRow(queryUser, email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	queryAgent := fmt.Sprintf("UPDATE %s a SET a.supervisor_id=uc.supervisor_id FROM %s uc INNER JOIN %s u ON uc.email = u.email WHERE a.id = u.id AND uc.email = $1", agentsTable, userCodesTable, usersTable)
	r.db.QueryRow(queryAgent, email)
	return id, nil
}

func (r *AuthPostgres) MigrateFromTemporaryTable(email string) (int, error) {
	role, err := GetRoleByEmail(r, email)
	if err != nil {
		return 0, err
	}
	var userFlag bool = true
	emailFlag, err := r.IsEmailValid(email)
	if err != nil {
		return 0, err
	}
	if emailFlag {
		userFlag, err = r.IsUserValid(email)
	}
	if err != nil {
		return 0, err
	}
	var id int
	switch role {
	case "agent", "Agent":
		if !userFlag {
			id, err = UpdateAgent(r, email)
			if err != nil {
				return 0, err
			}
			break
		}
		queryUser := fmt.Sprintf("INSERT INTO %s (email, name, surname, patronymic, reg_date_time, encrypted_password, role, is_valid) SELECT email, name, surname, patronymic, reg_date_time, encrypted_password, role, true FROM %s WHERE email = $1 RETURNING id", usersTable, userCodesTable)
		row := r.db.QueryRow(queryUser, email)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		queryAgent := fmt.Sprintf("INSERT INTO %s (id, supervisor_id) SELECT users.id, usercodes.supervisor_id FROM %s JOIN users ON users.email = usercodes.email", agentsTable, userCodesTable)
		r.db.QueryRow(queryAgent)
	case "supervisor", "Supervisor":
		queryUser := fmt.Sprintf("INSERT INTO %s (email, name, surname, patronymic, reg_date_time, encrypted_password, role, is_valid) SELECT email, name, surname, patronymic, reg_date_time, encrypted_password, role, true FROM %s WHERE email = $1 RETURNING id", usersTable, userCodesTable)
		row := r.db.QueryRow(queryUser, email)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		queryAgent := fmt.Sprintf("INSERT INTO %s (id, initials) SELECT users.id, usercodes.initials FROM %s JOIN %s ON users.email = usercodes.email", supervisorsTable, userCodesTable, usersTable)
		r.db.QueryRow(queryAgent)
	}
	r.DeleteFromTempTableByEmail(email)
	return id, nil
}

func (r *AuthPostgres) GetRegistrationTimeByEmail(email string) (time.Time, error) {
	var regDateTime time.Time
	timeQuery := fmt.Sprintf("SELECT reg_date_time FROM %s WHERE email = $1", userCodesTable)
	row := r.db.QueryRow(timeQuery, email)
	if err := row.Scan(&regDateTime); err != nil {
		return time.Now(), err
	}
	return regDateTime, nil
}

func (r *AuthPostgres) GetUsersEmailsWithExpiredTime(timeNow time.Time, entryTime int64) ([]string, error) {
	var emails []string
	query := fmt.Sprintf("SELECT email FROM %s WHERE DATE_PART('epoch', $1 - reg_date_time)::int/60 > $2", userCodesTable)
	rows, err := r.db.Query(query, timeNow, entryTime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func (r *AuthPostgres) IsTempTableHaveUser(email string) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1) AS result", userCodesTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if !flag {
		return false, nil
	}
	return true, nil
}

func (r *AuthPostgres) GetCodeByEmail(email string) (string, error) {
	var code string
	query := fmt.Sprintf("SELECT code FROM %s WHERE email = $1", userCodesTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&code); err != nil {
		return "", err
	}
	return code, nil
}

func (r *AuthPostgres) GetUser(email, password string) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, role FROM  %s WHERE email = $1 AND encrypted_password = $2", usersTable)
	err := r.db.Get(&user, query, email, password)
	return &user, err
}

func (r *AuthPostgres) GetPassword(email string) (string, error) {
	var password string
	query := fmt.Sprintf("SELECT encrypted_password FROM %s WHERE email = $1", usersTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}

func (r *AuthPostgres) IsEmailValid(email string) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1) AS result", usersTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if !flag {
		return false, nil
	}
	return true, nil
}

func (r *AuthPostgres) PostRefreshToken(userID int, token string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (id, token) VALUES ($1, $2) RETURNING id", refreshTokensTable)
	row := r.db.QueryRow(query, userID, token)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) IsUserHaveRefreshToken(userID int) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1) AS result", refreshTokensTable)
	row := r.db.QueryRow(query, userID)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if !flag {
		return false, nil
	}
	return true, nil
}

func (r *AuthPostgres) UpdateRefreshToken(userID int, token string) (int, error) {
	var id int
	query := fmt.Sprintf("UPDATE %s SET token = $1 WHERE id = $2 RETURNING id", refreshTokensTable)
	row := r.db.QueryRow(query, token, userID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CompareRefreshTokens(token string, id int) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT CASE WHEN token = $1 THEN 1 ELSE 0 END as result FROM %s WHERE id = $2", refreshTokensTable)
	row := r.db.QueryRow(query, token, id)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if !flag {
		return false, nil
	}
	return true, nil
}

func (r *AuthPostgres) IsUserValid(email string) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1 AND is_valid = true) AS result", usersTable)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if !flag {
		return false, nil
	}
	return true, nil
}
