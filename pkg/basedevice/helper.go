package basedevice

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	"reflect"
)

// GetControlsInfo Общая функция для получения информации о контролах
func GetControlsInfo(controls interface{}) []control.Info {
	var infoList []control.Info

	// Получаем значение и тип переданной структуры
	controlsValue := reflect.ValueOf(controls).Elem()
	controlsType := controlsValue.Type()

	// Проходимся по всем полям структуры
	for i := 0; i < controlsValue.NumField(); i++ {
		field := controlsValue.Field(i)

		// Проверяем, что поле является указателем и не nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// Проверяем, реализует ли поле метод GetInfo
			method := field.MethodByName("GetInfo")
			if method.IsValid() {
				// Вызываем метод GetInfo
				info := method.Call(nil)[0].Interface().(control.Info)
				infoList = append(infoList, info)
			} else {
				fmt.Printf("Field %s does not implement GetInfo\n", controlsType.Field(i).Name)
			}
		}
	}

	return infoList
}
