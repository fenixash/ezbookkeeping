package hdfc

import (
	"time"

	"golang.org/x/text/encoding/unicode"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// hdfcBankStatementTxtFileImporter defines the structure of HDFC bank statement text file importer for transaction data
type hdfcBankStatementTxtFileImporter struct {
}

// Initialize a HDFC bank statement transaction data text file importer singleton instance
var (
	HdfcBankStatementTxtFileImporter = &hdfcBankStatementTxtFileImporter{}
)

// ParseImportedData returns the imported data by parsing the HDFC bank statement text data
func (c *hdfcBankStatementTxtFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	fallbackDecoder := unicode.UTF8.NewDecoder()
	csvReader := csv.CreateNewCsvReaderFromData(data, unicode.BOMOverride(fallbackDecoder), ',', true)
	csvDataTable, err := csvReader.ReadToCsvFileBasicDataTable(ctx, false)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTable, err := createNewHdfcBankStatementTransactionBasicDataTable(ctx, csvDataTable)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewCommonDataTableFromBasicDataTable(dataTable)
	transactionRowParser := createNewHdfcBankStatementTransactionDataRowParser()
	transactionDataTable := datatable.CreateNewTransactionDataTableFromCommonDataTable(commonDataTable, hdfcBankStatementTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(hdfcBankStatementTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
