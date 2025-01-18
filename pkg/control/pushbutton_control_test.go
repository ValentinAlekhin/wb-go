package control

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	"github.com/ValentinAlekhin/wb-go/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPushbuttonControl_Push(t *testing.T) {
	t.Parallel()

	client, server, destroy := testutils.GetClientWithBroker()
	testutils.AddOnHandler(server)
	defer destroy()

	meta := Meta{}
	device := testutils.RandString(10)
	controlName := testutils.RandString(10)

	// Создаем PushbuttonControl
	pushbuttonControl := NewPushbuttonControl(client, device, controlName, meta)

	// Проверяем начальное состояние
	initialValue := pushbuttonControl.control.GetValue()
	assert.Equal(t, "", initialValue, "Начальное значение должно быть пустым")

	// Нажимаем на кнопку (вызываем метод Push)
	pushbuttonControl.Push()

	// Даем время для обработки изменения значения
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение изменилось на "true"
	updatedValue := pushbuttonControl.control.GetValue()
	assert.Equal(t, conventions.CONV_META_BOOL_TRUE, updatedValue, "После нажатия кнопки значение должно быть '1'")

	// Пытаемся снова нажать кнопку
	pushbuttonControl.Push()

	// Даем время для обработки
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что значение не изменилось
	updatedValueAgain := pushbuttonControl.control.GetValue()
	assert.Equal(t, conventions.CONV_META_BOOL_TRUE, updatedValueAgain, "После второго нажатия значение не должно измениться")
}
