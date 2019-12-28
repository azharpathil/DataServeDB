// Copyright (c) 2019 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package dbtable

import (
	"errors"
	"fmt"

	"DataServeDB/dbsystem"
	db_rules "DataServeDB/dbsystem/rules"
)

//TODO: move it to error messages (single location)

// public functions

// private functions

func validateTableName(tableName string) error {
	if !db_rules.TableNameRulesCheck(tableName) {
		return fmt.Errorf("invalid table name '%s'", tableName)
	}
	return nil
}

func validateFieldName(fieldName string) error {
	if !db_rules.TableFieldNameRulesCheck(fieldName) {
		return fmt.Errorf("invalid table field name '%s'", fieldName)
	}
	return nil
}

func validateFieldMetaData(fi *createTableExternalInterfaceFieldInfo, pkKeyName string, pkIsSet *bool) (*tableFieldProperties, error) {

	fp := newTableFieldProperties()

	if e := validateFieldName(fi.FieldName); e != nil {
		return nil, e
	}

	fp.FieldName = fi.FieldName

	if dt, e := getDbType(fi.FieldType); e != nil {
		return nil, e
	} else {
		fp.FieldType = dt
	}

	if dbsystem.SystemCasingHandler.AreEqual(fp.FieldName, pkKeyName) {
		if *pkIsSet {
			return nil, errors.New("table can only have one primary key")
		}
		*pkIsSet = true
	}

	return fp, nil
}

//validates and creates main object, reasons:
//1) better code, since adding fields automatically checks certain constraints.
//2) optimization, since most of the time validation is followed by creation.
//- HY 26-Dec-2019
func validateCreateTableMetaData(createTableData *createTableExternalInterface) (*tableMain, error) {
	//first quick checks

	if e := validateTableName(createTableData.TableName); e != nil {
		return nil, e
	}

	//quick checks end

	pkIsSet := false
	dbTbl := newTableMain(createTableData.TableName)

	for _, fi := range createTableData.TableFields {
		//_ = i

		var fp *tableFieldProperties
		var e error

		if fp, e = validateFieldMetaData(&fi, createTableData.PrimaryKeyName, &pkIsSet); e != nil {
			return nil, e
		}

		if e = dbTbl.TableFieldsMetaData.add(fp, dbsystem.SystemCasingHandler); e != nil {
			return nil, e
		}

	}

	if !pkIsSet {
		return nil, errors.New("table must have primary key")
	}

	return dbTbl, nil
}

// NOTE: TableRow is a map, so no need to pass it as pointer
// WARNING: TableRow (by field name) is not returned unless function succeeds. So don't override r in calling function.
func validateRowData(t *tableMain, r TableRow) (TableRow, tableRowByInternalIds, error)  {
	rowByInternalId, e := fromLabeledByFieldNames(r, t, dbsystem.SystemCasingHandler)
	if e != nil {
		return nil, nil, e
	}

	rowConvertedWithCorrectTypes, e := toLabeledByFieldNames(rowByInternalId, t)
	if e != nil {
		return nil, nil, e
	}

	return rowConvertedWithCorrectTypes, rowByInternalId, nil
}

