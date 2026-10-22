package minifunc

func MergeMaps(dest, src map[string]interface{}) {
	for key, srcVal := range src {
		if srcMap, srcIsMap := srcVal.(map[string]interface{}); srcIsMap {
			if destVal, destExists := dest[key]; destExists {
				if destMap, destIsMap := destVal.(map[string]interface{}); destIsMap {
					MergeMaps(destMap, srcMap)
					continue
				}
			}
		}
		dest[key] = srcVal
	}
}
