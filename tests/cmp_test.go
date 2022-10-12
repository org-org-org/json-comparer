package tests

import (
	"fmt"
	JsonComparer "github.com/org-org-org/json-comparer"
	"testing"
)

func TestCmp(t *testing.T) {
	cmp := JsonComparer.NewComparer("uuid")
	cmp.IgnoreListSequence(true)
	ok, err := cmp.CompareJson(s1, s2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ok)
}

var (
	s1 = `{
		"code": 0,
		"msg": "操作成功",
		"data": {
			"uuid": "064c295F-12BB-22A1-bf2B-A4d5E376e65A",
			"name": "曹军",
			"createdAt": "1976-11-04 16:55:11"
		}
	}`
	s2 = `{
		"code": 0,
		"msg": "操作成功",
		"data": {
			"uuid": "064c295F-12BB-22A1-bf2B-A4d5E376e65A",
			"createdAt": "1976-11-04 16:55:11",
			"name": "曹军"
		}
	}`
	list1 = `{
    "code": 0,
    "msg": "操作成功",
    "data": {
        "list": [
			{
                "name": "代矿标族必又规",
                "uuid": "bEb932eB-ABAE-7CeA-9dAa-3C117Fcc8b3A",
                "enterprise": {
                    "name": "都大军开相收",
                    "uuid": "5F81cFFc-AB6d-b7Bd-eAad-959Fcb5cD7cD"
                },
                "floor": {
                    "number": 1,
                    "uuid": "6c1eDeAc-7fFc-dc4C-504B-e7cCD52beFF2"
                }
            },
			{
                "name": "代矿标族必又规2",
                "uuid": "bEb932eB-ABAE-7CeA-9dAa-3C117Fcc8b3A"
            }
        ],
        "page": 1,
        "size": 10,
        "total": 20
    }
}`
	list2 = `{
    "code": 0,
    "msg": "操作成功",
    "data": {
        "list": [
			{
                "name": "代矿标族必又规2",
                "uuid": "bEb932eB-ABAE-7CeA-9dAa-3C117Fcc8b3A"
            },
			{
                "name": "代矿标族必又规",
                "uuid": "bEb932eB-ABAE-7CeA-9dAa-3C117Fcc8b3A",
                "enterprise": {
                    "name": "都大军开相收",
                    "uuid": "5F81cFFc-AB6d-b7Bd-eAad-959Fcb5cD7cD"
                },
                "floor": {
                    "number": 1,
                    "uuid": "6c1eDeAc-7fFc-dc4C-504B-e7cCD52beFF2"
                }
            }
        ],
        "page": 1,
        "size": 10,
        "total": 20
    }
}`
)
