package indexer

func aggByStatus(adrs []ADR) map[string][]string {
	return collectByFile(adrs, func(adr ADR) string { return adr.Status })
}

func aggByTags(adrs []ADR) map[string][]string {
	return collectListByFile(adrs, func(adr ADR) []string { return adr.Tags })
}

func aggByDriver(adrs []ADR) map[string][]string {
	return collectListByFile(adrs, func(adr ADR) []string { return adr.Driver })
}

func aggByInvolvement(adrs []ADR) map[string][]string {
	involvementByFile := make(map[string][]string)
	for _, adr := range adrs {
		allInvolved := append(append(append(adr.Driver, adr.Deciders...), adr.Consulted...), adr.Informed...)
		for _, person := range allInvolved {
			if person == "" {
				person = "undefined"
			}
			involvementByFile[person] = append(involvementByFile[person], allInvolved...)
		}
	}
	return involvementByFile
}

func aggByExtraProperties(adrs []ADR) map[string][]string {
	extraByFile := make(map[string][]string)
	for _, adr := range adrs {
		for key := range adr.Extra {
			extraByFile[key] = append(extraByFile[key], adr.File)
		}
	}
	return extraByFile
}

func getDocTermMatrix(adrs []ADR) map[string][]string {
	docTermMatrix := make(map[string][]string)
	for _, adr := range adrs {
		for _, term := range adr.Terms {
			docTermMatrix[term] = append(docTermMatrix[term], adr.File)
		}
	}
	return docTermMatrix
}

func collectByFile(adrs []ADR, selector func(ADR) string) map[string][]string {
	collections := make(map[string][]string)
	for _, adr := range adrs {
		key := selector(adr)
		if key == "" {
			key = "undefined"
		}
		collections[key] = append(collections[key], adr.File)
	}
	return collections
}

func collectListByFile(adrs []ADR, selector func(ADR) []string) map[string][]string {
	collections := make(map[string][]string)
	for _, adr := range adrs {
		keys := selector(adr)
		for _, key := range keys {
			if key == "" {
				key = "undefined"
			}
			collections[key] = append(collections[key], adr.File)
		}
	}
	return collections
}
