package home

func (h *HomeService) Index() interface{} {
	return map[string]interface{}{
		"status": true,
	}
}
