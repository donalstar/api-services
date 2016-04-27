package legacy

func cannedResponsesToIface(cannedResponses []CannedResponse) []interface{} {

	if len(cannedResponses) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(cannedResponses))

	for i, v := range cannedResponses {
		ifs[i] = v
	}
	return ifs
}
