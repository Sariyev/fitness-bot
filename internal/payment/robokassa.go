package payment

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type RobokassaConfig struct {
	MerchantLogin string
	Password1     string
	Password2     string
	IsTest        bool
}

// PendingPaymentCreator is a minimal interface for creating pending payment records.
type PendingPaymentCreator interface {
	CreatePendingPayment(ctx context.Context, userID int64, amountKZT int, provider string) (int64, error)
}

type RobokassaProvider struct {
	cfg  RobokassaConfig
	repo PendingPaymentCreator
}

const robokassaBaseURL = "https://auth.robokassa.kz/Merchant/Index.aspx"

func NewRobokassaProvider(cfg RobokassaConfig, repo PendingPaymentCreator) *RobokassaProvider {
	return &RobokassaProvider{cfg: cfg, repo: repo}
}

func (r *RobokassaProvider) CreatePayment(ctx context.Context, userID int64, amountKZT int) (*PaymentResult, error) {
	// Create pending payment to get auto-increment ID (= InvId for Robokassa)
	paymentID, err := r.repo.CreatePendingPayment(ctx, userID, amountKZT, "robokassa")
	if err != nil {
		return nil, fmt.Errorf("create pending payment: %w", err)
	}

	outSum := fmt.Sprintf("%d", amountKZT)
	invID := fmt.Sprintf("%d", paymentID)

	// Signature: MD5(MerchantLogin:OutSum:InvId:Password#1)
	sig := md5hex(fmt.Sprintf("%s:%s:%s:%s", r.cfg.MerchantLogin, outSum, invID, r.cfg.Password1))

	params := url.Values{}
	params.Set("MerchantLogin", r.cfg.MerchantLogin)
	params.Set("OutSum", outSum)
	params.Set("InvId", invID)
	params.Set("Description", "Полный доступ к фитнес-платформе")
	params.Set("SignatureValue", sig)
	params.Set("Encoding", "utf-8")
	if r.cfg.IsTest {
		params.Set("IsTest", "1")
	}

	redirectURL := robokassaBaseURL + "?" + params.Encode()

	return &PaymentResult{
		ProviderTxID: invID,
		Status:       "pending",
		RedirectURL:  redirectURL,
	}, nil
}

// VerifyCallback verifies the Robokassa ResultURL callback signature (uses Password#2).
func (r *RobokassaProvider) VerifyCallback(params map[string]string) (int64, error) {
	outSum := params["OutSum"]
	invID := params["InvId"]
	receivedSig := strings.ToLower(params["SignatureValue"])

	// Signature: MD5(OutSum:InvId:Password#2)
	expectedSig := md5hex(fmt.Sprintf("%s:%s:%s", outSum, invID, r.cfg.Password2))

	if receivedSig != expectedSig {
		return 0, fmt.Errorf("invalid signature: got %s, expected %s", receivedSig, expectedSig)
	}

	id, err := strconv.ParseInt(invID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid InvId: %w", err)
	}

	return id, nil
}

func md5hex(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
