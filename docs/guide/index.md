<script setup>
import { data } from '../loaders/cli.data.ts'
const list = Object.entries(data).map(([cmd, value]) => ({ cmd, description: value.synopsis }))
</script>

# Руководство

## Список команд

<ul>
<li v-for="item of list"> <a :href="`/guide/commands/${item.cmd}`"><code>{{item.cmd}}</code> - {{item.description}}</a> </li>
</ul>
