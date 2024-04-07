package handlers

import (
	"errors"
	"regexp"
)

func getLatAndLonFromWKT(wkt string) (string, string, error) {
	re := regexp.MustCompile(`POINT\(([-+]?[0-9]*\.?[0-9]+) ([-+]?[0-9]*\.?[0-9]+)\)`)
	matches := re.FindStringSubmatch(wkt)

	if len(matches) != 3 {
		return "", "", errors.New("invalid WKT format")
	}

	return matches[1], matches[2], nil
}
