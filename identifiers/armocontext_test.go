package identifiers

import (
	"sort"
	"testing"
)

func TestIsDesignatorsMatchContext(t *testing.T) {
	desi := &PortalDesignator{Attributes: map[string]string{
		"aaa":   "bbb",
		"cc34":  "ert56",
		"mmn90": "zvbnvb",
	},
		WLID: "asdf"}
	ctxSlice := DesignatorToArmoContext(desi, "some2")
	if !IsDesignatorsMatchContext(ctxSlice, desi, "some2") {
		t.Errorf("excpected match")
	}
	if !IsDesignatorsMatchContext(ctxSlice, desi, "some222") {
		t.Errorf("excpected match")
	}

	desi.Attributes["bbb"] = "vxcvcx"
	if !IsDesignatorsMatchContext(ctxSlice, desi, "some2") {
		t.Errorf("excpected match")
	}

	desi.Attributes["aaa"] = "vxcvcx"
	if IsDesignatorsMatchContext(ctxSlice, desi, "some2") {
		t.Errorf("excpected missmatch")
	}
	desi.Attributes["aaa"] = "bbb"
	desi.WLID = "mmmmmmmmmmmmmmk"
	if IsDesignatorsMatchContext(ctxSlice, desi, "some2") {
		t.Errorf("excpected missmatch")
	}
}

func TestDesignatorToArmoContext(t *testing.T) {
	desiCtx := DesignatorToArmoContext(&PortalDesignator{Attributes: map[string]string{
		"aaa":   "bbb",
		"cc34":  "ert56",
		"mmn90": "zvbnvb",
	}}, "some")
	sort.SliceStable(desiCtx, func(i, j int) bool {
		return desiCtx[i].Attribute < desiCtx[j].Attribute
	})
	if desiCtx[0].Value != "bbb" || desiCtx[1].Attribute != "cc34" || desiCtx[2].Source != "some.attributes" {
		t.Errorf("wrong ctx:%v", desiCtx)
	}

	desiCtx = DesignatorToArmoContext(&PortalDesignator{Attributes: map[string]string{
		"aaa":   "bbb",
		"cc34":  "ert56",
		"mmn90": "zvbnvb",
	}}, "")
	sort.SliceStable(desiCtx, func(i, j int) bool {
		return desiCtx[i].Attribute < desiCtx[j].Attribute
	})
	if desiCtx[0].Value != "bbb" || desiCtx[1].Attribute != "cc34" || desiCtx[2].Source != "attributes" {
		t.Errorf("wrong ctx:%v", desiCtx)
	}

	desiCtx = DesignatorToArmoContext(&PortalDesignator{
		WLID:     "aaabdd",
		WildWLID: "fdsfdsd",
		SID:      "sssioop"}, "some1")
	sort.SliceStable(desiCtx, func(i, j int) bool {
		return desiCtx[i].Attribute > desiCtx[j].Attribute
	})
	if desiCtx[0].Value != "aaabdd" || desiCtx[0].Attribute != "wlid" || desiCtx[0].Source != "some1" {
		t.Errorf("wrong WLID ctx:%v", desiCtx)
	}
	if desiCtx[1].Value != "fdsfdsd" || desiCtx[1].Attribute != "wildwlid" || desiCtx[1].Source != "some1" {
		t.Errorf("wrong WildWLID ctx:%v", desiCtx)
	}
	if desiCtx[2].Value != "sssioop" || desiCtx[2].Attribute != "sid" || desiCtx[2].Source != "some1" {
		t.Errorf("wrong SID ctx:%v", desiCtx)
	}
}
