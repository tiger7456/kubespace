/*




Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dataselect

import (
	"sort"
)

// GenericDataCell describes the interface of the data cell that contains all the necessary methods needed to perform
// complex data selection
// GenericDataSelect takes a list of these interfaces and performs selection operation.
// Therefore as long as the list is composed of GenericDataCells you can perform any data selection!
type DataCell interface {
	// GetPropertyAtIndex returns the property of this data cell.
	// Value returned has to have Compare method which is required by Sort functionality of DataSelect.
	GetProperty(PropertyName) ComparableValue
}

// ComparableValue hold any value that can be compared to its own kind.
type ComparableValue interface {
	// Compares self with other value. Returns 1 if other value is smaller, 0 if they are the same, -1 if other is larger.
	Compare(ComparableValue) int
	// Returns true if self value contains or is equal to other value, false otherwise.
	Contains(ComparableValue) bool
}

// SelectableData contains all the required data to perform data selection.
// It implements sort.Interface so its sortable under sort.Sort
// You can use its Select method to get selected GenericDataCell list.
type DataSelector struct {
	// GenericDataList hold generic data cells that are being selected.
	GenericDataList []DataCell
	// DataSelectQuery holds instructions for data select.
	DataSelectQuery *DataSelectQuery
}

// GenericDataSelectWithFilter 获取 GenericDataCells 和 DataSelectQuery 的列表，并按照 dsQuery 的指示返回选定的数据。
func GenericDataSelectWithFilter(dataList []DataCell, dsQuery *DataSelectQuery) ([]DataCell, int) {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}
	// Pipeline is Filter -> Sort -> CollectMetrics -> Paginate
	filtered := SelectableData.Filter()
	filteredTotal := len(filtered.GenericDataList)
	processed := filtered.Sort().Paginate()
	return processed.GenericDataList, filteredTotal
}

// Filter the data inside as instructed by DataSelectQuery and returns itself to allow method chaining.
func (self *DataSelector) Filter() *DataSelector {
	filteredList := []DataCell{}

	for _, c := range self.GenericDataList {
		matches := true
		for _, filterBy := range self.DataSelectQuery.FilterQuery.FilterByList {
			v := c.GetProperty(filterBy.Property)
			if v == nil || !v.Contains(filterBy.Value) {
				matches = false
				break
			}
		}
		if matches {
			filteredList = append(filteredList, c)
		}
	}

	self.GenericDataList = filteredList
	return self
}

// Implementation of sort.Interface so that we can use built-in sort function (sort.Sort) for sorting SelectableData

// Len returns the length of data inside SelectableData.
func (self DataSelector) Len() int { return len(self.GenericDataList) }

// Swap swaps 2 indices inside SelectableData.
func (self DataSelector) Swap(i, j int) {
	self.GenericDataList[i], self.GenericDataList[j] = self.GenericDataList[j], self.GenericDataList[i]
}

// Less compares 2 indices inside SelectableData and returns true if first index is larger.
func (self DataSelector) Less(i, j int) bool {
	for _, sortBy := range self.DataSelectQuery.SortQuery.SortByList {
		a := self.GenericDataList[i].GetProperty(sortBy.Property)
		b := self.GenericDataList[j].GetProperty(sortBy.Property)
		// ignore sort completely if property name not found
		if a == nil || b == nil {
			break
		}
		cmp := a.Compare(b)
		if cmp == 0 { // values are the same. Just continue to next sortBy
			continue
		} else { // values different
			return (cmp == -1 && sortBy.Ascending) || (cmp == 1 && !sortBy.Ascending)
		}
	}
	return false
}

// Sort sorts the data inside as instructed by DataSelectQuery and returns itself to allow method chaining.
func (self *DataSelector) Sort() *DataSelector {
	sort.Sort(*self)
	return self
}

// Paginates the data inside as instructed by DataSelectQuery and returns itself to allow method chaining.
func (self *DataSelector) Paginate() *DataSelector {
	pQuery := self.DataSelectQuery.PaginationQuery
	dataList := self.GenericDataList
	startIndex, endIndex := pQuery.GetPaginationSettings(len(dataList))

	// Return all items if provided settings do not meet requirements
	if !pQuery.IsValidPagination() {
		return self
	}
	// Return no items if requested page does not exist
	if !pQuery.IsPageAvailable(len(self.GenericDataList), startIndex) {
		self.GenericDataList = []DataCell{}
		return self
	}

	self.GenericDataList = dataList[startIndex:endIndex]
	return self
}

// IsValidPagination returns true if pagination has non negative parameters
func (p *PaginationQuery) IsValidPagination() bool {
	return p.ItemsPerPage >= 0 && p.Page >= 0
}

// IsPageAvailable returns true if at least one element can be placed on page. False otherwise
func (p *PaginationQuery) IsPageAvailable(itemsCount, startingIndex int) bool {
	return itemsCount > startingIndex && p.ItemsPerPage > 0
}

// GenericDataSelect takes a list of GenericDataCells and DataSelectQuery and returns selected data as instructed by dsQuery.
func GenericDataSelect(dataList []DataCell, dsQuery *DataSelectQuery) []DataCell {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}
	return SelectableData.Sort().Paginate().GenericDataList
}
