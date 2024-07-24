package snapshot

import (
	"github.com/doublecloud/tross/transfer_manager/go/pkg/abstract"
	"github.com/doublecloud/tross/transfer_manager/go/pkg/base"
	"github.com/doublecloud/tross/transfer_manager/go/pkg/providers/oracle/schema"
)

type OracleDataObject struct {
	table  *schema.Table
	closed bool
}

func NewOracleDataObject(table *schema.Table) *OracleDataObject {
	return &OracleDataObject{
		table:  table,
		closed: false,
	}
}

func (dataObject *OracleDataObject) Table() *schema.Table {
	return dataObject.table
}

// Begin of base DataObject interface

func (dataObject *OracleDataObject) Name() string {
	return dataObject.table.Name()
}

func (dataObject *OracleDataObject) FullName() string {
	return dataObject.table.FullName()
}

func (dataObject *OracleDataObject) Next() bool {
	if dataObject.closed {
		return false
	}
	dataObject.closed = true
	return true
}

func (dataObject *OracleDataObject) Err() error {
	return nil
}

func (dataObject *OracleDataObject) Close() {
	dataObject.closed = true
}

func (dataObject *OracleDataObject) Part() (base.DataObjectPart, error) {
	return dataObject, nil
}

func (dataObject *OracleDataObject) ToOldTableID() (*abstract.TableID, error) {
	tableID := &abstract.TableID{
		Namespace: dataObject.table.Schema(),
		Name:      dataObject.table.Name(),
	}
	return tableID, nil
}

// End of base DataObject interface

// Begin of base DataObjectPart interface

func (dataObject *OracleDataObject) ToOldTableDescription() (*abstract.TableDescription, error) {
	tableDesc := &abstract.TableDescription{
		Name:   dataObject.table.Name(),
		Schema: dataObject.table.Schema(),
		Filter: "",
		EtaRow: 0,
		Offset: 0,
	}
	return tableDesc, nil
}

func (dataObject *OracleDataObject) ToTablePart() (*abstract.TableDescription, error) {
	return dataObject.ToOldTableDescription()
}

// End of base DataObjectPart interface
