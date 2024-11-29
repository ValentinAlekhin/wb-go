<script setup>
import { useData } from 'vitepress'
const { params } = useData()
import { data } from '../../loaders/cli.data.ts'

const pageData = data[params.value.cmd]

const tableHeader = { name: 'Название', shorthand: 'Сокращение', default_value: 'Значение по умолчанию', usage: 'Описание' }
let options = pageData.options.map(opt => {
return  Object.keys(tableHeader).reduce((acc, key) => {
acc[key] = opt[key]
return acc
}, {})
})

options = [ Object.values(tableHeader), ...options]
</script>

# Команда `{{params.cmd}}`

{{pageData.description}}

## Использование

**Команда**

```shell-vue
{{pageData.usage}}
```

**Параметры**

<table>
  <tr v-for="row of options">
    <th v-for="opt of row" >{{opt}}</th>
  </tr>
</table>