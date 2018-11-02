package chain

import (
	//"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
)

type IdxDouble struct {
	context  *ApplyContext
	itrCache *iteratorCache
}

func NewIdxDouble(c *ApplyContext) *IdxDouble {
	return &IdxDouble{
		context:  c,
		itrCache: NewIteratorCache(),
	}
}

func (i *IdxDouble) store(scope int64, table int64, payer int64, id int64, secondary *float64) int {
	EosAssert(common.AccountName(payer) != common.AccountName(0), &InvalidTablePayer{}, "must specify a valid account to pay for new record")
	tab := i.context.FindOrCreateTable(int64(i.context.Receiver), scope, table, payer)

	obj := entity.SecondaryObjectDouble{
		TId:          tab.ID,
		PrimaryKey:   uint64(id),
		SecondaryKey: *secondary,
		Payer:        common.AccountName(payer),
	}

	err := i.context.DB.Insert(&obj)
	if err != nil {
		return -1
	}

	i.context.DB.Modify(tab, func(t *entity.TableIdObject) {
		t.Count++
	})

	i.context.UpdateDbUsage(common.AccountName(payer), int64(common.BillableSizeV("index_double_object")))

	i.itrCache.cacheTable(tab)
	return i.itrCache.add(obj)
}

func (i *IdxDouble) remove(iterator int) {
	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)
	i.context.UpdateDbUsage(obj.Payer, -int64(common.BillableSizeV("index_double_object")))

	tab := i.itrCache.getTable(obj.TId)
	EosAssert(tab.Code == i.context.Receiver, &TableAccessViolation{}, "db access violation")

	i.context.DB.Modify(tab, func(t *entity.TableIdObject) {
		t.Count--
	})

	i.context.DB.Remove(obj)
	if tab.Count == 0 {
		i.context.DB.Remove(tab)
	}
	i.itrCache.remove(iterator)
}

func (i *IdxDouble) update(iterator int, payer int64, secondary *float64) {
	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)
	objTable := i.itrCache.getTable(obj.TId)
	EosAssert(objTable.Code == i.context.Receiver, &TableAccessViolation{}, "db access violation")

	accountPayer := common.AccountName(payer)
	if accountPayer == common.AccountName(0) {
		accountPayer = obj.Payer
	}

	billingSize := int64(common.BillableSizeV("index_double_object"))
	if obj.Payer != accountPayer {
		i.context.UpdateDbUsage(obj.Payer, -billingSize)
		i.context.UpdateDbUsage(accountPayer, +billingSize)
	}

	i.context.DB.Modify(obj, func(o *entity.SecondaryObjectDouble) {
		o.SecondaryKey = *secondary
		o.Payer = accountPayer
	})
}

func (i *IdxDouble) findSecondary(code int64, scope int64, table int64, secondary *float64, primary *uint64) int {
	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.SecondaryObjectDouble{TId: tab.ID, SecondaryKey: *secondary}
	err := i.context.DB.Find("bySecondary", obj, &obj)

	*primary = obj.PrimaryKey

	if err != nil {
		return tableEndItr
	}
	return i.itrCache.add(&obj)
}

func (i *IdxDouble) lowerbound(code int64, scope int64, table int64, secondary *float64, primary *uint64) int {
	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.SecondaryObjectDouble{TId: tab.ID, SecondaryKey: *secondary}

	idx, _ := i.context.DB.GetIndex("bySecondary", &obj)
	itr, _ := idx.LowerBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}

	objLowerbound := entity.SecondaryObjectDouble{}
	itr.Data(&objLowerbound)
	if objLowerbound.TId != tab.ID {
		return tableEndItr
	}

	*primary = objLowerbound.PrimaryKey
	*secondary = objLowerbound.SecondaryKey

	return i.itrCache.add(&objLowerbound)
}

func (i *IdxDouble) upperbound(code int64, scope int64, table int64, secondary *float64, primary *uint64) int {
	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.SecondaryObjectDouble{TId: tab.ID, SecondaryKey: *secondary}

	idx, _ := i.context.DB.GetIndex("bySecondary", &obj)
	itr, _ := idx.UpperBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}

	objUpperbound := entity.SecondaryObjectDouble{}
	itr.Data(&objUpperbound)
	if objUpperbound.TId != tab.ID {
		return tableEndItr
	}

	*primary = objUpperbound.PrimaryKey
	*secondary = objUpperbound.SecondaryKey

	return i.itrCache.add(&objUpperbound)
}

func (i *IdxDouble) end(code int64, scope int64, table int64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}
	return i.itrCache.cacheTable(tab)
}

func (i *IdxDouble) next(iterator int, primary *uint64) int {

	if iterator < -1 {
		return -1
	}
	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)

	idx, _ := i.context.DB.GetIndex("bySecondary", obj)
	itr := idx.IteratorTo(obj)

	itr.Next()
	objNext := entity.SecondaryObjectDouble{}
	itr.Data(&objNext)

	if idx.CompareEnd(itr) || objNext.TId != obj.TId {
		return i.itrCache.getEndIteratorByTableID(obj.TId)
	}

	*primary = objNext.PrimaryKey
	return i.itrCache.add(&objNext)

}

