package hdfc

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const hdfcBankStatementTransactionTimeColumnName = "Date"
const hdfcBankStatementTransactionNarrationColumnName = "Narration"
const hdfcBankStatementTransactionDebitAmountColumnName = "Debit Amount"
const hdfcBankStatementTransactionCreditAmountColumnName = "Credit Amount"

var hdfcBankStatementTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

var hdfcBankStatementTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:  utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE: utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
}

// hdfcBankStatementTransactionDataRowParser defines the structure of HDFC bank statement transaction data row parser
type hdfcBankStatementTransactionDataRowParser struct {
}

// Parse returns the converted transaction data row
func (p *hdfcBankStatementTransactionDataRowParser) Parse(ctx core.Context, user *models.User, dataRow datatable.CommonDataTableRow, rowId string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	transactionTimeValue := strings.TrimSpace(dataRow.GetData(hdfcBankStatementTransactionTimeColumnName))

	if !utils.IsValidDayMonthYearLongOrShortDateFormat(transactionTimeValue) {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_row_parser.Parse] cannot parse transaction time \"%s\" of transaction in \"%s\"", transactionTimeValue, rowId)
		return nil, false, errs.NewImportError(errs.ErrTransactionTimeInvalid, rowId, transactionTimeValue)
	}

	transactionTime, err := getHdfcBankStatementTransactionLongDateTime(transactionTimeValue)

	if err != nil {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_row_parser.Parse] cannot parse transaction time \"%s\" of transaction in \"%s\", because %s", transactionTimeValue, rowId, err.Error())
		return nil, false, errs.NewImportError(errs.ErrTransactionTimeInvalid, rowId, transactionTimeValue)
	}

	debitAmountValue := strings.TrimSpace(dataRow.GetData(hdfcBankStatementTransactionDebitAmountColumnName))
	creditAmountValue := strings.TrimSpace(dataRow.GetData(hdfcBankStatementTransactionCreditAmountColumnName))

	debitAmount, err := utils.ParseAmount(debitAmountValue)

	if err != nil {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_row_parser.Parse] cannot parse debit amount \"%s\" of transaction in \"%s\", because %s", debitAmountValue, rowId, err.Error())
		return nil, false, errs.NewImportError(errs.ErrAmountInvalid, rowId, debitAmountValue)
	}

	creditAmount, err := utils.ParseAmount(creditAmountValue)

	if err != nil {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_row_parser.Parse] cannot parse credit amount \"%s\" of transaction in \"%s\", because %s", creditAmountValue, rowId, err.Error())
		return nil, false, errs.NewImportError(errs.ErrAmountInvalid, rowId, creditAmountValue)
	}

	var transactionType models.TransactionType
	amount := int64(0)

	if debitAmount > 0 && creditAmount == 0 { // withdrawal
		transactionType = models.TRANSACTION_TYPE_EXPENSE
		amount = debitAmount
	} else if debitAmount == 0 && creditAmount > 0 { // deposit
		transactionType = models.TRANSACTION_TYPE_INCOME
		amount = creditAmount
	} else {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_row_parser.Parse] cannot parse transaction in \"%s\", because both debit amount \"%s\" and credit amount \"%s\" are invalid", rowId, debitAmountValue, creditAmountValue)
		return nil, false, errs.NewImportError(errs.ErrAmountInvalid, rowId, debitAmountValue+" / "+creditAmountValue)
	}

	data := make(map[datatable.TransactionDataTableColumn]string, len(hdfcBankStatementTransactionSupportedColumns))
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = transactionTime
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = hdfcBankStatementTransactionTypeNameMapping[transactionType]
	data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
	data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
	data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
	data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
	data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = strings.TrimSpace(dataRow.GetData(hdfcBankStatementTransactionNarrationColumnName))

	return data, true, nil
}

// getHdfcBankStatementTransactionLongDateTime returns the long date time in the date format of HDFC bank statement, that is day/month/year
func getHdfcBankStatementTransactionLongDateTime(date string) (string, error) {
	date = strings.ReplaceAll(date, ".", "-")
	date = strings.ReplaceAll(date, "/", "-")
	date = strings.ReplaceAll(date, "'", "-")
	items := strings.Split(date, "-")

	if len(items) != 3 || items[0] == "" || items[1] == "" || items[2] == "" {
		return "", errs.ErrTransactionTimeInvalid
	}

	day := items[0]
	month := items[1]
	year := items[2]

	return utils.FormatYearMonthDayToLongDateTime(year, month, day)
}

// createNewHdfcBankStatementTransactionDataRowParser returns HDFC bank statement transaction data row parser
func createNewHdfcBankStatementTransactionDataRowParser() datatable.CommonTransactionDataRowParser {
	return &hdfcBankStatementTransactionDataRowParser{}
}
