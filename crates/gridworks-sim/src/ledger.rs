use serde::{Deserialize, Serialize};
use std::collections::BTreeSet;
use std::fmt::{Display, Formatter};

pub const MONEY_SCALE: i64 = 1;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum LedgerError {
    InvalidIdentifier(String),
    InvalidCurrency(String),
    InvalidAmount,
    AmountOverflow,
    DuplicateAccount(String),
    UnknownAccount(String),
    DuplicateTransaction(String),
    DuplicateIdempotency(String),
    ImbalancedTransaction,
    MixedCurrency,
    ImmutableTransaction(String),
    UnknownReversal(String),
}
impl Display for LedgerError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "{self:?}")
    }
}
impl std::error::Error for LedgerError {}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum AccountClass {
    Asset,
    Liability,
    Equity,
    Revenue,
    Expense,
    ClearingSystem,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum LineSide {
    Debit,
    Credit,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct LedgerAccount {
    pub account_id: String,
    pub owner_ref: String,
    pub account_class: AccountClass,
    pub currency: String,
    pub active: bool,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct JournalLine {
    pub line_id: String,
    pub transaction_id: String,
    pub sequence: u32,
    pub account_id: String,
    pub currency: String,
    pub side: LineSide,
    pub amount_minor: i64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct JournalTransaction {
    pub transaction_id: String,
    pub transaction_type: String,
    pub effective_time_ms: u64,
    pub idempotency_key: String,
    pub source_ref: String,
    pub currency: String,
    pub reversal_of: Option<String>,
    pub lines: Vec<JournalLine>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum LedgerCommand {
    Post {
        transaction_id: String,
        transaction_type: String,
        idempotency_key: String,
        source_ref: String,
        currency: String,
        lines: Vec<JournalLineInput>,
    },
    Reverse {
        transaction_id: String,
        reversal_id: String,
        idempotency_key: String,
    },
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct JournalLineInput {
    pub line_id: String,
    pub sequence: u32,
    pub account_id: String,
    pub side: LineSide,
    pub amount_minor: i64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum LedgerEvent {
    Posted {
        transaction_id: String,
        debits_minor: i64,
        credits_minor: i64,
    },
    Duplicate {
        idempotency_key: String,
        transaction_id: String,
    },
    Reversed {
        reversal_id: String,
        original_transaction_id: String,
    },
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(deny_unknown_fields)]
pub struct LedgerState {
    pub accounts: Vec<LedgerAccount>,
    pub transactions: Vec<JournalTransaction>,
}

impl LedgerState {
    pub fn canonicalized(&self) -> Self {
        let mut state = self.clone();
        state
            .accounts
            .sort_by(|a, b| a.account_id.cmp(&b.account_id));
        state
            .transactions
            .sort_by(|a, b| a.transaction_id.cmp(&b.transaction_id));
        for tx in &mut state.transactions {
            tx.lines
                .sort_by(|a, b| a.sequence.cmp(&b.sequence).then(a.line_id.cmp(&b.line_id)));
        }
        state
    }
    pub fn validate(&self) -> Result<(), LedgerError> {
        let mut accounts = BTreeSet::new();
        for account in &self.accounts {
            valid_id(&account.account_id)?;
            valid_id(&account.owner_ref)?;
            valid_currency(&account.currency)?;
            if !accounts.insert(&account.account_id) {
                return Err(LedgerError::DuplicateAccount(account.account_id.clone()));
            }
        }
        let mut txs = BTreeSet::new();
        let mut keys = BTreeSet::new();
        for tx in &self.transactions {
            self.validate_transaction(tx)?;
            if !txs.insert(&tx.transaction_id) {
                return Err(LedgerError::DuplicateTransaction(tx.transaction_id.clone()));
            }
            if !keys.insert(&tx.idempotency_key) {
                return Err(LedgerError::DuplicateIdempotency(
                    tx.idempotency_key.clone(),
                ));
            }
        }
        Ok(())
    }
    fn validate_transaction(&self, tx: &JournalTransaction) -> Result<(), LedgerError> {
        valid_id(&tx.transaction_id)?;
        valid_id(&tx.transaction_type)?;
        valid_id(&tx.source_ref)?;
        valid_currency(&tx.currency)?;
        let account_ids = self
            .accounts
            .iter()
            .map(|a| a.account_id.as_str())
            .collect::<BTreeSet<_>>();
        let mut sequences = BTreeSet::new();
        let mut debits = 0i64;
        let mut credits = 0i64;
        for line in &tx.lines {
            valid_id(&line.line_id)?;
            if !sequences.insert(line.sequence) {
                return Err(LedgerError::ImmutableTransaction(tx.transaction_id.clone()));
            }
            if !account_ids.contains(line.account_id.as_str()) {
                return Err(LedgerError::UnknownAccount(line.account_id.clone()));
            }
            let account = self
                .accounts
                .iter()
                .find(|account| account.account_id == line.account_id)
                .ok_or_else(|| LedgerError::UnknownAccount(line.account_id.clone()))?;
            if account.currency != tx.currency || line.currency != tx.currency {
                return Err(LedgerError::MixedCurrency);
            }
            if line.amount_minor <= 0 {
                return Err(LedgerError::InvalidAmount);
            }
            match line.side {
                LineSide::Debit => {
                    debits = debits
                        .checked_add(line.amount_minor)
                        .ok_or(LedgerError::AmountOverflow)?
                }
                LineSide::Credit => {
                    credits = credits
                        .checked_add(line.amount_minor)
                        .ok_or(LedgerError::AmountOverflow)?
                }
            }
        }
        if debits == 0 || debits != credits {
            return Err(LedgerError::ImbalancedTransaction);
        }
        Ok(())
    }
    pub fn post(
        &mut self,
        effective_time_ms: u64,
        command: &LedgerCommand,
    ) -> Result<LedgerEvent, LedgerError> {
        self.validate()?;
        match command {
            LedgerCommand::Post {
                transaction_id,
                transaction_type,
                idempotency_key,
                source_ref,
                currency,
                lines,
            } => {
                if let Some(existing) = self
                    .transactions
                    .iter()
                    .find(|tx| tx.idempotency_key == *idempotency_key)
                {
                    return Ok(LedgerEvent::Duplicate {
                        idempotency_key: idempotency_key.clone(),
                        transaction_id: existing.transaction_id.clone(),
                    });
                }
                let tx = JournalTransaction {
                    transaction_id: transaction_id.clone(),
                    transaction_type: transaction_type.clone(),
                    effective_time_ms,
                    idempotency_key: idempotency_key.clone(),
                    source_ref: source_ref.clone(),
                    currency: currency.clone(),
                    reversal_of: None,
                    lines: lines
                        .iter()
                        .map(|line| JournalLine {
                            line_id: line.line_id.clone(),
                            transaction_id: transaction_id.clone(),
                            sequence: line.sequence,
                            account_id: line.account_id.clone(),
                            currency: currency.clone(),
                            side: line.side,
                            amount_minor: line.amount_minor,
                        })
                        .collect(),
                };
                let mut staged = self.clone();
                staged.transactions.push(tx.clone());
                staged.validate()?;
                *self = staged;
                let total = tx
                    .lines
                    .iter()
                    .filter(|l| l.side == LineSide::Debit)
                    .try_fold(0i64, |sum, l| {
                        sum.checked_add(l.amount_minor)
                            .ok_or(LedgerError::AmountOverflow)
                    })?;
                Ok(LedgerEvent::Posted {
                    transaction_id: tx.transaction_id,
                    debits_minor: total,
                    credits_minor: total,
                })
            }
            LedgerCommand::Reverse {
                transaction_id,
                reversal_id,
                idempotency_key,
            } => {
                if let Some(existing) = self
                    .transactions
                    .iter()
                    .find(|tx| tx.idempotency_key == *idempotency_key)
                {
                    return Ok(LedgerEvent::Duplicate {
                        idempotency_key: idempotency_key.clone(),
                        transaction_id: existing.transaction_id.clone(),
                    });
                }
                let original = self
                    .transactions
                    .iter()
                    .find(|tx| tx.transaction_id == *transaction_id)
                    .cloned()
                    .ok_or_else(|| LedgerError::UnknownReversal(transaction_id.clone()))?;
                if original.reversal_of.is_some() {
                    return Err(LedgerError::ImmutableTransaction(transaction_id.clone()));
                }
                let mut staged = self.clone();
                staged.transactions.push(JournalTransaction {
                    transaction_id: reversal_id.clone(),
                    transaction_type: "ledger.reversal".to_owned(),
                    effective_time_ms,
                    idempotency_key: idempotency_key.clone(),
                    source_ref: original.source_ref.clone(),
                    currency: original.currency.clone(),
                    reversal_of: Some(transaction_id.clone()),
                    lines: original
                        .lines
                        .iter()
                        .map(|line| JournalLine {
                            line_id: format!("{}.reversal", line.line_id),
                            transaction_id: reversal_id.clone(),
                            sequence: line.sequence,
                            account_id: line.account_id.clone(),
                            currency: line.currency.clone(),
                            side: match line.side {
                                LineSide::Debit => LineSide::Credit,
                                LineSide::Credit => LineSide::Debit,
                            },
                            amount_minor: line.amount_minor,
                        })
                        .collect(),
                });
                staged.validate()?;
                *self = staged;
                Ok(LedgerEvent::Reversed {
                    reversal_id: reversal_id.clone(),
                    original_transaction_id: transaction_id.clone(),
                })
            }
        }
    }
    pub fn balance_minor(&self, account_id: &str, currency: &str) -> Result<i64, LedgerError> {
        let mut balance = 0i64;
        for tx in &self.transactions {
            for line in &tx.lines {
                if line.account_id == account_id {
                    if line.currency != currency {
                        return Err(LedgerError::MixedCurrency);
                    }
                    let signed = match line.side {
                        LineSide::Debit => line.amount_minor,
                        LineSide::Credit => -line.amount_minor,
                    };
                    balance = balance
                        .checked_add(signed)
                        .ok_or(LedgerError::AmountOverflow)?;
                }
            }
        }
        Ok(balance)
    }
}
fn valid_id(value: &str) -> Result<(), LedgerError> {
    if value.is_empty() {
        Err(LedgerError::InvalidIdentifier(value.to_owned()))
    } else {
        Ok(())
    }
}
fn valid_currency(value: &str) -> Result<(), LedgerError> {
    if value.len() == 3 && value.chars().all(|c| c.is_ascii_uppercase()) {
        Ok(())
    } else {
        Err(LedgerError::InvalidCurrency(value.to_owned()))
    }
}

pub fn synthetic_ledger_fixture() -> LedgerState {
    LedgerState {
        accounts: vec![
            LedgerAccount {
                account_id: "account.system.clearing".to_owned(),
                owner_ref: "principal.gridworks.system".to_owned(),
                account_class: AccountClass::ClearingSystem,
                currency: "CRD".to_owned(),
                active: true,
            },
            LedgerAccount {
                account_id: "account.synthetic.wallet".to_owned(),
                owner_ref: "principal.synthetic.player".to_owned(),
                account_class: AccountClass::Asset,
                currency: "CRD".to_owned(),
                active: true,
            },
        ],
        transactions: Vec::new(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    fn post(id: &str, key: &str, amount: i64) -> LedgerCommand {
        LedgerCommand::Post {
            transaction_id: id.to_owned(),
            transaction_type: "ledger.synthetic_grant".to_owned(),
            idempotency_key: key.to_owned(),
            source_ref: "event.synthetic.grant".to_owned(),
            currency: "CRD".to_owned(),
            lines: vec![
                JournalLineInput {
                    line_id: format!("{id}.debit"),
                    sequence: 1,
                    account_id: "account.system.clearing".to_owned(),
                    side: LineSide::Debit,
                    amount_minor: amount,
                },
                JournalLineInput {
                    line_id: format!("{id}.credit"),
                    sequence: 2,
                    account_id: "account.synthetic.wallet".to_owned(),
                    side: LineSide::Credit,
                    amount_minor: amount,
                },
            ],
        }
    }
    #[test]
    fn balanced_posting_and_derived_balances_are_exact() {
        let mut ledger = synthetic_ledger_fixture();
        let event = ledger
            .post(0, &post("tx.grant", "key.grant", 1_000))
            .unwrap();
        assert!(matches!(
            event,
            LedgerEvent::Posted {
                debits_minor: 1_000,
                credits_minor: 1_000,
                ..
            }
        ));
        assert_eq!(
            ledger.balance_minor("account.system.clearing", "CRD"),
            Ok(1_000)
        );
        assert_eq!(
            ledger.balance_minor("account.synthetic.wallet", "CRD"),
            Ok(-1_000)
        );
    }
    #[test]
    fn imbalance_currency_and_overflow_reject_without_mutation() {
        let mut ledger = synthetic_ledger_fixture();
        let mut bad = post("tx.bad", "key.bad", 100);
        if let LedgerCommand::Post { ref mut lines, .. } = bad {
            lines[1].amount_minor = 99;
        }
        let before = ledger.clone();
        assert_eq!(
            ledger.post(0, &bad),
            Err(LedgerError::ImbalancedTransaction)
        );
        assert_eq!(ledger, before);
        let mut mixed = post("tx.mixed", "key.mixed", 100);
        if let LedgerCommand::Post {
            ref mut currency, ..
        } = mixed
        {
            *currency = "USD".to_owned();
        }
        assert_eq!(ledger.post(0, &mixed), Err(LedgerError::MixedCurrency));
        assert_eq!(ledger, before);
        let mut overflow = post("tx.overflow", "key.overflow", i64::MAX);
        if let LedgerCommand::Post { ref mut lines, .. } = overflow {
            lines.push(JournalLineInput {
                line_id: "tx.overflow.extra".to_owned(),
                sequence: 3,
                account_id: "account.system.clearing".to_owned(),
                side: LineSide::Debit,
                amount_minor: 1,
            });
        }
        assert!(matches!(
            ledger.post(0, &overflow),
            Err(LedgerError::ImbalancedTransaction | LedgerError::AmountOverflow)
        ));
        assert_eq!(ledger, before);
    }
    #[test]
    fn duplicate_posting_and_reversal_are_append_only() {
        let mut ledger = synthetic_ledger_fixture();
        ledger
            .post(0, &post("tx.grant", "key.grant", 1_000))
            .unwrap();
        let duplicate = ledger
            .post(0, &post("tx.other", "key.grant", 1_000))
            .unwrap();
        assert!(
            matches!(duplicate, LedgerEvent::Duplicate { transaction_id, .. } if transaction_id == "tx.grant")
        );
        let original = ledger.transactions[0].clone();
        let reversal = ledger
            .post(
                1,
                &LedgerCommand::Reverse {
                    transaction_id: "tx.grant".to_owned(),
                    reversal_id: "tx.reversal".to_owned(),
                    idempotency_key: "key.reversal".to_owned(),
                },
            )
            .unwrap();
        assert!(matches!(reversal, LedgerEvent::Reversed { .. }));
        assert_eq!(ledger.transactions[0], original);
        assert_eq!(
            ledger.balance_minor("account.system.clearing", "CRD"),
            Ok(0)
        );
        assert_eq!(
            ledger.balance_minor("account.synthetic.wallet", "CRD"),
            Ok(0)
        );
    }
}
