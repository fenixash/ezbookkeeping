package hdfc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestHdfcBankStatementTxtFileImporterParseImportedData_MinimumValidData(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, []byte(
		"\n"+
			"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			" 01/10/26  ,UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI           ,01/10/26 ,        210.00     ,          0.00     ,0000111111111111       ,     99790.11\n"+
			" 03/10/26  ,UPI-FAKE PERSON-FAKEPERSON@OKICICI-ICIC0000123-222222222222-UPI           ,03/10/26 ,          0.00     ,        500.00     ,0000222222222222       ,    100290.11\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 0, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2026-10-01 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(21000), allNewTransactions[0].Amount)
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)
	assert.Equal(t, "UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI", allNewTransactions[0].Comment)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, "2026-10-03 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(50000), allNewTransactions[1].Amount)
	assert.Equal(t, "", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)
	assert.Equal(t, "UPI-FAKE PERSON-FAKEPERSON@OKICICI-ICIC0000123-222222222222-UPI", allNewTransactions[1].Comment)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "", allNewAccounts[0].Name)
	assert.Equal(t, "INR", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)
}

func TestHdfcBankStatementTxtFileImporterParseImportedData_ParseHdfcBankStatementFileData(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	allNewTransactions, _, allNewSubExpenseCategories, allNewSubIncomeCategories, _, _, err := importer.ParseImportedData(context, user, []byte(
		"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			" 01/10/26  ,UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI           ,01/10/26 ,        210.00     ,          0.00     ,0000111111111111       ,     99790.11\n"+
			" 02/10/26  ,POS 512345XXXXXX6789 000000 02OCT26 12:31:50 FAKE STORE                   ,02/10/26 ,        299.00     ,          0.00     ,0000000000000000       ,     99491.11\n"+
			" 03/10/26  ,UPI-FAKE PERSON-FAKEPERSON@OKICICI-ICIC0000123-222222222222-UPI           ,03/10/26 ,          0.00     ,       5650.00     ,0000222222222222       ,    105141.11\n"+
			"10/10/26  ,ACH D- HDFC BANK LTD-333333333                                            ,10/10/26 ,      11916.00     ,          0.00     ,0000333333333333       ,     93225.11\n"+
			"18/10/26  ,UPI-FAKE PERSON TWO-FAKEPERSON2@OKAXIS-IDIB000K123-444444444444-UPI       ,18/10/26 ,          0.00     ,       1200.00     ,0000444444444444       ,     94425.11\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))

	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2026-10-01 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(21000), allNewTransactions[0].Amount)
	assert.Equal(t, "UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI", allNewTransactions[0].Comment)

	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2026-10-02 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(29900), allNewTransactions[1].Amount)
	assert.Equal(t, "POS 512345XXXXXX6789 000000 02OCT26 12:31:50 FAKE STORE", allNewTransactions[1].Comment)

	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[2].Type)
	assert.Equal(t, "2026-10-03 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(565000), allNewTransactions[2].Amount)
	assert.Equal(t, "UPI-FAKE PERSON-FAKEPERSON@OKICICI-ICIC0000123-222222222222-UPI", allNewTransactions[2].Comment)

	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[3].Type)
	assert.Equal(t, "2026-10-10 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(1191600), allNewTransactions[3].Amount)
	assert.Equal(t, "ACH D- HDFC BANK LTD-333333333", allNewTransactions[3].Comment)

	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[4].Type)
	assert.Equal(t, "2026-10-18 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime), time.UTC))
	assert.Equal(t, int64(120000), allNewTransactions[4].Amount)
	assert.Equal(t, "UPI-FAKE PERSON TWO-FAKEPERSON2@OKAXIS-IDIB000K123-444444444444-UPI", allNewTransactions[4].Comment)

	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)
}

func TestHdfcBankStatementTxtFileImporterParseImportedData_StopReadingDataAfterInvalidDataRow(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			" 01/10/26  ,UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI           ,01/10/26 ,        210.00     ,          0.00     ,0000111111111111       ,     99790.11\n"+
			" 03/10/26  ,UPI-FAKE PERSON-FAKEPERSON@OKICICI-ICIC0000123-222222222222-UPI           ,03/10/26 ,          0.00     ,        500.00     ,0000222222222222       ,    100290.11\n"+
			"\n"+
			"STATEMENT SUMMARY  :-\n"+
			"Opening Balance,,,,Debits,Credits,Closing Bal\n"+
			"100000.11,,,,210.00,500.00,100290.11\n"+
			"---  End Of Statement ---\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(allNewTransactions))
}

func TestHdfcBankStatementTxtFileImporterParseImportedData_ParseInvalidFileHeader(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"Date,Amount\n"+
			"01/10/26,210.00\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)
}

func TestHdfcBankStatementTxtFileImporterParseImportedData_NoTransactionData(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			"\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestHdfcBankStatementTxtFileImporterParseImportedData_ParseInvalidAmount(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	// invalid debit amount
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			" 01/10/26  ,UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI           ,01/10/26 ,       210.00 INR   ,          0.00     ,0000111111111111       ,     99790.11\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	// both debit amount and credit amount are zero
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			" 01/10/26  ,UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI           ,01/10/26 ,          0.00     ,          0.00     ,0000111111111111       ,     99790.11\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestHdfcBankStatementTxtFileImporterParseImportedData_ParseFewerFieldsInDataRow(t *testing.T) {
	importer := HdfcBankStatementTxtFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "INR",
	}

	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"  Date     ,Narration                                                                 ,Value Dat,Debit Amount       ,Credit Amount      ,Chq/Ref Number   ,Closing Balance\n"+
			" 01/10/26  ,UPI-FAKE MERCHANT-FAKEMERCHANT@YBL-YESB0YBLUPI-111111111111-UPI           ,01/10/26 ,        210.00     \n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.EqualError(t, err, errs.ErrFewerFieldsInDataRowThanInHeaderRow.Message)
}
