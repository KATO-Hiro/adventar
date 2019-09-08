<template>
  <ul>
    <li v-for="entry in calendar.entries" :key="entry.day">
      <div class="date">12/{{ entry.day }}</div>
      <div class="user">
        <UserIcon :user="entry.owner" size="24" />
        {{ entry.owner.name }}
      </div>
      <div class="body">
        <div v-if="entry.comment"><font-awesome-icon icon="comment" /> {{ entry.comment }}</div>
        <div v-if="entry.title && !isFutureEntry(entry)"><font-awesome-icon icon="file" /> {{ entry.title }}</div>
        <div v-if="entry.url && !isFutureEntry(entry)">
          <font-awesome-icon icon="link" />
          <a :href="entry.url">{{ entry.url }}</a>
        </div>
      </div>
      <div class="image">
        <img :src="entry.imageUrl" v-if="entry.imageUrl && !isFutureEntry(entry)" width="100" height="100" />
      </div>
    </li>
  </ul>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";
import UserIcon from "~/components/UserIcon.vue";
import { Calendar } from "~/types/adventar";

@Component({
  components: { UserIcon }
})
export default class extends Vue {
  @Prop() readonly calendar: Calendar;

  isFutureEntry(_): boolean {
    // TODO
    return false;
  }
}
</script>

<style scoped>
table {
  border-collapse: collapse;
  border-spacing: 0;
  font-size: 15px;
  color: #666;
  width: 100%;
  margin-top: 50px;
}

td,
th {
  vertical-align: top;
  padding: 20px 10px;
  border-top: 1px solid #e3e4e3;
  background-color: #fcfcfc;
}

.date {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
  width: 70px;
}

.userIcon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  vertical-align: middle;
  margin-right: 5px;
}

.body {
}

.body a {
  word-break: break-all;
}

.body > div {
  padding-bottom: 5px;
}

.body svg {
  margin-right: 5px;
  width: 15px;
}

.image {
  width: 100px;
  text-align: center;
  vertical-align: top;
}
</style>
