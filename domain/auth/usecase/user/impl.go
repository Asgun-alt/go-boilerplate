package authuser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go-boilerplate/domain/auth/constant"
	"go-boilerplate/domain/auth/model"
	"go-boilerplate/infrastructure/security"
	"go-boilerplate/pkg/utils"
	"path"
	"strconv"
	"text/template"
	"time"

	pkgconstant "go-boilerplate/pkg/constant"

	gomail "gopkg.in/mail.v2"

	"go-boilerplate/pkg/session"
)

func (u *usecase) RegisterUser(
	sess *session.Session,
	req *model.UserRegistrationRequest,
) (response *model.UserRegistrationResponse, err error) {
	hashPassword, err := security.EncryptPassword(
		u.cfgAuthCredential.SecretKey,
		u.cfgAuthCredential.IV,
		req.ConfirmPassword)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return
	}

	username := fmt.Sprintf(
		constant.UsernameIncompletedProfileUser,
		time.Now().Unix())
	userReq := req.ToCreateUserRequest(username, hashPassword)
	err = u.userUsecase.CreateUser(sess, userReq)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	err = u.SendVerificationCode(sess, username, req.Email)
	if err != nil {
		return nil, err
	}

	response = &model.UserRegistrationResponse{
		Message: constant.SuccessRegisterUser,
	}
	return
}

func (u *usecase) SendVerificationCode(sess *session.Session, username, email string) error {
	if email == "" {
		err := errors.New("email cannot be empty")
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	otpData, err := u.authRepository.GetOTPFromRedis(sess, email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	var otpNumber string
	if otpData == nil {
		otpNumber, err = utils.GenerateRandomNumber(constant.MagicNumberOTP)
		if err != nil {
			sess.SetError(pkgconstant.ErrorGeneral, err)
			return err
		}

		otpData = &model.OTPRedisRequest{
			OTP:   otpNumber,
			Count: 1,
			Email: email,
		}
	} else {
		if otpData.Count < constant.OTPMaxRequest {
			otpData.Count++
		} else if otpData.Count == constant.OTPMaxRequest {
			err = fmt.Errorf(constant.ErrRequestOTP, email)
			sess.SetError(pkgconstant.ErrorToManyRequest, err)
			return err
		}
	}

	otpExpired := time.Duration(constant.OTPExpired * int(time.Minute))
	if err = u.authRepository.SetOTPToRedis(
		sess,
		*otpData,
		otpExpired,
	); err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	var filepath = path.Join("templates", "user_verification_email.html")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	emailPayload := model.OTPVerificationEmailPayload{
		Username:  username,
		UserEmail: email,
		OTPCode:   otpNumber,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, emailPayload); err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}
	mailPort, err := strconv.Atoi(u.SmtpConfig.MailPort)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(u.SmtpConfig.MailAdminRecipient, u.SmtpConfig.MailFrom))
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("Verify Your %s App Email", u.SmtpConfig.MailFrom))
	m.SetBody("text/html", tpl.String())
	d := gomail.NewDialer(u.SmtpConfig.MailSMTP, mailPort, u.SmtpConfig.MailUsername, u.SmtpConfig.MailPassword)
	err = d.DialAndSend(m)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	return nil
}

func (u *usecase) LoginUser(
	sess *session.Session,
	req *model.UserLoginRequest,
) (*model.UserLoginResponse, error) {
	user, err := u.userUsecase.GetUserByEmail(sess, req.Email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	if user == nil {
		err = constant.UserNotFound
		sess.SetError(pkgconstant.ErrorNotFound, err)
		return nil, err
	}

	hashPassword, err := security.EncryptPassword(
		u.cfgAuthCredential.SecretKey,
		u.cfgAuthCredential.IV,
		req.Password)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	exists, err := utils.ComparePassword(user.Password, hashPassword)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	if !exists {
		err = constant.WrongEmailPassword
		sess.SetError(pkgconstant.ErrorUnauthorized, err)
		return nil, err
	}

	var isVerified bool
	if user.IsEmailVerified {
		isVerified = true
	}

	var (
		accessToken string
		exp         int
	)

	accessToken, err = u.authRepository.GetAccessTokenFromRedis(sess, user.UserID.String())
	if err != nil {
		accessToken, exp, err = u.generateToken(sess, user.UserID.String(), user.Username)
		if err != nil {
			sess.SetError(pkgconstant.ErrorGeneral, err)
			return nil, err
		}
	} else {
		claims, err := u.jwt.VerifyToken(accessToken)
		if err != nil {
			accessToken, exp, err = u.generateToken(sess, user.UserID.String(), user.Username)
			if err != nil {
				sess.SetError(pkgconstant.ErrorGeneral, err)
				return nil, err
			}
		} else if claims != nil {
			exp = int(claims.ExpiredAt)
		}
	}

	return &model.UserLoginResponse{
		Token:      accessToken,
		Expired:    exp,
		IsVerified: isVerified,
	}, nil
}

func (u usecase) VerifyUser(
	sess *session.Session,
	req *model.UserVerifyRequest,
) (*model.UserVerifyResponse, error) {
	user, err := u.userUsecase.GetUserByEmail(sess, req.Email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	if user == nil {
		err = constant.UserNotFound
		sess.SetError(pkgconstant.ErrorNotFound, err)
		return nil, err
	}

	if user.IsEmailVerified {
		err = fmt.Errorf("user has been verified")
		sess.SetError(pkgconstant.ErrorDupCheck, err)
		return nil, err
	}

	var countVerify int
	verifyResult, err := u.authRepository.GetVerifyOTPFromRedis(sess, req.Email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	if verifyResult != nil {
		if verifyResult.Count < constant.OTPMaxRequest {
			countVerify = verifyResult.Count
		} else if verifyResult.Count == constant.OTPMaxRequest {
			err = fmt.Errorf(constant.ErrVerifyOTP, req.Email)
			sess.SetError(pkgconstant.ErrorToManyRequest, err)
			return nil, err
		}

		countVerify++
	}

	result, err := u.authRepository.GetOTPFromRedis(sess, req.Email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	if result == nil {
		err = errors.New(constant.ErrOTPNotFound)
		sess.SetError(pkgconstant.ErrorRequest, err)
		return nil, err
	}

	if result.OTP != req.OTP {
		reqSet := model.OTPRedisVerifyRequest{
			Email: req.Email,
			Count: countVerify,
		}

		otpExpired := time.Duration(constant.OTPExpired * int(time.Minute))
		if err = u.authRepository.SetMaxVerifyOTPToRedis(sess,
			reqSet,
			otpExpired); err != nil {
			sess.SetError(pkgconstant.ErrorGeneral, err)
			return nil, err
		}

		err = errors.New(constant.ErrInvalidOTP)
		sess.SetError(pkgconstant.ErrorRequest, err)
		return nil, err
	}

	go u.removeCache(sess, req.Email)

	err = u.userUsecase.UpdateVerifiedUser(sess, req.Email, true)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	return &model.UserVerifyResponse{
		Mesage: constant.UserVerified,
	}, nil
}

func (u usecase) removeCache(sess *session.Session, email string) {
	sess.Ctx = context.Background()
	_ = u.authRepository.DeleleteOTPFromRedis(sess, email)
	_ = u.authRepository.DeleteVerifyOTPFromRedis(sess, email)
}
