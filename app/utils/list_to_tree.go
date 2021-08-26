package utils


func ListToTree(data []map[string]interface{}) map[uint]map[string]interface{} {
	res := make(map[uint]map[string]interface{})
	for _, v := range data {
		id := uint(v["id"].(float64))
		parentID := uint(v["parent_id"].(float64))

		if res[id] != nil {
			v["children"] = res[id]["children"]
			res[id] = v
		} else {
			//v["children"] = []map[string]interface{}{}
			res[id] = v
		}
		if res[parentID] != nil {
			if res[parentID]["children"] !=nil{
				res[parentID]["children"] = append(
					res[parentID]["children"].([]map[string]interface{}),
					res[id],
				)
			} else {
				res[parentID]["children"] = []map[string]interface{}{res[id]}
			}

		} else {
			res[parentID] = map[string]interface{}{
				"children": []map[string]interface{}{
					res[id],
				},
			}
		}
	}
	//res[0]["children"].([]map[string]interface{})

	return res
}

func listToTree(data []map[string]interface{}) map[uint]map[string]interface{} {
	res := make(map[uint]map[string]interface{})
	for _, v := range data {
		id := uint(v["id"].(float64))
		parentID := uint(v["parent_id"].(float64))

		if res[id] != nil {
			v["children"] = res[id]["children"]
			res[id] = v
		} else {
			//v["children"] = []map[string]interface{}{}
			res[id] = v
		}
		if res[parentID] != nil {
			if res[parentID]["children"] !=nil{
				res[parentID]["children"] = append(
					res[parentID]["children"].([]map[string]interface{}),
					res[id],
				)
			} else {
				res[parentID]["children"] = []map[string]interface{}{res[id]}
			}

		} else {
			res[parentID] = map[string]interface{}{
				"children": []map[string]interface{}{
					res[id],
				},
			}
		}
	}
	//res[0]["children"].([]map[string]interface{})

	return res
}