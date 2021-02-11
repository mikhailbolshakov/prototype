package impl

import (
	"context"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"strings"
)

func (s *serviceImpl) AskBot(ctx context.Context, r *domain.AskBotRequest) (*domain.AskBotResponse, error) {

	rs := &domain.AskBotResponse{}

	var answerMap = map[string]func() string{
		"привет":        s.greeting,
		"добрый день":   s.greeting,
		"доктор онлайн": s.doctorOnlineInfo,
		"ргс":           s.rgsInfo,
	}

	msg := strings.ToLower(r.Message)
	for word, f := range answerMap {
		if strings.Contains(msg, word) {
			rs.Answer = f()
			rs.Found = true
			break
		}
	}

	return rs, nil
}

func (s *serviceImpl) greeting() string {
	return "Добрый день, я бот РГС!\nЧем я могу Вам помочь?"
}

func (s *serviceImpl) doctorOnlineInfo() string {
	return "**«Доктор онлайн»** — полис для получения удаленной медицинской помощи.\nВы в любое время сможете связаться с нужными специалистами, чтобы получить консультацию.\n«Приемы» проводятся онлайн, на сайте https://med.moi-service.ru."
}

func (s *serviceImpl) rgsInfo() string {
	return "**ПАО СК «Росгосстрах»** — флагман отечественного рынка страхования.\nНа территории Российской Федерации действуют около 1 500 офисов и представительств компании, порядка 300 центров и пунктов урегулирования убытков.\nВ компании работает около 50 тысяч сотрудников и страховых агентов.\n\n«Росгосстрах» является безусловным лидером и системообразующей основой рынка страхования Российской Федерации.\n\n Подробнее: https://www.rgs.ru/"
}

