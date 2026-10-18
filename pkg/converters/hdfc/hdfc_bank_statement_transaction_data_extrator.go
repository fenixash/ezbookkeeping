package hdfc

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// createNewHdfcBankStatementTransactionBasicDataTable returns the transaction data table from the original data table of HDFC bank statement file
func createNewHdfcBankStatementTransactionBasicDataTable(ctx core.Context, originalDataTable datatable.BasicDataTable) (datatable.BasicDataTable, error) {
	iterator := originalDataTable.DataRowIterator()
	foundDataHeaderLine := false
	headerColumnCount := 0
	fromIndex := 0
	toIndex := originalDataTable.DataRowCount()
	currentIndex := -1

	for iterator.HasNext() {
		row := iterator.Next()
		currentIndex++

		if !foundDataHeaderLine {
			if !isHdfcBankStatementDataHeaderRow(row) {
				continue
			}

			foundDataHeaderLine = true
			headerColumnCount = row.ColumnCount()
			fromIndex = currentIndex
			continue
		}

		// the data rows are always between the data header row and the first row which is not a transaction data row
		if !isHdfcBankStatementTransactionDataRow(row) {
			if !isEmptyDataRow(row) {
				log.Warnf(ctx, "[hdfc_bank_statement_transaction_data_extrator.createNewHdfcBankStatementTransactionBasicDataTable] stop reading transaction data in \"%s\", because \"%s\" is not a valid transaction data row", iterator.CurrentRowId(), row.GetData(0))
			}

			toIndex = currentIndex
			break
		}

		if row.ColumnCount() < headerColumnCount {
			log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_extrator.createNewHdfcBankStatementTransactionBasicDataTable] cannot parse \"%s\", because may missing some columns (column count %d in data row is less than header column count %d)", iterator.CurrentRowId(), row.ColumnCount(), headerColumnCount)
			return nil, errs.NewImportError(errs.ErrFewerFieldsInDataRowThanInHeaderRow, iterator.CurrentRowId())
		}
	}

	if !foundDataHeaderLine {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_extrator.createNewHdfcBankStatementTransactionBasicDataTable] cannot parse import data, because cannot find data header line")
		return nil, errs.ErrInvalidFileHeader
	}

	if toIndex-fromIndex < 2 {
		log.Errorf(ctx, "[hdfc_bank_statement_transaction_data_extrator.createNewHdfcBankStatementTransactionBasicDataTable] cannot parse import data, because there is no transaction data row")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return datatable.CreateSubBasicTable(originalDataTable, fromIndex, toIndex, true), nil
}

func isHdfcBankStatementDataHeaderRow(row datatable.BasicDataTableRow) bool {
	if row == nil || row.ColumnCount() < 2 {
		return false
	}

	if strings.TrimSpace(row.GetData(0)) != hdfcBankStatementTransactionTimeColumnName ||
		strings.TrimSpace(row.GetData(1)) != hdfcBankStatementTransactionNarrationColumnName {
		return false
	}

	hasDebitAmountColumn := false
	hasCreditAmountColumn := false

	for i := 0; i < row.ColumnCount(); i++ {
		columnName := strings.TrimSpace(row.GetData(i))

		if columnName == hdfcBankStatementTransactionDebitAmountColumnName {
			hasDebitAmountColumn = true
		} else if columnName == hdfcBankStatementTransactionCreditAmountColumnName {
			hasCreditAmountColumn = true
		}
	}

	return hasDebitAmountColumn && hasCreditAmountColumn
}

func isHdfcBankStatementTransactionDataRow(row datatable.BasicDataTableRow) bool {
	if row == nil || row.ColumnCount() < 1 {
		return false
	}

	return utils.IsValidDayMonthYearLongOrShortDateFormat(strings.TrimSpace(row.GetData(0)))
}

func isEmptyDataRow(row datatable.BasicDataTableRow) bool {
	if row == nil || row.ColumnCount() < 1 {
		return true
	}

	for i := 0; i < row.ColumnCount(); i++ {
		if strings.TrimSpace(row.GetData(i)) != "" {
			return false
		}
	}

	return true
}
