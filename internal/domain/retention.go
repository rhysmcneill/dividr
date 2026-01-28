package domain

import (
	"time"
)

// -- RETENTION POLICY (The "Law") --

// TransactionTTL is the maximum lifespan of a draft transaction.
// Legal Basis: GDPR Storage Limitation.
const TransactionTTL = 90 * 24 * time.Hour

// PurgeOnSuccess: Once HMRC accepts a return, the raw bank data is destroyed.
const PurgeOnSuccess = true

// RetentionPolicyVersion tracks the active legal rule set.
// Update this string manually if we change the constants above
const RetentionPolicyVersion = "v1.0 (Strict-Ephemeral)" // REMEMVER: Update when changing policy

// -- DATA LIFECYCLE --

// TransactionStatus tracks the sorting process of a row.
type TransactionStatus string

const (
	StatusUnprocessed TransactionStatus = "unprocessed" // Raw upload
	StatusBusiness    TransactionStatus = "business"    // Tax relevant
	StatusPersonal    TransactionStatus = "personal"    // Private (Ignore)
	StatusExcluded    TransactionStatus = "excluded"    // Duplicate/Transfer (Ignore)
)

func IsValidStatus(s string) bool {
	switch TransactionStatus(s) {
	case StatusUnprocessed, StatusBusiness, StatusPersonal, StatusExcluded:
		return true
	}
	return false
}

// -- HMRC CATEGORIES (Tax Buckets) --

// Trade (SA103)
const (
	CatTradeTurnover = "turnover"
	CatTradeExpenses = "expenses_total"
	CatTradeGoods    = "cost_of_goods"
	CatTradeTravel   = "travel_costs"
	CatTradePremises = "premises_costs"
	CatTradeAdmin    = "admin_costs"
)

// Property (SA105)
const (
	CatPropRentIncome = "rent_income"
	CatPropRepairs    = "repairs_maintenance"
	CatPropFinance    = "finance_costs"
	CatPropLegal      = "legal_prof_fees"
	CatPropServices   = "services_provided"
	CatPropOther      = "other_expenses"
)

func IsValidCategory(c string) bool {
	switch c {
	case CatTradeTurnover, CatTradeExpenses, CatTradeGoods, CatTradeTravel, CatTradePremises, CatTradeAdmin:
		return true
	case CatPropRentIncome, CatPropRepairs, CatPropFinance, CatPropLegal, CatPropServices, CatPropOther:
		return true
	}
	return false
}
