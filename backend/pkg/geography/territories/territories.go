package geography

func GetAllTerritories() ([]string, error) {
	// This function should return a list of all territories.
	// For now, we will return a hardcoded list.
	territories := []string{
		"USA",
		"Canada",
		"Mexico",
		"Brazil",
		"Argentina",
	}
	return territories, nil
}
