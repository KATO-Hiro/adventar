<template>
  <div>
    <GlobalHeader />

    <template v-if="calendar">
      <header :style="{ backgroundColor: calendarColor() }">
        <div class="inner">
          <h2>{{ calendar.title }} Advent Calendar {{ calendar.year }}</h2>
          <div>登録数 {{ calendar.entries.length }}/25人</div>
          <div>
            作成者 <UserIcon :user="calendar.owner" size="22" />
            {{ calendar.owner.name }}
          </div>
          <nuxt-link class="editBtn" :to="`/calendars/${calendar.id}/edit`" v-if="isOwnCalendar(calendar)">
            <font-awesome-icon icon="edit"></font-awesome-icon> 編集
          </nuxt-link>
        </div>
      </header>
      <main>
        <div>
          <VueMarkdown>{{ calendar.description }}</VueMarkdown>
          <CalendarTable cells="cells"></CalendarTable>
          <table class="entryList">
            <tr v-for="entry in calendar.entries" :key="entry.day">
              <th class="date">12/{{ entry.day }}</th>
              <td class="user">
                <UserIcon :user="entry.owner" size="24" />
                {{ entry.owner.name }}
              </td>
              <td class="body">
                <div v-if="entry.comment"><font-awesome-icon icon="comment" /> {{ entry.comment }}</div>
                <div v-if="entry.title && !isFutureEntry(entry)">
                  <font-awesome-icon icon="file" /> {{ entry.title }}
                </div>
                <div v-if="entry.url && !isFutureEntry(entry)">
                  <font-awesome-icon icon="link" />
                  <a :href="entry.url">{{ entry.url }}</a>
                </div>
              </td>
              <td class="image">
                <img :src="entry.imageUrl" v-if="entry.imageUrl && !isFutureEntry(entry)" width="100" height="100" />
              </td>
            </tr>
          </table>
        </div>
      </main>
    </template>
  </div>
</template>

<script lang="ts">
import dayjs from "dayjs";
import { Component, Vue } from "nuxt-property-decorator";
import VueMarkdown from "vue-markdown";
import { getCalendar, createEntry, updateEntry, deleteEntry } from "~/lib/GrpcClient";
import * as RestClient from "~/lib/RestClient";
import { calendarColor } from "~/lib/utils/Colors";
import { Calendar, Entry } from "~/types/adventar";
import { getToken } from "~/lib/Auth";
import GlobalHeader from "~/components/GlobalHeader.vue";
import UserIcon from "~/components/UserIcon.vue";

function getRows(calendar: Calendar): any[] {
  const year = calendar.year;
  const endDay = dayjs(new Date(year, 11, 25));
  let currentDay = dayjs(new Date(year, 11, 1)).startOf("week");

  const entryMapByDay: { [key: number]: Entry } = {};
  if (calendar !== null && calendar.entries !== undefined) {
    calendar.entries.forEach(entry => {
      if (entry.day) {
        entryMapByDay[entry.day] = entry;
      }
    });
  }

  const rows: any[] = [];
  while (currentDay <= endDay) {
    const cells: any[] = [];
    for (let i = 0; i < 7; i++) {
      const day = currentDay.date();
      const entry = entryMapByDay[day] || null;
      const entryable = currentDay.month() === 11 && day <= 25;
      const cell = { day, entry, entryable };
      cells.push(cell);
      currentDay = currentDay.add(1, "day");
    }
    rows.push(cells);
  }

  return rows;
}

@Component({
  components: { GlobalHeader, UserIcon, VueMarkdown }
})
export default class extends Vue {
  calendar: Calendar | null = null;
  rows: any[];

  async asyncData({ params }) {
    if (process.server) {
      const calendar = await RestClient.getCalendar(params.id);
      const rows = getRows(calendar);
      return { calendar, rows };
    }
  }

  async mounted() {
    this.calendar = await getCalendar(Number(this.$route.params.id));
    this.rows = getRows(this.calendar);
  }

  async refetchCalendar() {
    this.calendar = await getCalendar(this.calendar!.id);
    this.rows = getRows(this.calendar);
  }

  calendarColor(): string {
    return calendarColor(this.calendar!.id);
  }

  isOwnCalendar(calendar: Calendar): boolean {
    if (!this.$store.state.user) return false;
    if (!calendar.owner) return false;
    return calendar.owner.id === this.$store.state.user.id;
  }

  isFutureEntry(_): boolean {
    // TODO
    return false;
  }
}
</script>

<style lang="scss" scoped>
main > div {
  padding-top: 10px;
}

header {
  color: #fff;
}

header > .inner {
  max-width: $content-max-width;
  padding: 30px 12px;
  margin: 0 auto;
  position: relative;
}

header .userIcon {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  vertical-align: middle;
}

header h2 {
  margin: 0 0 20px 0;
  font-size: 25px;
  font-weight: bold;
}

header .editBtn {
  padding: 10px 20px;
  font-size: 12px;
  display: inline-block;
  border-radius: 5px;
  cursor: pointer;
  text-decoration: none;
  border: none;
  outline: none;
  background: rgba(255, 255, 255, 0.85);
  color: #333;
  position: absolute;
  right: 12px;
  bottom: 30px;
}

.description p {
  margin: 100px;
}

.entryList {
  border-collapse: collapse;
  border-spacing: 0;
  font-size: 15px;
  color: #666;
  width: 100%;
  margin-top: 50px;
}

.entryList td,
.entryList th {
  vertical-align: top;
  padding: 20px 10px;
  border-top: 1px solid #e3e4e3;
  background-color: #fcfcfc;
}

.entryList .date {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
  width: 70px;
}

.entryList .userIcon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  vertical-align: middle;
  margin-right: 5px;
}

.entryList .body {
  width: 570px;
}

.entryList .body > div {
  padding-bottom: 5px;
}

.entryList .body svg {
  margin-right: 5px;
  width: 15px;
}

.entryList .image {
  width: 100px;
  text-align: center;
  vertical-align: top;
}
</style>
