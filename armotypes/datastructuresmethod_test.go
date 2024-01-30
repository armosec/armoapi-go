package armotypes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateOrderBy(t *testing.T) {
	req := &V2ListRequest{}
	req.ValidateOrderBy("dddd")
	if req.OrderBy != "dddd:desc" {
		t.Errorf("order by validation failed")
	}
	req.ValidateOrderBy("hyuii")
	if req.OrderBy != "dddd:desc" {
		t.Errorf("order by validation failed")
	}
	req = &V2ListRequest{}
	req.ValidateOrderBy("")
	if req.OrderBy != "" {
		t.Errorf("order by validation failed")
	}
}

func TestValidatePageProperties(t *testing.T) {
	req := &V2ListRequest{}
	req.ValidatePageProperties(235)
	if *req.PageNum != 0 || *req.PageSize != 235 {
		t.Errorf("page properties validation failed")
	}
	req.fixedPageNum = false
	*req.PageNum = 43
	req.ValidatePageProperties(154)
	if *req.PageNum != 42 || *req.PageSize != 154 {
		t.Errorf("page properties validation failed")
	}
	//validate again wihtout reseting the validate flag - make sure ignored
	req.ValidatePageProperties(500)
	if *req.PageNum != 42 || *req.PageSize != 154 {
		t.Errorf("page properties validation failed")
	}
	req.fixedPageNum = false
	req.ValidatePageProperties(-45)
	if *req.PageNum != 41 || *req.PageSize != 154 {
		t.Errorf("page properties validation failed")
	}

}

func TestGetFieldsNames(t *testing.T) {
	req := &V2ListRequest{}
	fields := req.GetFieldsNames()
	if len(fields) != 0 {
		t.Errorf("wrong len(fields):%v", len(fields))
	}
	req.InnerFilters = []map[string]string{{"aaa": "dsds", "nnmou": "vvvvv"}, {"aaa": "mmmm", "ccc": "dddd"}}
	req.OrderBy = "bbb:desc,oiu:asc,uiotree:desc"
	fields = req.GetFieldsNames()
	if len(fields) != 7 {
		t.Errorf("wrong len(fields):%v", len(fields))
	}
	for _, field := range []string{"aaa", "ccc", "nnmou", "bbb", "oiu", "uiotree"} {
		found := false
		for _, fieldA := range fields {
			if fieldA == field {
				found = true
			}
		}
		if !found {
			t.Errorf("field %s not found in %v", field, fields)
		}
	}
}

func TestReplaceFieldsToKeywords(t *testing.T) {
	req := &V2ListRequest{}
	req.InnerFilters = []map[string]string{{"aaa": "dsds", "nnmou": "vvvvv"}, {"aaa": "mmmm", "ccc": "dddd"}, {"aaa.new": "fdsfdsf"}}
	req.OrderBy = "bbb:desc,aaa:asc,oiu:asc,uiotree:desc"
	req.ReplaceFieldsToKeywords(map[string]string{"aaa": "aaa.new", "uiotree": "uiotree.older"})
	if req.OrderBy != "bbb:desc,aaa.new:asc,oiu:asc,uiotree.older:desc" {
		t.Errorf("wrong fixed order by string")
	}
	if req.InnerFilters[0]["aaa.new"] != "dsds" || req.InnerFilters[0]["aaa"] != "" {
		t.Errorf("ReplaceFieldsToKeywords of aaa failed")
	}
	if req.InnerFilters[1]["aaa.new"] != "mmmm" || req.InnerFilters[1]["mmmm"] != "" {
		t.Errorf("ReplaceFieldsToKeywords of aaa failed")
	}
	if req.InnerFilters[0]["nnmou"] != "vvvvv" || req.InnerFilters[1]["ccc"] != "dddd" {
		t.Errorf("ReplaceFieldsToKeywords of aaa failed")
	}
	if req.InnerFilters[2]["aaa.new"] != "fdsfdsf" {
		t.Errorf("ReplaceFieldsToKeywords of aaa.new failed")
	}
}

func TestDurationString(t *testing.T) {
	d := Duration(60 * time.Second)
	assert.Equal(t, "1m", d.String())

	d = Duration(70 * time.Second)
	assert.Equal(t, "1m10s", d.String())
}