func (i *IdxDouble) previous(iterator int, primary *uint64) int {
	idx, _ := i.context.DB.GetIndex("bySecondary", &entity.SecondaryObjectDouble{})

	if iterator < -1 {
		tab := i.itrCache.findTablebyEndIterator(iterator)
		EosAssert(tab != nil, &InvalidTableTterator{}, "not a valid end iterator")

		objTId := entity.SecondaryObjectDouble{TId: tab.ID}

		itr, _ := idx.UpperBound(&objTId)
		if idx.CompareIterator(idx.Begin(), idx.End()) || idx.CompareBegin(itr) {
			return -1
		}

		itr.Prev()
		objPrev := entity.KeyValueObject{}
		itr.Data(&objPrev)

		if objPrev.TId != tab.ID {
			return -1
		}

		*primary = objPrev.PrimaryKey
		return i.itrCache.add(&objPrev)
	}

	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)
	itr := idx.IteratorTo(obj)

	if idx.CompareBegin(itr) {
		return -1
	}

	itr.Prev()
	objPrev := entity.SecondaryObjectDouble{}
	itr.Data(&objPrev)

	if objPrev.TId != obj.TId {
		return -1
	}
	*primary = objPrev.PrimaryKey
	return i.itrCache.add(&objPrev)
}

func (i *IdxDouble) findPrimary(code int64, scope int64, table int64, secondary *float64, primary *uint64) int {
	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.SecondaryObjectDouble{TId: tab.ID, PrimaryKey: *primary}
	err := i.context.DB.Find("byPrimary", obj, &obj)

	*secondary = obj.SecondaryKey

	if err != nil {
		return tableEndItr
	}
	return i.itrCache.add(&obj)
}

func (i *IdxDouble) lowerboundPrimary(code int64, scope int64, table int64, primary *uint64) int {
	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.SecondaryObjectDouble{TId: tab.ID, PrimaryKey: *primary}
	idx, _ := i.context.DB.GetIndex("byPrimary", &obj)

	itr, _ := idx.LowerBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}

	//objLowerbound := (*types.SecondaryObjectI64)(itr.GetObject())
	objLowerbound := entity.SecondaryObjectDouble{}
	itr.Data(&objLowerbound)

	if objLowerbound.TId != tab.ID {
		return tableEndItr
	}

	return i.itrCache.add(&objLowerbound)
}

func (i *IdxDouble) upperboundPrimary(code int64, scope int64, table int64, primary *uint64) int {
	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.SecondaryObjectDouble{TId: tab.ID, PrimaryKey: *primary}
	idx, _ := i.context.DB.GetIndex("byPrimary", &obj)
	itr, _ := idx.UpperBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}
	//objUpperbound := (*types.SecondaryObjectI64)(itr.GetObject())
	objUpperbound := entity.SecondaryObjectDouble{}
	itr.Data(&objUpperbound)
	if objUpperbound.TId != tab.ID {
		return tableEndItr
	}

	i.itrCache.cacheTable(tab)
	return i.itrCache.add(&objUpperbound)
}

func (i *IdxDouble) nextPrimary(iterator int, primary *uint64) int {
	if iterator < -1 {
		return -1
	}
	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)
	idx, _ := i.context.DB.GetIndex("byPrimary", obj)

	itr := idx.IteratorTo(obj)

	itr.Next()
	objNext := entity.SecondaryObjectDouble{}
	itr.Data(&objNext)

	if idx.CompareEnd(itr) || objNext.TId != obj.TId {
		return i.itrCache.getEndIteratorByTableID(obj.TId)
	}

	*primary = objNext.PrimaryKey
	return i.itrCache.add(&objNext)

}

func (i *IdxDouble) previousPrimary(iterator int, primary *uint64) int {

	idx, _ := i.context.DB.GetIndex("byPrimary", &entity.SecondaryObjectDouble{})

	if iterator < -1 {
		tab := i.itrCache.findTablebyEndIterator(iterator)
		EosAssert(tab != nil, &InvalidTableTterator{}, "not a valid end iterator")

		objTId := entity.SecondaryObjectDouble{TId: tab.ID}

		itr, _ := idx.UpperBound(&objTId)
		if idx.CompareIterator(idx.Begin(), idx.End()) || idx.CompareBegin(itr) {
			return -1
		}

		itr.Prev()
		objPrev := entity.SecondaryObjectDouble{}
		itr.Data(&objPrev)

		if objPrev.TId != tab.ID {
			return -1
		}

		*primary = objPrev.PrimaryKey
		return i.itrCache.add(&objPrev)
	}

	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)
	itr := idx.IteratorTo(obj)

	if idx.CompareBegin(itr) {
		return -1
	}

	itr.Prev()
	objNext := entity.SecondaryObjectDouble{}
	itr.Data(&objNext)

	if objNext.TId != obj.TId {
		return -1
	}
	*primary = objNext.PrimaryKey
	return i.itrCache.add(&objNext)
}

func (i *IdxDouble) get(iterator int, secondary *float64, primary *uint64) {
	obj := (i.itrCache.get(iterator)).(*entity.SecondaryObjectDouble)

	*primary = obj.PrimaryKey
	*secondary = obj.SecondaryKey
}
