# APIs for Web applications

## GET /search
### Request (URL Query)
| field | type   | possible value | note                                     |
|-------|--------|----------------|------------------------------------------|
| key   | string | "tissue paper" | NOT NULL                                 |
| mart  | string | "carrefour"    | If null, search all                      |
| num   | int    | 20             | If null, search all                      |
| order | int    | 0, 1           | 0: by price; default<br>1: by popularity | 

### Response (application/json; array of objects)
| field | type   | possible value | note                     |
|-------|--------|----------------|--------------------------|
| name  | string | "paper"        | product name             |
| image | string | "https://..."  | url of the product image |
| page  | string | "https://..."  | url of the product page  |
| price | int    | 89             | price of the product     |
| mart  | srting | "carrefour"    | name of the mart         |

### Example
GET /search?key=抽取衛生紙&mart=wellcome&num=2&order=0

```json
[
	{
		"name": "[得意]抽取衛生紙100抽x10包/袋",
		"image": "https://sbd-ec.wellcome.com.tw/fileHandler/show/2306?photoSize=480x480",
		"page": "https://sbd-ec.wellcome.com.tw/product/view/erR",
		"price": 100,
		"mart": "wellcome"
	}, {
		"name": "[舒潔]拉拉炫彩抽取衛生紙110抽x8包/串",
		"image": "https://sbd-ec.wellcome.com.tw/fileHandler/show/18059?photoSize=480x480",
		"page": "https://sbd-ec.wellcome.com.tw/product/view/QWoq",
		"price": 109,
		"mart": "wellcome"
	}
]
```

## GET /marts
### Request
no parameter needed

### Response (application/json; array of objects)
| field | type   | possible value | note                     |
|-------|--------|----------------|--------------------------|
| id    | string | "rt"           | abbreviation of the mart |
| name  | string | "RT-Mart"      | full name of the mart    |
| cur   | string | "TWD"          | currency of the prices   |

### Example
GET /marts

```json
[
	{
		"id": "rt",
		"name": "RT-Mart",
		"cur": "TWD"
	},{
		"id": "carrefour",
		"name": "Carrefour(TW)",
		"cur": "TWD"
	}
]
```