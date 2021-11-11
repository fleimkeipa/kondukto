package utils

func CheckSeverity(result map[string]interface{}) string {
	totals := result["metrics"].(map[string]interface{})["_totals"]
	for i, v := range totals.(map[string]interface{}) {
		if (i == "SEVERITY.HIGH" || i == "CONFIDENCE.HIGH") && v.(float64) > 1 {
			return "severity high metrics is bigger than 1"
		}
		if (i == "SEVERITY.MEDIUM" || i == "CONFIDENCE.MEDIUM") && v.(float64) > 1 {
			return "severity medium metrics is bigger than 2"
		}
		if (i == "SEVERITY.LOW" || i == "CONFIDENCE.LOW") && v.(float64) > 3 {
			return "severity low metrics is bigger than 3"
		}
	}
	return ""
}
